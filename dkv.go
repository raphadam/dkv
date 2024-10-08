package dkv

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

const raftTimeout = 10 * time.Second

var ErrKeyDoesNotExist = errors.New("key does not exist")

func init() {
	gob.Register(&Command{})
}

type CommandType int

const (
	SET CommandType = iota
	GET
	DEL
)

type Command struct {
	Cmd CommandType
	Key string
	Val string
}

type DKV struct {
	kv   map[string]string
	mu   sync.RWMutex
	raft *raft.Raft
}

func newDKV(single bool, raftAddr string) (*DKV, error) {
	d := &DKV{
		kv: make(map[string]string),
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftAddr)

	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, err
	}

	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()
	snapshot := raft.NewInmemSnapshotStore()

	transport, err := raft.NewTCPTransport(raftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	r, err := raft.NewRaft(config, d, logStore, stableStore, snapshot, transport)
	if err != nil {
		return nil, err
	}
	d.raft = r

	if single {
		d.raft.BootstrapCluster(raft.Configuration{Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		}})
	}

	return d, nil
}

type NotLeaderError struct {
	LeaderAddr string
}

func (e *NotLeaderError) Error() string {
	return "not leader"
}

func (d *DKV) Set(k string, v string) error {
	if d.raft.State() != raft.Leader {
		addr, _ := d.raft.LeaderWithID()
		return &NotLeaderError{LeaderAddr: string(addr)}
	}

	cmd := Command{
		Cmd: SET,
		Key: k,
		Val: v,
	}

	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(&cmd)
	if err != nil {
		return err
	}

	future := d.raft.Apply(buf.Bytes(), raftTimeout)
	return future.Error()
}

func (d *DKV) Get(k string) (string, error) {
	// TODO: eventual consistency?
	// if d.raft.State() != raft.Leader {
	// 	addr, _ := d.raft.LeaderWithID()

	// 	return "", &NotLeaderError{LeaderAddr: string(addr)}
	// }

	d.mu.RLock()
	defer d.mu.RUnlock()

	v, ok := d.kv[k]
	if !ok {
		return "", ErrKeyDoesNotExist
	}

	return v, nil
}

func (d *DKV) Del(k string) error {
	if d.raft.State() != raft.Leader {
		addr, _ := d.raft.LeaderWithID()

		return &NotLeaderError{LeaderAddr: string(addr)}
	}

	cmd := Command{
		Cmd: DEL,
		Key: k,
	}

	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(&cmd)
	if err != nil {
		return err
	}

	future := d.raft.Apply(buf.Bytes(), raftTimeout)
	return future.Error()
}

func (d *DKV) Join(raftAddr string) error {
	if d.raft.State() != raft.Leader {
		addr, _ := d.raft.LeaderWithID()
		return &NotLeaderError{LeaderAddr: string(addr)}
	}

	futureConfig := d.raft.GetConfiguration()
	if err := futureConfig.Error(); err != nil {
		log.Println("unable to  get config")
		return err
	}

	config := futureConfig.Configuration()
	for _, server := range config.Servers {
		if server.ID == raft.ServerID(raftAddr) || server.Address == raft.ServerAddress(raftAddr) {

			if server.ID == raft.ServerID(raftAddr) && server.Address == raft.ServerAddress(raftAddr) {
				return nil
			}

			future := d.raft.RemoveServer(server.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error %s %s", raftAddr, err)
			}
		}
	}

	future := d.raft.AddVoter(raft.ServerID(raftAddr), raft.ServerAddress(raftAddr), 0, 0)
	if err := future.Error(); err != nil {
		return err
	}

	return nil
}

func (d *DKV) Apply(l *raft.Log) interface{} {
	var req Command

	err := gob.NewDecoder(bytes.NewReader(l.Data)).Decode(&req)
	if err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	switch req.Cmd {
	case SET:
		d.kv[req.Key] = req.Val

	case DEL:
		delete(d.kv, req.Key)

	default:
		log.Fatal("not handled operation")
	}

	return nil
}

func (d *DKV) Snapshot() (raft.FSMSnapshot, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	clone := maps.Clone(d.kv)
	snapshot := &snapshot{kv: clone}

	return snapshot, nil
}

func (d *DKV) Restore(rc io.ReadCloser) error {
	var store map[string]string

	err := gob.NewDecoder(rc).Decode(&store)
	if err != nil {
		return err
	}

	d.kv = store
	return nil
}

type snapshot struct {
	kv map[string]string
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	err := gob.NewEncoder(sink).Encode(&s.kv)
	if err != nil {
		sink.Cancel()
		return err
	}

	return sink.Close()
}

func (s *snapshot) Release() {
}

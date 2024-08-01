package dkv

// create load balancer
// create gateway
// sharding
// encrypted password manager cli
// a tool to easly check the hash of programs
// promotheus
// ecnrypted file storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/raft"
)

const raftTimeout = 10 * time.Second

type joinRequest struct {
	Addr string `json:"addr"`
	Id   string `json:"id"`
}

type joinResponse struct {
}

type CommandRequest struct {
	Cmd string `json:"cmd"`
	Key string `json:"key"`
	Val string `json:"val"`
}

type CommandResponse struct {
}

type node struct {
	kv   map[string]string
	mu   sync.RWMutex
	raft *raft.Raft
	// log slog.Logger
}

func Serve(single bool, httpAddr string, raftAddr string) error {
	n := &node{
		kv: make(map[string]string),
	}

	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftAddr)

	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return err
	}

	logStore := raft.NewInmemStore()
	stableStore := raft.NewInmemStore()
	snapshot := raft.NewInmemSnapshotStore()

	transport, err := raft.NewTCPTransport(raftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	r, err := raft.NewRaft(config, n, logStore, stableStore, snapshot, transport)
	if err != nil {
		return err
	}
	n.raft = r

	if single {
		n.raft.BootstrapCluster(raft.Configuration{Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		}})
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", n.HttpCommand)
	mux.HandleFunc("POST /join", n.HttpJoin)
	return http.ListenAndServe(httpAddr, mux)
}

func (n *node) Apply(l *raft.Log) interface{} {
	var req CommandRequest

	err := json.Unmarshal(l.Data, &req)
	if err != nil {
		log.Fatal("it's supposed to work")
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	n.kv[req.Key] = req.Val
	return nil
}

func (n *node) Snapshot() (raft.FSMSnapshot, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	clone := maps.Clone(n.kv)
	snapshot := &snapshot{kv: clone}

	return snapshot, nil
}

func (n *node) Restore(rc io.ReadCloser) error {
	var store map[string]string

	err := json.NewDecoder(rc).Decode(&store)
	if err != nil {
		return err
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	n.kv = store

	return nil
}

type snapshot struct {
	kv map[string]string
}

func (s *snapshot) Persist(sink raft.SnapshotSink) error {
	b, err := json.Marshal(s.kv)
	if err != nil {
		return err
	}

	if _, err := sink.Write(b); err != nil {
		sink.Cancel()
		return err
	}

	return sink.Close()
}

func (s *snapshot) Release() {
}

func (n *node) Set(key string, val string) error {
	if n.raft.State() != raft.Leader {
		log.Println("not the leader")
		return fmt.Errorf("not the leader")
	}
	log.Println("i am the leader")

	// addr, id := n.raft.LeaderWithID()
	// log.Println(addr, id)
	// raft.Leader

	// n.raft.State()

	// n.raft.Apply()
	return nil
}

func (n *node) HttpCommand(w http.ResponseWriter, r *http.Request) {
	cmd := CommandRequest{}

	err := json.NewDecoder(r.Body).Decode(&cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(cmd)

	switch cmd.Cmd {
	case "set":
		n.Set(cmd.Key, cmd.Val)
	case "get":
		log.Println("get")
	case "del":
		log.Println("del")
	}
}

func (n *node) HttpJoin(w http.ResponseWriter, r *http.Request) {
	req := joinRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println(req)
	// log.Println(cmd)

	// switch cmd.Cmd {
	// case "set":
	// 	n.Set(cmd.Key, cmd.Val)
	// case "get":
	// 	log.Println("get")
	// case "del":
	// 	log.Println("del")
	// }
}

// data, err := io.ReadAll(r.Body)
// if err != nil {
// 	w.WriteHeader(http.StatusBadRequest)
// 	return
// }

// future := n.consensun.Apply(data, 5*time.Second)
// if future.Error() != nil {
// 	w.WriteHeader(http.StatusnodeUnavailable)
// }

// w.WriteHeader(http.StatusOK)
// w.Write([]byte("it is working great set"))

// func (n *node) get(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("it is working great get"))
// }

// func Serve(addr string) error {
// 	mux := http.NewServeMux()

// 	s := node{
// 		kv: make(map[string]string),
// 		// consensus: raft.NewRaft(raft.DefaultConfig()),
// 	}

// 	mux.HandleFunc("POST /set", n.set)
// 	mux.HandleFunc("GET /get", n.get)

//		return http.ListenAndServe(addr, mux)
//	}

package dkv

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/hashicorp/serf/serf"
)

type CommandRequest struct {
	Cmd string `json:"cmd"`
	Key string `json:"key"`
	Val string `json:"val"`
}

type CommandResponse struct {
	Value string `json:"val"`
}

type CommandError struct {
	Error string `json:"error"`
}

type node struct {
	store    *DKV
	serf     *serf.Serf
	eventCh  chan serf.Event
	serfAddr string
}

func Serve(single bool, restAddr string, serfAddr string, raftAddr string, bootstrap []string) error {
	// setup serf
	addr, err := net.ResolveTCPAddr("tcp", serfAddr)
	if err != nil {
		return err
	}

	// TODO: who is responsible to close ?
	eventCh := make(chan serf.Event)

	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port
	config.EventCh = eventCh
	config.Tags = map[string]string{
		"serfAddr": serfAddr,
		"raftAddr": raftAddr,
		"restAddr": restAddr,
	}
	config.NodeName = serfAddr

	serf, err := serf.Create(config)
	if err != nil {
		return err
	}

	// setup raft
	store, err := newDKV(single, raftAddr)
	if err != nil {
		return err
	}

	n := &node{
		serf:     serf,
		store:    store,
		eventCh:  eventCh,
		serfAddr: serfAddr,
	}

	go n.handleEvent()

	if bootstrap != nil {
		i, err := n.serf.Join(bootstrap, true)
		log.Println("joined: ", i, "error:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", n.handleCommand)
	return http.ListenAndServe(restAddr, mux)
}

// func (n *node) Leave() error {
// 	return n.serf.Leave()
// }

func (n *node) handleEvent() {
	for e := range n.eventCh {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if n.serf.LocalMember().Name == member.Name {
					continue
				}

				raftAddr := member.Tags["raftAddr"]
				n.store.Join(raftAddr)
			}
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if n.serf.LocalMember().Name == member.Name {
					continue
				}

				log.Println("someont is leaving")
			}
		}
	}
}

func (n *node) handleCommand(w http.ResponseWriter, r *http.Request) {
	req := CommandRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Cmd {
	case "set":
		n.handleSet(req, w)

	case "get":
		n.handleGet(req, w)

	case "del":
		n.handleDel(req, w)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (n *node) handleSet(req CommandRequest, w http.ResponseWriter) {
	err := n.store.Set(req.Key, req.Val)
	if err != nil {

		if errors.Is(err, &NotLeaderError{}) {
			log.Fatal("error not leader")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (n *node) handleGet(req CommandRequest, w http.ResponseWriter) {
	res, err := n.store.Get(req.Key)
	if err != nil {

		if errors.Is(err, &NotLeaderError{}) {
			log.Fatal("error not leader")
			return
		}

		if errors.Is(err, ErrKeyDoesNotExist) {
			log.Println("does not exist")
			return
		}

		log.Fatal(err)
	}

	data, err := json.Marshal(&CommandResponse{
		Value: res,
	})
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (n *node) handleDel(req CommandRequest, w http.ResponseWriter) {
	err := n.store.Del(req.Key)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

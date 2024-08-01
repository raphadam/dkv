package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/raphadam/dkv"
)

type CommandRequest struct {
	Cmd string `json:"cmd"`
	Key string `json:"key"`
	Val string `json:"val"`
}

type CommandResponse struct {
}

type Rest struct {
	store *dkv.DKV
}

func Serve(addr string, store *dkv.DKV) error {
	rs := Rest{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", rs.Handle)
	return http.ListenAndServe(addr, mux)
}

func (rs *Rest) Handle(w http.ResponseWriter, r *http.Request) {
	req := CommandRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Fatal("unable to make request")
	}

	switch req.Cmd {
	case "set":
		err = rs.store.Set(req.Key, req.Val)
		if err != nil {
			log.Fatal(err)
		}
	// case "get":
	// 	err = rs.store.Get(req.Key)
	// case "del":

	default:
		log.Fatal("not handled case")
	}
}

/*


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

httpAddr
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




*/

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

// func (n *node) HttpCommand(w http.ResponseWriter, r *http.Request) {
// 	cmd := CommandRequest{}

// 	err := json.NewDecoder(r.Body).Decode(&cmd)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	log.Println(cmd)

// 	switch cmd.Cmd {
// 	case "set":
// 		n.Set(cmd.Key, cmd.Val)
// 	case "get":
// 		log.Println("get")
// 	case "del":
// 		log.Println("del")
// 	}
// }

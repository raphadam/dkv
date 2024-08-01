package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/raphadam/dkv"
)

type JoinRequest struct {
	NodeID string `json:"nodeID"`
	Addr   string `json:"addr"`
}

type JoinResponse struct {
	Error string `json:"error"`
}

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

type Rest struct {
	store *dkv.DKV
}

func Serve(addr string, store *dkv.DKV) error {
	rs := Rest{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", rs.HandleCommand)
	mux.HandleFunc("POST /join", rs.HandleJoin)
	return http.ListenAndServe(addr, mux)
}

func (rs *Rest) HandleJoin(w http.ResponseWriter, r *http.Request) {
	req := JoinRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Fatal("unable to decode request")
	}

	err = rs.store.Join(req.NodeID, req.Addr)
	if err != nil {
		log.Fatal("unable to decode request")
	}
}

func (rs *Rest) HandleCommand(w http.ResponseWriter, r *http.Request) {
	req := CommandRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Cmd {
	case "set":
		rs.handleSet(req, w)

	case "get":
		rs.handleGet(req, w)

	case "del":
		rs.handleDel(req, w)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (rs *Rest) handleSet(req CommandRequest, w http.ResponseWriter) {
	err := rs.store.Set(req.Key, req.Val)
	if err != nil {

		if errors.Is(err, &dkv.NotLeaderError{}) {
			log.Fatal("error not leader")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rs *Rest) handleGet(req CommandRequest, w http.ResponseWriter) {
	res, err := rs.store.Get(req.Key)
	if err != nil {

		if errors.Is(err, &dkv.NotLeaderError{}) {
			log.Fatal("error not leader")
			return
		}

		if errors.Is(err, dkv.KEY_DOES_NOT_EXIST) {
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

func (rs *Rest) handleDel(req CommandRequest, w http.ResponseWriter) {
	err := rs.store.Del(req.Key)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
}

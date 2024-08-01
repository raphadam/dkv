package main

import (
	"flag"
	"log"

	"github.com/raphadam/dkv"
	"github.com/raphadam/dkv/rest"
)

var joinAddr string

func init() {
	flag.StringVar(&joinAddr, "join", "", "set to join node")
}

func main() {
	flag.Parse()

	store, err := dkv.New(true, "localhost:50001")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(rest.Serve(":40001", store))
}

// func main() {
// 	config := raft.DefaultConfig()

// 	raft.ServerID

// 	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:12345")
// 	if err != nil {
// 		log.Fatal("unable to resolve tcp addr")
// 	}

// 	transport, err := raft.NewTCPTransport("127.0.0.1:12345", addr, 3, 5*time.Second, os.Stderr)
// 	if err != nil {
// 		log.Fatal("unable to create transport", err)
// 	}

// 	kv := dkv.New()

// 	log.Printf("config local id %#v", config)

// 	node, err := raft.NewRaft(
// 		config,
// 		kv,
// 		raft.NewInmemStore(),
// 		raft.NewInmemStore(),
// 		raft.NewDiscardSnapshotStore(),
// 		transport,
// 	)
// 	if err != nil {
// 		log.Fatal("unable to create node", err)
// 	}

// 	node.BootstrapCluster(raft.Configuration{
// 		Servers: []raft.Server{
// 			{
// 				ID:      config.LocalID,
// 				Address: transport.LocalAddr(),
// 			},
// 		},
// 	})

// 	data, err := json.Marshal(&dkv.Command{Key: "mynameis", Val: "slimshady"})
// 	if err != nil {
// 		log.Fatal("unable to marshal command", err)
// 	}

// 	future := node.Apply(data, 5*time.Second)
// 	if err := future.Error(); err != nil {
// 		log.Fatal("error trying to apply", err)
// 	}

// 	log.Printf("Set command applied, current state: %v", kv)
// }

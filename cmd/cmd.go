package main

import (
	"flag"
	"log"
	"time"

	"github.com/raphadam/dkv"
)

var joinAddr string

func init() {
	flag.StringVar(&joinAddr, "join", "", "set to join node")
}

func main() {
	// flag.Parse()
	go func() {
		err := dkv.Serve(true, "127.0.0.1:30001", "127.0.0.1:40001", "127.0.0.1:50001", []string{})
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)

	go func() {
		err := dkv.Serve(false, "127.0.0.1:30002", "127.0.0.1:40002", "127.0.0.1:50002", []string{
			"127.0.0.1:40001",
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("TAHT")
	time.Sleep(1 * time.Hour)
}

// err = AskJoin("127.0.0.1:50002", "127.0.0.1:40001")
// if err != nil {
// 	log.Fatal(err)
// }

// log.Fatal(rest.Serve(":40002", store))

// func AskJoin(me string, httpOther string) error {
// 	// req := rest.JoinRequest{
// 	// 	NodeID: me,
// 	// 	Addr:   me,
// 	// }

// 	// b, err := json.Marshal(req)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// resp, err := http.Post(fmt.Sprintf("http://%s/join", httpOther), "application-type/json", bytes.NewReader(b))
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// resp.Body.Close()
// 	return nil
// }

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

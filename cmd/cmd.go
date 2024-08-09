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
	time.Sleep(2 * time.Second)

	go func() {
		err := dkv.Serve(false, "127.0.0.1:30003", "127.0.0.1:40003", "127.0.0.1:50003", []string{
			"127.0.0.1:40001",
		})
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)

	go func() {
		err := dkv.Serve(false, "127.0.0.1:30004", "127.0.0.1:40004", "127.0.0.1:50004", []string{
			"127.0.0.1:40003",
		})
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)

	go func() {
		err := dkv.Serve(false, "127.0.0.1:30005", "127.0.0.1:40005", "127.0.0.1:50005", []string{
			"127.0.0.1:40002",
		})
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)

	go func() {
		err := dkv.Serve(false, "127.0.0.1:30006", "127.0.0.1:40006", "127.0.0.1:50006", []string{
			"127.0.0.1:40005",
		})
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Hour)
}

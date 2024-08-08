package dkv

import (
	"log"
	"net"
	"sync"

	"github.com/hashicorp/serf/serf"
)

type Node struct {
	serf     *serf.Serf
	events   chan serf.Event
	members  map[string]serf.Member
	mu       sync.RWMutex
	bindAddr string
}

func NewNode(bindAddr string, tags map[string]string, bootstrap []string) (*Node, error) {
	addr, err := net.ResolveTCPAddr("tcp", bindAddr)
	if err != nil {
		return nil, err
	}

	n := &Node{
		members:  make(map[string]serf.Member),
		bindAddr: bindAddr,
	}
	n.events = make(chan serf.Event)

	config := serf.DefaultConfig()
	config.Init()

	config.MemberlistConfig.BindAddr = addr.IP.String()
	config.MemberlistConfig.BindPort = addr.Port
	config.EventCh = n.events
	config.Tags = tags
	config.NodeName = bindAddr

	n.serf, err = serf.Create(config)
	if err != nil {
		return nil, err
	}

	go n.eventHandler()

	if bootstrap != nil {
		i, err := n.serf.Join(bootstrap, true)
		log.Println("joined: ", i, "error:", err)
	}

	return n, nil
}

func (n *Node) eventHandler() {

	for e := range n.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if n.serf.LocalMember().Name == member.Name {
					continue
				}

				log.Printf("[%s]: %d joined", n.bindAddr, member.Port)

			}

		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if n.serf.LocalMember().Name == member.Name {
					continue
				}

				log.Printf("[%s]: %d left", n.bindAddr, member.Port)
			}

		}
	}
}

func (n *Node) Leave() error {
	return n.serf.Leave()
}

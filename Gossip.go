package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
)

type GossipProtocol struct {
	list          *memberlist.Memberlist
	config        *memberlist.Config
	blockchain    *Blockchain
	myAddress     int
	shareNumPeers int
	maxHops       int
	broadcasts    *memberlist.TransmitLimitedQueue
	mtx           sync.RWMutex
}

var networkShareRequests = make(map[int]bool)

func (g *GossipProtocol) SetupGossip(bootstrapUrl ...string) {
	config := memberlist.DefaultLocalConfig()
	config.Name = uuid.New().String()
	config.BindPort = g.blockchain.serverConfig.Port
	config.Events = &eventDelegate{}
	config.Delegate = &delegate{
		gossip: g,
	}

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	_, err = list.Join(bootstrapUrl)
	if err != nil {
		if !strings.Contains(err.(error).Error(), "Failed to join") {
			panic("Failed to join cluster: " + err.Error())
		}
	}

	g.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return list.NumMembers()
		},
		RetransmitMult: 3,
	}

	g.config = config
	g.list = list

	fmt.Printf("Local member %s:%d\n", list.LocalNode().Addr, list.LocalNode().Port)
}

type WorkPacket struct {
	Data  []byte `json:"data"`
	Block int    `json:"block"`
	Hops  int    `json:"hops"`
}

func (g *GossipProtocol) CascadeShare(data *WorkPacket) {
	log.Println("Cascading share", data.Block)
	if networkShareRequests[data.Block] {
		log.Println("Already requested", data.Block)
		return
	}

	if data.Hops >= g.maxHops {
		log.Println("Max hops reached", data.Block)
		return
	}

	data.Hops += 1

	members := g.list.Members()

	serializedData, err := json.Marshal(*data)
	if err != nil {
		log.Println("Failed to marshal data: " + err.Error())
		return
	}

	fmt.Println("There are", len(members), "members")
	for i := 0; i < g.shareNumPeers; i++ {
		a := rand.Intn(len(members))
		if members[a].FullAddress().Addr == g.list.LocalNode().FullAddress().Addr {
			log.Println("Skipping local node")
			continue
		}
		log.Println("Sending to " + members[a].Name)

		address := strings.Split(members[a].Address(), ":")[0] + ":4444/gossip"
		_, err := SendRequest(address, bytes.NewBuffer(serializedData))
		if err != nil {
			log.Println("Failed to send request: " + err.Error())
			continue
		}
	}

	networkShareRequests[data.Block] = true
}

// broadcasts.QueueBroadcast(&broadcast{
// 	msg:    append([]byte("d"), b...),
// 	notify: nil,
// })

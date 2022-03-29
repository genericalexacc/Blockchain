package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
)

type GossipProtocol struct {
	list          *memberlist.Memberlist
	config        *memberlist.Config
	blockchain    *Blockchain
	shareNumPeers int
	maxHops       int
	broadcasts    *memberlist.TransmitLimitedQueue
}

var networkShareRequests = make(map[int]bool)

func (g *GossipProtocol) SetupGossip(bootstrapUrl ...string) {
	config := memberlist.DefaultLocalConfig()
	config.Name = uuid.New().String()
	config.BindPort = g.blockchain.serverConfig.Port
	config.Events = &eventDelegate{}
	config.Delegate = &delegate{
		blockchain: g.blockchain,
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
	Node Node `json:"node"`
	Hops int  `json:"hops"`
}

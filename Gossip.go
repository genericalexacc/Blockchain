package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"

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
}

var networkShareRequests = make(map[int]bool)

func (g *GossipProtocol) SetupGossip(bootstrapUrl string) {
	config := memberlist.DefaultLocalConfig()
	config.Name = uuid.New().String()
	config.BindPort = g.myAddress

	list, err := memberlist.Create(config)
	if err != nil {
		panic("Failed to create memberlist: " + err.Error())
	}

	_, err = list.Join([]string{bootstrapUrl})
	if err != nil {
		if !strings.Contains(err.(error).Error(), "Failed to join") {
			panic("Failed to join cluster: " + err.Error())
		}
	}

	g.config = config
	g.list = list
}

type WorkPacket struct {
	Data  []byte `json:"data"`
	Block int    `json:"block"`
	Hops  int    `json:"hops"`
}

func (g *GossipProtocol) handleGossipedBlock(bodyAsBytes []byte) error {
	var wp WorkPacket

	err := json.Unmarshal(bodyAsBytes, &wp)
	if err != nil {
		log.Println("Failed to unmarshal work packet: " + err.Error())
		return err
	}

	err = g.blockchain.AddHash(wp.Data, wp.Block)
	if err != nil {
		log.Println("Failed to add block: " + err.Error())
		return err
	}

	go g.CascadeShare(&wp)

	return nil
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

	serializedData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal data: " + err.Error())
		return
	}

	for i := 0; i < g.shareNumPeers; i++ {
		a := rand.Intn(len(members))

		log.Println("Sending to " + members[a].Name)
		err := g.list.SendReliable(members[a], serializedData)
		log.Println(err)
	}

	networkShareRequests[data.Block] = true
}

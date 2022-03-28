package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func startOtherNode() {
	go func() {
		genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")
		gp2 := GossipProtocol{myAddress: 8081}
		gp2.SetupGossip("0.0.0.0:" + strconv.Itoa(globalConfig.Port))

		blockChainB := Blockchain{
			start:          genesisNode,
			lastBlock:      genesisNode,
			difficulty:     2,
			gossipProtocol: &gp2,
		}

		blockChainB.AddBlock(CreateNodeFromText("Send 0.1 BTC to Alice from Bob"))
		blockChainB.AddBlock(CreateNodeFromText("Send 0.1 BTC to Charles for Alice"))

		gp2.ListenForRequests()

		<-make(chan bool)
	}()
}

func TestGossip(t *testing.T) {
	genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")

	gp := GossipProtocol{
		myAddress: globalConfig.Port,
	}
	gp.SetupGossip("0.0.0.0:8081")

	blockChain := Blockchain{
		start:          genesisNode,
		lastBlock:      genesisNode,
		difficulty:     2,
		gossipProtocol: &gp,
	}

	go startOtherNode()

	fmt.Println(gp.list.Join([]string{"0.0.0.0:8080"}))
	go gp.ListenForRequests()

	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Alice from Bob"))
	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Charles for Alice"))
	time.Sleep(time.Second * 3)
	blockChain.Work()

	<-make(chan bool)
}

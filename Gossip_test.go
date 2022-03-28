package main

import (
	"fmt"
	"testing"
	"time"
)

func startOtherNode() {
	go func() {
		genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")
		gp := GossipProtocol{
			myAddress:     4443,
			maxHops:       2,
			shareNumPeers: 2,
		}

		blockChain := Blockchain{
			start:          genesisNode,
			lastBlock:      genesisNode,
			difficulty:     2,
			gossipProtocol: &gp,
			serverConfig: Config{
				Port:     8081,
				HttpPort: 4443,
			},
		}
		gp.blockchain = &blockChain
		gp.SetupGossip("0.0.0.0:8080")

		blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Alice from Bob"))
		blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Charles from Alice"))
		fmt.Println(gp.list.Join([]string{"0.0.0.0:8080"}))

		go func() {
			time.Sleep(time.Second * 10)
			blockChain.PrintAll()
		}()

		gp.ListenForRequests()

		<-make(chan bool)
	}()
}

func TestGossip(t *testing.T) {
	gp := GossipProtocol{
		myAddress:     4444,
		maxHops:       2,
		shareNumPeers: 2,
	}

	genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")

	blockChain := Blockchain{
		start:          genesisNode,
		lastBlock:      genesisNode,
		difficulty:     2,
		gossipProtocol: &gp,
		serverConfig: Config{
			Port:     8080,
			HttpPort: 4444,
		},
	}

	gp.blockchain = &blockChain
	gp.SetupGossip()
	go startOtherNode()
	go gp.ListenForRequests()

	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Alice from Bob"))
	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Charles from Alice"))
	time.Sleep(time.Second * 3)
	blockChain.Work()

	<-make(chan bool)
}

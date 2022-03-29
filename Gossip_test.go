package main

import (
	"fmt"
	"testing"
	"time"
)

func startOtherNode() {
	go func() {
		genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")

		blockChain := Blockchain{
			start:      genesisNode,
			lastBlock:  genesisNode,
			difficulty: 2,
			gossipProtocol: &GossipProtocol{
				maxHops:       2,
				shareNumPeers: 2,
			},
			serverConfig: Config{
				Port: 8081,
			},
			length: 1,
		}

		blockChain.gossipProtocol.blockchain = &blockChain
		blockChain.gossipProtocol.SetupGossip("0.0.0.0:8080")
		fmt.Println(blockChain.gossipProtocol.list.Join([]string{"0.0.0.0:8080"}))

		go func() {
			for {
				fmt.Println("There are", len(blockChain.gossipProtocol.list.Members()), "peers.")
				fmt.Println("I have", blockChain.length, "blocks.")
				time.Sleep(time.Second * 1)
			}
		}()

		<-make(chan bool)
	}()
}

func TestGossip(t *testing.T) {
	genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")

	blockChain := Blockchain{
		start:      genesisNode,
		lastBlock:  genesisNode,
		difficulty: 2,
		gossipProtocol: &GossipProtocol{
			maxHops:       2,
			shareNumPeers: 2,
		},
		serverConfig: Config{
			Port: 8080,
		},
		length: 1,
	}

	blockChain.gossipProtocol.blockchain = &blockChain
	blockChain.gossipProtocol.SetupGossip()

	go startOtherNode()

	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Alice from Bob"), true)
	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Charles from Alice"), true)
	blockChain.Work()

	<-make(chan bool)
}

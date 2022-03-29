package main

import (
	"time"
)

func main() {
	initString := []byte("Genesis Block")
	var arr [1024]byte
	copy(arr[:], initString[:])

	genesisNode := &Node{
		Val:        arr,
		Next:       nil,
		ParentHash: getHash([]byte("GENESIS")),
	}

	blockChain := Blockchain{
		start:      genesisNode,
		lastBlock:  genesisNode,
		difficulty: 2,
	}

	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Alice from Bob"), true)
	blockChain.AddBlock(CreateNodeFromText("Send 0.1 BTC to Charles for Alice"), true)

	a := time.Now()
	blockChain.Work()
	blockChain.PrintAll()

	println("Time taken:", time.Since(a).Milliseconds(), "ms")
}

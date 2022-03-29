package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const TestingDifficulty = 3

func testWork(t *testing.T) {
	genesisNode := CreateNodeFromText("GENESIS TRANSACTION.")
	genesisNode.ParentHash = getHash([]byte("GENESIS"))

	blockChain := Blockchain{
		start:          genesisNode,
		lastBlock:      genesisNode,
		difficulty:     TestingDifficulty,
		gossipProtocol: &GossipProtocol{},
	}

	for i := 0; i < 10; i++ {
		newRandomBytes := make([]byte, 128)
		rand.Read(newRandomBytes)
		randomText := string(newRandomBytes)
		blockChain.AddBlock(CreateNodeFromText(randomText), true)
	}

	a := time.Now()
	blockChain.Work()
	fmt.Println(time.Since(a).Milliseconds())

	currentBlock := blockChain.start

	for currentBlock != nil && currentBlock.Next != nil {
		hashable := []byte{}
		hashable = append(hashable, currentBlock.Next.Val[:]...)
		hashable = append(hashable, currentBlock.ParentHash...)
		hashable = append(hashable, currentBlock.Next.ParentHash...)

		println(hex.EncodeToString(getHash(hashable)))

		if hex.EncodeToString(getHash(hashable))[:TestingDifficulty] != strings.Repeat("0", TestingDifficulty) {
			t.Error("Hash is not valid")
		}

		currentBlock = currentBlock.Next
	}

}

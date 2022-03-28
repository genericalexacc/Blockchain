package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

type Blockchain struct {
	start          *Node
	lastBlock      *Node
	difficulty     int
	gossipProtocol *GossipProtocol
}

func (b *Blockchain) PrintAll() {
	n := b.start
	for n != nil {
		println("TRANSACTION")
		println("Message: " + string(n.Val[:]))
		println("Hash: " + hex.EncodeToString([]byte(n.ParentHash)))
		println("")
		n = n.Next
	}
}

func (b *Blockchain) AddBlockFromString(data [1024]byte) {
	newBlock := &Node{
		Val:      data,
		Next:     nil,
		Previous: b.lastBlock,
	}

	b.lastBlock.Next = newBlock
	b.lastBlock = newBlock
}

func (b *Blockchain) AddBlock(block *Node) {
	block.Previous = b.lastBlock
	b.lastBlock.Next = block
	b.lastBlock = block
}

func (b *Blockchain) Work() {
	current := b.start
	for current.Next != nil {
		if current.Next.ParentHash != nil {
			current = current.Next
			continue
		}
		b.computeHash(current.Next, current.ParentHash)
		fmt.Println("Found new block")

		wp := WorkPacket{
			Data:  current.Next.ParentHash,
			Block: current.Next.NodeIndex,
		}

		b.gossipProtocol.CascadeShare(&wp)
		current = current.Next
	}
}

func (b *Blockchain) computeHash(n *Node, parentHash []byte) {
	initHashable := []byte{}
	initHashable = append(initHashable, n.Val[:]...)
	initHashable = append(initHashable, parentHash...)

	currentHashable := []byte{}
	copy(currentHashable, initHashable[:])

	for {
		newRandomBytes := make([]byte, 32)
		rand.Read(newRandomBytes)
		currentHashable = append(initHashable, newRandomBytes...)

		currentHash := hex.EncodeToString(getHash(currentHashable))
		if currentHash[:b.difficulty] == strings.Repeat("0", b.difficulty) {
			n.ParentHash = newRandomBytes
			break
		}
	}
}

func (b *Blockchain) AddHash(data []byte, blockIndex int) error {
	current := b.start
	for current != nil {
		if current.NodeIndex == blockIndex {
			hashable := []byte{}
			hashable = append(hashable, current.Val[:]...)
			hashable = append(hashable, current.Previous.ParentHash...)
			hashable = append(hashable, current.ParentHash...)

			if hex.EncodeToString(getHash(hashable))[:b.difficulty] != strings.Repeat("0", b.difficulty) {
				return errors.New("Hash is not valid")
			}

			current.ParentHash = data
		}
		current = current.Next
		return nil
	}

	return errors.New("Block not found")
}

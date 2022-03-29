package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
)

type Blockchain struct {
	start          *Node
	lastBlock      *Node
	difficulty     int
	gossipProtocol *GossipProtocol
	serverConfig   Config
	length         int
	mtx            sync.RWMutex
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

	b.length++
	b.lastBlock.Next = newBlock
	b.lastBlock = newBlock
}

func (b *Blockchain) GetBlockById(id int) *Node {
	n := b.start
	for n != nil {
		if n.NodeIndex == id {
			return n
		}
		n = n.Next
	}
	return nil
}

func (b *Blockchain) AddBlock(block *Node, shouldBroadcast bool) {
	b.length++
	block.NodeIndex = b.lastBlock.NodeIndex + 1
	block.Previous = b.lastBlock
	b.lastBlock.Next = block
	b.lastBlock = block

	nodeCopy := *block
	nodeCopy.Next = nil
	nodeCopy.Previous = nil

	wp := WorkPacket{
		Node: nodeCopy,
		Hops: 0,
	}

	out, err := json.Marshal(wp)
	if err != nil {
		log.Fatal(err)
	}

	if shouldBroadcast {
		fmt.Println("Broadcasting new block")
		b.gossipProtocol.broadcasts.QueueBroadcast(&broadcast{
			msg:    append([]byte("a"), out...),
			notify: nil,
		})
	}
}

func (b *Blockchain) Work() {
	current := b.start
	for current.Next != nil {
		b.computeHash(current.Next, current.ParentHash)
		fmt.Println("Found new block")
		fmt.Println(string(current.Next.Val[:]))

		nodeCopy := *current.Next
		nodeCopy.Next = nil
		nodeCopy.Previous = nil

		wp := WorkPacket{
			Node: nodeCopy,
			Hops: 0,
		}

		out, err := json.Marshal(wp)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Sending to all peers")
		b.gossipProtocol.broadcasts.QueueBroadcast(&broadcast{
			msg:    append([]byte("d"), out...),
			notify: nil,
		})

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

func (b *Blockchain) AddHash(node Node) error {
	current := b.start
	for current != nil {
		if current.NodeIndex == node.NodeIndex {
			hashable := []byte{}
			hashable = append(hashable, current.Val[:]...)
			hashable = append(hashable, current.Previous.ParentHash...)
			hashable = append(hashable, current.ParentHash...)

			if hex.EncodeToString(getHash(hashable))[:b.difficulty] != strings.Repeat("0", b.difficulty) {
				log.Println("Hash is not valid")
				return errors.New("Hash is not valid")
			}

			current.ParentHash = node.ParentHash
		}
		current = current.Next
		return nil
	}

	return errors.New("Block not found")
}

func (b *Blockchain) Pack() []Node {
	n := b.start
	var nodes []Node
	for n != nil {
		nodes = append(nodes, *n)
		n = n.Next
	}
	return nodes
}

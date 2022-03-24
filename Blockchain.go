package main

import (
	"encoding/hex"
	"math/rand"
	"strings"
)

type Blockchain struct {
	start      *Node
	lastBlock  *Node
	difficulty int
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
		Val:  data,
		Next: nil,
	}

	b.lastBlock.Next = newBlock
	b.lastBlock = newBlock
}

func (b *Blockchain) AddBlock(block *Node) {
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
		newRandomBytes := make([]byte, 1024)
		rand.Read(newRandomBytes)
		currentHashable = append(initHashable, newRandomBytes...)

		currentHash := hex.EncodeToString(getHash(currentHashable))
		if currentHash[:b.difficulty] == strings.Repeat("0", b.difficulty) {
			n.ParentHash = newRandomBytes
			break
		}
	}
}

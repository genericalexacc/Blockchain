package main

import (
	"encoding/hex"
	"fmt"
)

type Node struct {
	NodeIndex  int
	Val        [1024]byte
	Next       *Node
	Previous   *Node
	ParentHash []byte
}

func CreateNodeFromText(data string) *Node {
	var arr [1024]byte
	copy(arr[:], data)
	return &Node{
		Val:        arr,
		Next:       nil,
		Previous:   nil,
		ParentHash: nil,
	}
}

func (n *Node) Print() {
	fmt.Println("TRANSACTION")
	fmt.Println("Message: " + string(n.Val[:]))
	fmt.Println("Hash: " + hex.EncodeToString([]byte(n.ParentHash)))
	fmt.Println("Previous:", n == nil)
	fmt.Println("Next:", &n.Next)
}

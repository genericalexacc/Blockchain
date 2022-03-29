package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type Node struct {
	NodeIndex  int        `json:"node_index"`
	Val        [1024]byte `json:"val"`
	Next       *Node      `json:"next"`
	Previous   *Node      `json:"previous"`
	ParentHash []byte     `json:"parent_hash"`
}

func (u *Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		NodeIndex  int        `json:"node_index"`
		Val        [1024]byte `json:"val"`
		ParentHash []byte     `json:"parent_hash"`
	}{
		NodeIndex:  u.NodeIndex,
		Val:        u.Val,
		ParentHash: u.ParentHash,
	})
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

func (n *Node) PrintAll() {
	current := n
	for current != nil {
		current.Print()
		current = current.Next
	}
}

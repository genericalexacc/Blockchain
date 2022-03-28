package main

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

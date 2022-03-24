package main

type Node struct {
	Val        [1024]byte
	Next       *Node
	ParentHash []byte
}

func CreateNodeFromText(data string) *Node {
	var arr [1024]byte
	copy(arr[:], data)
	return &Node{
		Val:        arr,
		Next:       nil,
		ParentHash: nil,
	}
}

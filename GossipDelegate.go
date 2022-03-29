package main

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/memberlist"
)

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

type delegate struct {
	blockchain *Blockchain
}

func init() {
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	fmt.Println("Received message: ", string(b))
	if len(b) == 0 {
		return
	}

	g := d.blockchain

	switch b[0] {
	case 'd':
		var update *WorkPacket
		fmt.Println("Received a data packet")
		if err := json.Unmarshal(b[1:], &update); err != nil {
			fmt.Println(err)
			return
		}

		g.mtx.Lock()
		g.AddHash(update.Node)
		g.mtx.Unlock()
		break
	case 'a':
		var update *WorkPacket
		fmt.Println("Received an addition data packet")
		if err := json.Unmarshal(b[1:], &update); err != nil {
			fmt.Println(err)
			return
		}

		g.mtx.Lock()
		g.AddBlock(&update.Node, false)
		g.mtx.Unlock()
		break
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return d.blockchain.gossipProtocol.broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	g := d.blockchain.gossipProtocol
	d.blockchain.mtx.RLock()
	m := g.blockchain.Pack()
	d.blockchain.mtx.RUnlock()
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error marshalling", err)
		return nil
	}
	return b
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}

	var m []Node
	if err := json.Unmarshal(buf, &m); err != nil {
		return
	}

	d.blockchain.mtx.Lock()
	defer d.blockchain.mtx.Unlock()

	fmt.Println("Merging remote state", len(m))
	for _, v := range m[1:] {
		fmt.Println("Adding node: ", v.NodeIndex, string(v.Val[:]))

		if d.blockchain.GetBlockById(v.NodeIndex) != nil {
			go d.blockchain.AddHash(v)
			continue
		}

		if v.NodeIndex > d.blockchain.length-2 {
			d.blockchain.AddBlock(&v, false)
		}
	}
}

type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Println("A node has joined: " + node.String())
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Println("A node has left: " + node.String())
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Println("A node was updated: " + node.String())
}

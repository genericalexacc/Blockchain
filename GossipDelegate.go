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
	gossip *GossipProtocol
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
	if len(b) == 0 {
		return
	}

	g := d.gossip

	switch b[0] {
	case 'd':
		var update *WorkPacket
		fmt.Println("Received a data packet")
		if err := json.Unmarshal(b[1:], &update); err != nil {
			return
		}

		g.mtx.Lock()
		g.blockchain.AddHash(update.Data, update.Block)
		g.mtx.Unlock()
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	g := d.gossip
	return g.broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	g := d.gossip
	g.mtx.RLock()
	m := g.blockchain
	g.mtx.RUnlock()
	b, _ := json.Marshal(m)
	return b
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	g := d.gossip
	var m []Node
	if err := json.Unmarshal(buf, &m); err != nil {
		return
	}
	g.mtx.Lock()
	for _, v := range m {
		g.blockchain.AddBlock(&v)
	}
	g.mtx.Unlock()
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

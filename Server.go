package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (g *GossipProtocol) ListenForRequests() {
	node := g.list.LocalNode()

	r := gin.Default()

	r.GET("/gossip/:block", func(c *gin.Context) {
		block, err := strconv.Atoi(c.Param("block"))
		if err != nil {
			c.JSON(400, gin.H{"message": "Invalid block number"})
			return
		}

		if block < 0 || block > g.blockchain.lastBlock.NodeIndex {
			c.JSON(400, gin.H{"message": "Block number out of range"})
			return
		}

		current := g.blockchain.start
		for i := 0; i < g.blockchain.lastBlock.NodeIndex; i++ {
			if i == current.NodeIndex {
				c.JSON(200, gin.H{"message": "Block found", "block": g.blockchain.lastBlock.Val[:]})
				return
			}
			current = current.Next
		}

		c.JSON(404, gin.H{"message": "Block not found"})
	})

	r.GET("/gossip/:block/:batch", func(c *gin.Context) {
		block, err := strconv.Atoi(c.Param("block"))
		if err != nil {
			c.JSON(400, gin.H{"message": "Invalid block number"})
			return
		}

		if block < 0 || block > g.blockchain.lastBlock.NodeIndex {
			c.JSON(400, gin.H{"message": "Block number out of range"})
			return
		}

		batch, err := strconv.Atoi(c.Param("batch"))
		if err != nil {
			c.JSON(400, gin.H{"message": "Invalid batch number"})
			return
		}

		if batch > 100 {
			c.JSON(400, gin.H{"message": "Batch number too large"})
			return
		}

		blocksBatch := []Node{}

		current := g.blockchain.start
		found := 0

		for i := 0; i < g.blockchain.lastBlock.NodeIndex; i++ {
			if i == current.NodeIndex || found > 0 {
				blocksBatch = append(blocksBatch, *current)
				found++
			}
			if found == batch {
				break
			}
			current = current.Next
		}

		c.JSON(200, gin.H{"message": "Blocks found", "blocks": blocksBatch})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(node.Address())
}

package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// getBlock is a handler for GET requests to /gossip/:block
// It receives a block number and returns the block
func getBlock(c *gin.Context) {
	g := c.Value("GOSSIP").(*GossipProtocol)

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
}

// getBlockBatch is a handler for GET requests to /gossip/:block/:batch
// It receives a block number and a batch size and returns a batch of blocks
func getBlockBatch(c *gin.Context) {
	g := c.Value("GOSSIP").(*GossipProtocol)

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
	counter := 0

	for i := 0; i < g.blockchain.lastBlock.NodeIndex; i++ {
		if i == current.NodeIndex || counter > 0 {
			blocksBatch = append(blocksBatch, *current)
			counter++
		}
		if counter == batch {
			break
		}
		current = current.Next
	}

	c.JSON(200, gin.H{"message": "Blocks found", "blocks": blocksBatch})
}

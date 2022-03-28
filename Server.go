package main

import (
	"github.com/gin-gonic/gin"
)

func (g *GossipProtocol) ListenForRequests() {
	r := gin.Default()
	r.Use(func(ctx *gin.Context) { ctx.Set("GOSSIP", g) })

	r.GET("/gossip/:block", getBlock)
	r.GET("/gossip/:block/:batch", getBlockBatch)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	localNode := g.list.LocalNode()
	r.Run(localNode.Address())
}

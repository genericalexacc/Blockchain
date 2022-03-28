package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (g *GossipProtocol) ListenForRequests() {
	r := gin.Default()
	r.Use(func(ctx *gin.Context) { ctx.Set("GOSSIP", g) })

	r.GET("/gossip/:block", getBlock)
	r.GET("/gossip/:block/:batch", getBlockBatch)
	r.POST("/gossip", postBlock)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":" + strconv.Itoa(g.blockchain.serverConfig.HttpPort))
}

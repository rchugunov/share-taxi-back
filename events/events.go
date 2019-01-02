package events

import "github.com/gin-gonic/gin"

func TestEvent(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong1",
	})
}

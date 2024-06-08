package internal

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func exampleHandler(c *gin.Context) {
	keyStr := c.DefaultQuery("key", "value")
	timeStr := c.DefaultQuery("time", time.Now().Format(time.RFC3339Nano))

	timeObj, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("Invalid time format: %s", timeStr),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Hello, World!",
		"key":     keyStr,
		"time":    timeObj.Format(time.RFC3339Nano),
	})
}

func RunServer() {
	router := gin.Default()
	router.GET("/api/v1/example", exampleHandler)
	router.Run(":8080")
}

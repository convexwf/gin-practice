package route

import "github.com/gin-gonic/gin"

func SetupRouter(router *gin.Engine) {

	helloGroup := router.Group("/hello")
	{
		helloGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		helloGroup.POST("/echo", func(c *gin.Context) {
			requestBody, err := c.GetRawData()
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.Data(200, "application/json", requestBody)
		})
	}

	userGroup := router.Group("/user")
	{

	}

}

package roverserver

import (
	"github.com/gin-gonic/gin"
)

// Serve all routes and websocket of the rover server
func Serve() (err error) {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(200, "The rover is running its service at this address")
	})

	r.POST("/api/v1/sendcommand", sendcommand)

	go func(r *gin.Engine) {
		r.Run("0.0.0.0:80")
	}(r)

	return
}

func sendcommand(c *gin.Context) {
	cmd, _ := c.GetPostForm("command")
	c.String(200, "Sendcommand "+cmd+" executed")
}

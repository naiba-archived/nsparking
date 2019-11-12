package controller

import "github.com/gin-gonic/gin"

// ServeWeb ...
func ServeWeb() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {

	})
	r.Run(":80")
}

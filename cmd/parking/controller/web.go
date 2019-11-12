package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naiba/nsparking/model"
)

// ServeWeb ...
func ServeWeb() {
	r := gin.Default()
	r.LoadHTMLGlob("resource/template/*")
	r.Static("/static", "resource/static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GClient": model.GClient,
		})
	})
	r.POST("/up", up)
	r.Run(":80")
}

func up(c *gin.Context) {

}

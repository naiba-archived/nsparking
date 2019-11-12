package controller

import (
	"fmt"
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

type upReq struct {
	G  string `form:"g" binding:"required"`
	To string `from:"to" binding:"required;url"`
}

type upRsp struct {
	Success bool
	Msg     string
	Data    string
}

func up(c *gin.Context) {
	var ur upReq
	var up upRsp
	if err := c.ShouldBind(&ur); err != nil {
		up.Msg = fmt.Sprintf("输入数据有误：%s", err)
		c.JSON(http.StatusOK, up)
		return
	}

}

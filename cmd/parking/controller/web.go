package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naiba/com"

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
	To string `form:"to" binding:"url,required"`
}

type upRsp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg,omitempty"`
	Data    string `json:"data,omitempty"`
}

func up(c *gin.Context) {
	var ur upReq
	var up upRsp
	if err := c.ShouldBind(&ur); err != nil {
		up.Msg = fmt.Sprintf("输入数据有误：%s", err)
		c.JSON(http.StatusOK, up)
		return
	}

	if !captcha.Verify(ur.G, c.ClientIP()) {
		// up.Msg = fmt.Sprintf("人机验证未通过，请重试")
		// c.JSON(http.StatusOK, up)
		// return
	}

	var r model.Redirect
	r.IP = c.ClientIP()
	r.CreatedAt = time.Now()
	r.UserAgent = c.Request.Header.Get("User-Agent")
	r.To = ur.To
	r.Server = com.MD5(fmt.Sprintf("%s", r))[:10]
	if err := db.Create(r).Error; err != nil {
		up.Msg = fmt.Sprintf("数据库错误：%s", err)
		c.JSON(http.StatusOK, up)
		return
	}

	up.Success = true
	up.Data = r.Server
	c.JSON(http.StatusOK, up)
}

package controller

import (
	"fmt"
	"net/http"
	"strings"
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
	// 处理域名跳转
	r.Use(func(c *gin.Context) {
		if !strings.HasSuffix(c.Request.Host, model.Domain) {
			rd, err := getRedirectByDomain(c.Request.Host)
			if err != nil {
				c.String(http.StatusOK, "错误：未找到对应跳转")
				c.Abort()
				return
			}
			domain, prefix, suffix := parseDomain(c.Request.Host)
			rd.To = strings.ReplaceAll(rd.To, "[domain]", domain)
			rd.To = strings.ReplaceAll(rd.To, "[prefix]", prefix)
			rd.To = strings.ReplaceAll(rd.To, "[suffix]", suffix)
			c.Redirect(http.StatusFound, rd.To)
			c.Abort()
			return
		}
	})
	// 首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GClient": model.GClient,
		})
	})
	// 处理启动跳转
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
	up.Data = fmt.Sprintf("%s.%s", r.Server, model.Domain)
	c.JSON(http.StatusOK, up)
}

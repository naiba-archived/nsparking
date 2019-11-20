package controller

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/miekg/dns"
	"github.com/naiba/com"

	"github.com/naiba/nsparking/data"
	"github.com/naiba/nsparking/model"
)

// ServeWeb ...
func ServeWeb() {
	r := gin.Default()
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	var templatePrefix string
	if model.Domain == "localhost" {
		r.LoadHTMLGlob("resource/template/*")
		r.Static("/static", "resource/static")
	} else {
		templatePrefix = "/"
		r.SetHTMLTemplate(t)
		r.StaticFS("/static", data.StaticFS)
	}
	// 处理域名跳转
	r.Use(func(c *gin.Context) {
		if !strings.HasSuffix(c.Request.Host, model.Domain) {
			rd, err := getRedirectByDomain(c.Request.Host)
			if err != nil || rd.Mode != "url" {
				c.String(http.StatusOK, "错误：未找到对应跳转，您配置的是URL跳转模式吗？")
				c.Abort()
				return
			}
			domain, prefix, suffix := parseDomain(c.Request.Host)
			rd.Value = strings.ReplaceAll(rd.Value, "[domain]", domain)
			rd.Value = strings.ReplaceAll(rd.Value, "[prefix]", prefix)
			rd.Value = strings.ReplaceAll(rd.Value, "[suffix]", suffix)
			c.Redirect(http.StatusFound, rd.Value)
			c.Abort()
			return
		}
	})

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, templatePrefix+"index.html", gin.H{
			"GClient": model.GClient,
			"Domain":  model.Domain,
		})
	})
	// 处理启动跳转
	r.POST("/up", up)
	r.Run(":8080")
}

type upReq struct {
	G     string `binding:"required" json:"g,omitempty"`
	Mode  string `binding:"required" json:"mode,omitempty"`
	Value string `binding:"required" json:"value,omitempty"`

	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
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

	if ok := model.Modes[ur.Mode]; !ok {
		up.Msg = fmt.Sprintf("不支持的模式：%s", ur.Mode)
		c.JSON(http.StatusOK, up)
		return
	}

	switch ur.Mode {
	case "a":
		if !govalidator.IsIPv4(ur.Value) {
			up.Msg = fmt.Sprintf("A 记录不符合规范：%s", ur.Value)
			c.JSON(http.StatusOK, up)
			return
		}
	case "cname":
		ur.Value = dns.Fqdn(ur.Value)
		if !govalidator.IsDNSName(ur.Value) {
			up.Msg = fmt.Sprintf("主机记录不符合规范：%s", ur.Value)
			c.JSON(http.StatusOK, up)
			return
		}
	case "url":
		if !govalidator.IsURL(ur.Value) {
			up.Msg = fmt.Sprintf("URL 不符合规范：%s", ur.Value)
			c.JSON(http.StatusOK, up)
			return
		}
	}

	if !captcha.Verify(ur.G, c.ClientIP()) {
		up.Msg = fmt.Sprintf("人机验证未通过，请重试")
		c.JSON(http.StatusOK, up)
		return
	}

	var r model.Parking

	ur.Password = com.MD5(strings.TrimSpace(ur.Password))

	if ur.ID != "" {
		if ur.Password == "" {
			up.Msg = fmt.Sprintf("管理密码不能为空：%s", ur.Password)
			c.JSON(http.StatusOK, up)
			return
		}
		if err := db.Where("id = ? AND password = ?", ur.ID, ur.Password).First(&r).Error; err != nil {
			up.Msg = fmt.Sprintf("未找到该记录：%s", err)
			c.JSON(http.StatusOK, up)
			return
		}
	}

	r.IP = c.ClientIP()
	r.CreatedAt = time.Now()
	r.UserAgent = c.Request.Header.Get("User-Agent")
	r.Value = ur.Value
	r.Mode = ur.Mode
	r.Password = ur.Password

	if r.ID == "" {
		r.ID = com.MD5(fmt.Sprintf("%s", r))[:10]
	}

	if err := db.Save(r).Error; err != nil {
		up.Msg = fmt.Sprintf("数据库错误：%s", err)
		c.JSON(http.StatusOK, up)
		return
	}

	up.Success = true
	up.Data = fmt.Sprintf("%s.ns1.%s, %s.ns2.%s", r.ID, model.Domain, r.ID, model.Domain)
	c.JSON(http.StatusOK, up)
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range data.TemplateFS.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".html") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

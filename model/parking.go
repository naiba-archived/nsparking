package model

import "time"

// Modes ..
var Modes map[string]bool

func init() {
	Modes = make(map[string]bool)
	Modes["cname"] = true
	Modes["url"] = true
	Modes["a"] = true
}

// Parking ..
type Parking struct {
	ID       string `gorm:"type:char(10);PRIMARY_KEY"`
	Value    string //值 域名、IP、URL
	Mode     string `gorm:"type:char(5);IDNEX"` //跳转模式
	Password string //管理密码

	IP        string
	UserAgent string
	CreatedAt time.Time
}

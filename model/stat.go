package model

// Stat ..
type Stat struct {
	Date     string `gorm:"type:char(8)"`
	DNS      uint64 // DNS 请求次数
	Web      uint64 // 访问次数
	Redirect uint64 // 新增跳转
}

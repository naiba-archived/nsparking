package model

import "time"

// Redirect ..
type Redirect struct {
	Server    string `gorm:"type:char(10);PRIMARY_KEY"`
	To        string
	IP        string
	UserAgent string
	CreatedAt time.Time
}

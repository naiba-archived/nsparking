package model

import "time"

// Redirect ..
type Redirect struct {
	Domain    string
	To        string
	IP        string
	UserAgent string
	CreatedAt time.Time
}

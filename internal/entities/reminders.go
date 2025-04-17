package entities

import "time"

type reminder struct {
	Time    time.Time
	Title   string
	Content string
}

package entities

import "time"

type Reminder struct {
	Text  string
	Topic string
	Time  time.Time
}

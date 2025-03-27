package domain

import "time"

type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}

type Lesson struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Location  string
	Teacher   string
}

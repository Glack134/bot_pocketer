package domain

import "time"

// User - модель пользователя
type User struct {
	ID        int64
	Username  string
	CreatedAt time.Time
}

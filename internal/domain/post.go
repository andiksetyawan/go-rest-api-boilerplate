package domain

import "time"

type Post struct {
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package domain

import "time"

type Note struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"authorID"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

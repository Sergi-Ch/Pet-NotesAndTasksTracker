package domain

import (
	"errors"
	"time"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrInvalidTitle = errors.New("title must be at least 3 characters")
	ErrAccessDenied = errors.New("access denied") // Доступ запрещен
)

type Note struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"authorID"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

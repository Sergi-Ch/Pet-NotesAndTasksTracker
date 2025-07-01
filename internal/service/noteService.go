package service

import (
	"NotesAndTasks/internal/domain"
	"NotesAndTasks/internal/repository"
	"context"
	"time"
)

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(ctx context.Context, note *domain.Note) error {
	if len(note.Title) < 3 {
		return domain.ErrInvalidTitle
	}

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	return s.repo.Create(ctx, note)
}

func (s *NoteService) GetNoteById(ctx context.Context, id int) (*domain.Note, error) {
	return s.repo.GetById(ctx, id)
}

func (s *NoteService) GetAll(ctx context.Context, authorId int) ([]*domain.Note, error) {
	return s.repo.GetAll(ctx, authorId)
}

func (s *NoteService) Update(ctx context.Context, note *domain.Note) error {
	return s.repo.Update(ctx, note)
}

func (s *NoteService) Delete(ctx context.Context, id int, authorID int) error {
	return s.repo.Delete(ctx, id, authorID)
}

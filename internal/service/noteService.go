package service

import (
	"NotesAndTasks/internal/domain"
	"NotesAndTasks/internal/repository"
	"context"
	"time"
)

type NoteServ struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *NoteServ {
	return &NoteServ{repo: repo}
}

func (s *NoteServ) CreateNote(ctx context.Context, note *domain.Note) error {
	if len(note.Title) < 3 {
		return domain.ErrInvalidTitle
	}

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	return s.repo.Create(ctx, note)
}

func (s *NoteServ) GetNoteById(ctx context.Context, id int) (*domain.Note, error) {
	return s.repo.GetById(ctx, id)
}

func (s *NoteServ) GetAll(ctx context.Context) ([]*domain.Note, error) {
	return s.repo.GetAll(ctx)
}

func (s *NoteServ) Update(ctx context.Context, note *domain.Note) error {
	return s.repo.Update(ctx, note)
}

func (s *NoteServ) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

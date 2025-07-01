package repository

import (
	"NotesAndTasks/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type NoteRepository interface { //Нужно, чтобы потом дергать в main и дать возможность
	Create(ctx context.Context, note *domain.Note) error // Переписать методы под другую базу данных, но с таким же названием.
	GetById(ctx context.Context, id int) (*domain.Note, error)
	GetAll(ctx context.Context, authorID int) ([]*domain.Note, error)
	Update(ctx context.Context, note *domain.Note) error
	Delete(ctx context.Context, id int, authorID int) error
}

type NoteRepoPG struct { // Структура, заметки с полем типа подключения (в main мы подрубаемся и этот тип позволяет
	db *pgx.Conn // Работать с подключением. Либо conn либо пул соединений pgx.pool
}

func NewNotePG(db *pgx.Conn) *NoteRepoPG {
	return &NoteRepoPG{db: db}
}

func (r *NoteRepoPG) Create(ctx context.Context, note *domain.Note) error {

	query := `INSERT INTO notes (title, content, authorID) 
          VALUES ($1, $2, $3) 
          RETURNING id, createdAT, updatedAt` // обратная кавычка это для многострочных строк
	err := r.db.QueryRow(context.Background(), query, note.Title, note.Content, note.AuthorID).Scan( //почему мы не передаем время, мы ведь его в сервисе задаем.
		&note.ID,
		&note.CreatedAt,
		&note.UpdatedAt)
	return err

	//знаки $ это Placeholders  для безопасной работы, под них подставляются аргументы, которые идут после запроса в
	// queryRow (note.Title, note.Content, note.AuthorId.
	//Scan считывает return запроса в адрес полей структуры.
	// return просто экономит время, мы могли вызвать select чтобы получить данные только что созданной заметки.
}

func (r *NoteRepoPG) GetByID(ctx context.Context, id int) (**domain.Note, error) {
	query := `SELECT id, title, content, authorID, createdAt, updatedAt
		FROM notes
		WHERE id = $1`
	var note *domain.Note
	err := r.db.QueryRow(ctx, query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.AuthorID,
		&note.CreatedAt,
		&note.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}
	return &note, err
}

func (r *NoteRepoPG) GetAll(ctx context.Context, authorID int) ([]*domain.Note, error) {
	query := `SELECT id, title, content, authorID, createdAt, updatedAt FROM notes
        WHERE authorID = $1
		ORDER BY createdAt DESC `

	rows, err := r.db.Query(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer func() { rows.Close() }()

	notes := make([]*domain.Note, 0)

	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.AuthorID,
			&note.CreatedAt,
			&note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error %w", err)
	}
	return notes, nil
}

func (r *NoteRepoPG) Update(ctx context.Context, note *domain.Note) error {
	query := `UPDATE notes
           SET title = $2, content = $3, updatedAt = $4
           WHERE id = $1 and authorID = $5
           RETURNING updatedAt` // че за запрос такой веселый
	now := time.Now()
	if err := r.db.QueryRow(ctx, query,
		note.ID,
		note.Title,
		note.Content,
		now,
		note.AuthorID).Scan(&note.UpdatedAt); err != nil {
		return fmt.Errorf("error of scan %w", err)
	}
	return nil
}

func (r *NoteRepoPG) Delete(ctx context.Context, id int, authorID int) error {
	query := `DELETE FROM notes
WHERE id = $1 and authorID = $2`
	result, err := r.db.Exec(ctx, query, id, authorID) //что такое Exec?
	if err != nil {
		return fmt.Errorf("error of delete %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("the note didn't delete %w", err)
	}

	return nil
}

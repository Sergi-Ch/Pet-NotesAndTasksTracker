package main

import (
	"NotesAndTasks/internal/handler"
	"NotesAndTasks/internal/repository"
	"NotesAndTasks/internal/service"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
)

func main() {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	NoteRepo := repository.NewNotePG(dbPool)
	NoteService := service.NewNoteService(NoteRepo)
	NoteHandler := handler.NewNoteHandler(NoteService)

	r := chi.NewRouter()
	r.Post("/notes", NoteHandler.CreateNote)
	r.Get("/notes/{id}", NoteHandler.GetNote)

	port := "8080"
	log.Printf("Starting server on port>> ", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

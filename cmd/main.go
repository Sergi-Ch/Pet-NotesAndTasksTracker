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
	"time"
)

func main() {
	//err := godotenv.Load() // очень важно его в начале подгрузить
	//if err != nil {
	//	log.Printf("err of load env file>> %v\n", err)
	//}

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Явная проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Successfully connected to database!")

	defer dbPool.Close()

	NoteRepo := repository.NewNotePG(dbPool)
	NoteService := service.NewNoteService(NoteRepo)
	NoteHandler := handler.NewNoteHandler(NoteService)

	r := chi.NewRouter()
	r.Post("/notes", NoteHandler.CreateNote)
	r.Get("/notes/{id}", NoteHandler.GetNote)
	r.Get("/notes/all", NoteHandler.GetAll)
	r.Delete("/notes/delete/{id}", NoteHandler.Delete)
	port := ":8080"
	log.Printf("Starting server on port >> %s", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

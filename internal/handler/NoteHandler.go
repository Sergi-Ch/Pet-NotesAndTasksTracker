package handler

import (
	"NotesAndTasks/internal/domain"
	"NotesAndTasks/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type NoteHandler struct {
	service *service.NoteServ
}

func NewNoteHandler(s *service.NoteServ) *NoteHandler {
	return &NoteHandler{service: s}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note domain.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		switch err {
		case domain.ErrInvalidTitle:
			http.Error(w, "Title too short", http.StatusBadRequest)
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, "Error of encoding", http.StatusInternalServerError)
	}
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id") // что за urlParam?
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	note, err := h.service.GetNoteById(r.Context(), id)
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(note)
	if err != nil {
		http.Error(w, "Error of encoding", http.StatusInternalServerError)
	}
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, "Error of encoding", http.StatusInternalServerError)
	}

}

func (h *NoteHandler) UdateNote(w http.ResponseWriter, r http.Request) { //вообще в других слоях на вход идет указатель на заметку, а сейчас ии предлагает по id

}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}
	err = h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "error of delete", http.StatusNotFound)
	}
}

package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/punnch/go-todo/internal/core/apperrors"
	"github.com/punnch/go-todo/internal/services/db-service/todo"
)

type Server struct {
	router  *mux.Router
	addr    string
	service *todo.TodoService
}

func NewServer(addr string, service *todo.TodoService) *Server {
	return &Server{
		router:  mux.NewRouter(),
		addr:    addr,
		service: service,
	}
}

func (s *Server) StartSever() error {
	s.router.HandleFunc("/tasks", s.createTask).Methods("POST")
	s.router.HandleFunc("/tasks", s.getAllTasks).Methods("GET")
	s.router.HandleFunc("/tasks/{id}", s.getTask).Methods("GET")
	s.router.HandleFunc("/tasks/{id}", s.deleteTask).Methods("DELETE")
	s.router.HandleFunc("/tasks/{id}", s.completeTask).Methods("PATCH")

	if err := http.ListenAndServe(s.addr, s.router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	// Create struct to get user's payload
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
		return
	}

	task, err := s.service.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, NewErrorDTO(err))
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (s *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	idStr := q.Get("id")
	completedStr := q.Get("completed")

	var id *int
	if idStr != "" {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
			return
		}

		id = &idInt
	}

	var completed *bool
	if completedStr != "" {
		completedBool, err := strconv.ParseBool(completedStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
			return
		}

		completed = &completedBool
	}

	tasks, err := s.service.GetAllTasks(r.Context(), id, completed)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, NewErrorDTO(err))
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
		return
	}

	task, err := s.service.GetTask(r.Context(), id)
	if err != nil {
		errorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
		return
	}

	if err := s.service.DeleteTask(r.Context(), id); err != nil {
		errorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Completed *bool `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err))
		return
	}

	task, err := s.service.CompleteTask(r.Context(), id, req.Completed)
	if err != nil {
		errorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

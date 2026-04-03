package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/punnch/go-todo/internal/core/apperrors"
	"github.com/punnch/go-todo/internal/services/db-service/todo"
	"go.uber.org/zap"
)

type Server struct {
	router  *mux.Router
	addr    string
	service *todo.TodoService
	log     *zap.Logger
}

func NewServer(addr string, service *todo.TodoService, log *zap.Logger) *Server {
	return &Server{
		router:  mux.NewRouter(),
		addr:    fmt.Sprintf(":%s", addr),
		service: service,
		log:     log,
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
		s.log.Warn("invalid request body", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
		return
	}

	task, err := s.service.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		s.log.Error("failed to create task", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, NewErrorDTO(err), s.log)
		return
	}

	s.log.Info("task created", zap.Int("id", task.ID))
	writeJSON(w, http.StatusCreated, task, s.log)
}

func (s *Server) getAllTasks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	idStr := q.Get("id")
	completedStr := q.Get("completed")

	var id *int
	if idStr != "" {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			s.log.Warn("ivalid query param 'id'", zap.Error(err))
			writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
			return
		}

		id = &idInt
	}

	var completed *bool
	if completedStr != "" {
		completedBool, err := strconv.ParseBool(completedStr)
		if err != nil {
			s.log.Warn("ivalid query param 'completed'", zap.Error(err))
			writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
			return
		}

		completed = &completedBool
	}

	tasks, err := s.service.GetAllTasks(r.Context(), id, completed)
	if err != nil {
		s.log.Error("failed to get all tasks", zap.Error(err))
		writeJSON(w, http.StatusInternalServerError, NewErrorDTO(err), s.log)
		return
	}

	s.log.Info("all tasks gotten", zap.Int("amount", len(tasks)))
	writeJSON(w, http.StatusOK, tasks, s.log)
}

func (s *Server) getTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Warn("invalid id", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
		return
	}

	task, err := s.service.GetTask(r.Context(), id)
	if err != nil {
		s.log.Error("failed to get task", zap.Error(err))
		errorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	s.log.Info("task gotten", zap.Int("id", task.ID))
	writeJSON(w, http.StatusOK, task, s.log)
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Warn("invalid id", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
		return
	}

	if err := s.service.DeleteTask(r.Context(), id); err != nil {
		s.log.Error("failed to delete task", zap.Error(err))
		errorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	s.log.Info("task deleted", zap.Int("id", id))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Completed *bool `json:"completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Warn("invalid request body", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.log.Warn("invalid id", zap.Error(err))
		writeJSON(w, http.StatusBadRequest, NewErrorDTO(err), s.log)
		return
	}

	task, err := s.service.CompleteTask(r.Context(), id, req.Completed)
	if err != nil {
		s.log.Error("failed to complete task", zap.Error(err))
		errorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	s.log.Info("task completed", zap.Int("id", id))
	writeJSON(w, http.StatusOK, task, s.log)
}

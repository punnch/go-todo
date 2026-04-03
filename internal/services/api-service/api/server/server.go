package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router   *mux.Router
	addr     string
	handlers *Handler
}

func NewServer(addr string, handler *Handler) *Server {
	return &Server{
		router:   mux.NewRouter(),
		addr:     fmt.Sprintf(":%s", addr),
		handlers: handler,
	}
}

func (s *Server) StartServer() error {
	s.router.HandleFunc("/tasks", s.handlers.CreateTask).Methods("POST")
	s.router.HandleFunc("/tasks", s.handlers.GetAllTasks).Methods("GET")
	s.router.HandleFunc("/tasks/{id}", s.handlers.GetTask).Methods("GET")
	s.router.HandleFunc("/tasks/{id}", s.handlers.DeleteTask).Methods("DELETE")
	s.router.HandleFunc("/tasks/{id}", s.handlers.CompleteTask).Methods("PATCH")

	if err := http.ListenAndServe(s.addr, s.router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

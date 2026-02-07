package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	handlers *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{
		handlers: handler,
	}
}

func StartServer(r *Router) error {
	router := mux.NewRouter()

	router.HandleFunc("/tasks", r.handlers.CreateTask).Methods("POST")
	router.HandleFunc("/tasks", r.handlers.GetAllTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", r.handlers.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", r.handlers.DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", r.handlers.CompleteTask).Methods("PATCH")

	if err := http.ListenAndServe(":8080", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

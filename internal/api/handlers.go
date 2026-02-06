package api

import (
	"net/http"

	"github.com/punnch/go-todo/internal/todo"
)

type Handler struct {
	service *todo.TodoService
}

func NewHandler(todoService *todo.TodoService) *Handler {
	return &Handler{
		service: todoService,
	}
}

/*
path: /tasks
method: POST
info: json

succeed:
  - status: 201 Created
  - response body: json

fail:
  - status: 400, 409, 500...
  - response body: error + time
*/
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {

}

/*
path: /tasks
method: GET
info: json

succeed:
  - status: 200 OK
  - response body: json

fail:
  - status: 500...
  - response body: error + time
*/
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {

}

/*
path: /tasks/{id}
method: GET
info: json

succeed:
  - status: 200 OK
  - response body: json

fail:
  - status: 400, 404, 500...
  - response body: error + time
*/
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {

}

/*
path: /tasks/{id}
method: DELETE
info: json

succeed:
  - status: 204 No content
  - response body: -

fail:
  - status: 400, 404, 500...
  - response body: error + time
*/
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
}

/*
path: /tasks/{id}
method: PATCH
info: json

succeed:
  - status: 200 OK
  - response body: json

fail:
  - status: 400, 404, 500...
  - response body: error + time
*/
func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {

}

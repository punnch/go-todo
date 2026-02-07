package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/punnch/go-todo/internal/api/dto"
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
	var taskDTO dto.Task
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateToCreate(); err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.service.CreateTask(r.Context(), taskDTO.Title, taskDTO.Description)
	if err != nil {
		dto.ErrorCompareJSON(w, err, todo.ErrNotFound, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	b := dto.ToJSON(task)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response body:", err)
		return
	}
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
	tasks, err := h.service.GetAllTasks(r.Context())
	if err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	b := dto.ToJSON(tasks)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response body:", err)
		return
	}
}

/*
path: /tasks/{id}
method: GET
info: path

succeed:
- status: 200 OK
- response body: json

fail:
- status: 400, 404, 500...
- response body: error + time
*/
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		dto.ErrorCompareJSON(w, err, todo.ErrNotFound, http.StatusNotFound)
		return
	}

	b := dto.ToJSON(task)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response body:", err)
		return
	}
}

/*
path: /tasks/{id}
method: DELETE
info: path

succeed:
- status: 204 No content
- response body: -

fail:
- status: 400, 404, 500...
- response body: error + time
*/
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		dto.ErrorCompareJSON(w, err, todo.ErrNotFound, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
path: /tasks/{id}
method: PATCH
info: path + json

succeed:
- status: 200 OK
- response body: json

fail:
- status: 400, 404, 500...
- response body: error + time
*/
func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeTaskDTO dto.CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeTaskDTO); err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.service.CompleteTask(r.Context(), id)
	if err != nil {
		dto.ErrorCompareJSON(w, err, todo.ErrNotFound, http.StatusNotFound)
		return
	}

	b := dto.ToJSON(task)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response body:", err)
		return
	}
}

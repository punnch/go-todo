package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/punnch/go-todo/internal/core/apperrors"
	"github.com/punnch/go-todo/internal/services/api-service/api/dto"
	"github.com/punnch/go-todo/internal/services/api-service/client"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	todoClient *client.TodoClient
	log        *zap.Logger
}

func NewHandler(todoClient *client.TodoClient, log *zap.Logger) *Handler {
	return &Handler{
		todoClient: todoClient,
		log:        log,
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
		h.log.Warn("invalid request body", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateToCreate(); err != nil {
		h.log.Warn("failed task validation", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.todoClient.CreateTask(r.Context(), taskDTO.Title, taskDTO.Description)
	if err != nil {
		h.log.Error("failed to create task", zap.Error(err))
		dto.ErrorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	b := dto.ToJSON(task)
	if _, err := w.Write(b); err != nil {
		h.log.Error("failed to write http response body", zap.Error(err))
		return
	}

	h.log.Info("task created", zap.Int("id", task.ID))
}

/*
path: /tasks?id={id}&completed={completed}
method: GET
info: query params

succeed:
- status: 200 OK
- response body: json

fail:
- status: 500...
- response body: error + time
*/
func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	idStr := q.Get("id")
	completedStr := q.Get("completed")

	var id *int
	if idStr != "" {
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			h.log.Warn("invalid query param 'id'", zap.Error(err))
			dto.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		id = &idInt
	}

	var completed *bool
	if completedStr != "" {
		completedBool, err := strconv.ParseBool(completedStr)
		if err != nil {
			h.log.Warn("invalid query param 'completed'", zap.Error(err))
			dto.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		completed = &completedBool
	}

	tasks, err := h.todoClient.GetAllTasks(r.Context(), id, completed)
	if err != nil {
		h.log.Error("failed to get all tasks", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	b := dto.ToJSON(tasks)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.log.Error("failed wo write http response body", zap.Error(err))
		return
	}

	h.log.Info("all tasks gotten", zap.Int("amount", len(tasks)))
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
		h.log.Warn("invalid id", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.todoClient.GetTask(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get task", zap.Error(err))
		dto.ErrorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	b := dto.ToJSON(task)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.log.Error("failed to write http response body", zap.Error(err))
		return
	}

	h.log.Info("task gotten", zap.Int("id", task.ID))
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
		h.log.Warn("invalid id", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := h.todoClient.DeleteTask(r.Context(), id); err != nil {
		h.log.Error("failed to delete task", zap.Error(err))
		dto.ErrorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	h.log.Info("task deleted", zap.Int("id", id))
}

/*
path: /tasks/{id}
method: PATCH
info: path

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
		h.log.Warn("invalid request body", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Warn("invalid id", zap.Error(err))
		dto.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	task, err := h.todoClient.CompleteTask(r.Context(), id, completeTaskDTO.Completed)
	if err != nil {
		h.log.Error("failed to complete task", zap.Error(err))
		dto.ErrorCompareJSON(w, err, apperrors.ErrTaskNotFound, http.StatusNotFound)
		return
	}

	b := dto.ToJSON(task)

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		h.log.Error("failed to write http response body", zap.Error(err))
		return
	}

	h.log.Info("task completed", zap.Int("id", id))
}

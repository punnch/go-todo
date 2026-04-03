package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/punnch/go-todo/internal/core/apperrors"
	"github.com/punnch/go-todo/internal/core/domains/todo"
	"github.com/punnch/go-todo/internal/services/api-service/api/dto"
)

type TodoClient struct {
	baseURL string
	client  *http.Client
}

func NewTodoClient(port string) *TodoClient {
	return &TodoClient{
		baseURL: fmt.Sprintf("http://db-service:%s", port),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *TodoClient) CreateTask(ctx context.Context, title, description string) (todo.Task, error) {
	url := c.baseURL + "/tasks"

	// Initialize payload to create task
	payload := map[string]string{
		"title":       title,
		"description": description,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return todo.Task{}, err
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return todo.Task{}, err
	}

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return todo.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var clientErrorDTO dto.ClientErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&clientErrorDTO); err != nil {
			return todo.Task{}, err
		}

		return todo.Task{}, clientErrorDTO.ToError()
	}

	// Decode response from the server to task
	var task todo.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return todo.Task{}, err
	}

	return task, nil
}

func (c *TodoClient) GetAllTasks(ctx context.Context, id *int, completed *bool) ([]todo.Task, error) {
	url := c.baseURL + "/tasks"

	var params []string

	if id != nil {
		params = append(params, fmt.Sprintf("id=%d", *id))
	}

	if completed != nil {
		params = append(params, fmt.Sprintf("completed=%v", *completed))
	}

	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Don't check status code, cause it will be always OK
	var tasks []todo.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *TodoClient) GetTask(ctx context.Context, id int) (todo.Task, error) {
	url := c.baseURL + fmt.Sprintf("/tasks/%d", id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return todo.Task{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return todo.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return todo.Task{}, apperrors.ErrTaskNotFound
		}

		var clientErrorDTO dto.ClientErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&clientErrorDTO); err != nil {
			return todo.Task{}, err
		}

		return todo.Task{}, clientErrorDTO.ToError()
	}

	var task todo.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return todo.Task{}, err
	}

	return task, nil
}

func (c *TodoClient) DeleteTask(ctx context.Context, id int) error {
	url := c.baseURL + fmt.Sprintf("/tasks/%d", id)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusNotFound {
			return apperrors.ErrTaskNotFound
		}

		var clientErrorDTO dto.ClientErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&clientErrorDTO); err != nil {
			return err
		}

		return clientErrorDTO.ToError()
	}

	return nil
}

func (c *TodoClient) CompleteTask(ctx context.Context, id int, completed *bool) (todo.Task, error) {
	url := c.baseURL + fmt.Sprintf("/tasks/%d", id)

	payload := map[string]*bool{"completed": completed}

	data, err := json.Marshal(payload)
	if err != nil {
		return todo.Task{}, err
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return todo.Task{}, err
	}

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return todo.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return todo.Task{}, apperrors.ErrTaskNotFound
		}

		var clientErrorDTO dto.ClientErrorDTO
		if err := json.NewDecoder(resp.Body).Decode(&clientErrorDTO); err != nil {
			return todo.Task{}, err
		}

		return todo.Task{}, clientErrorDTO.ToError()
	}

	// Decode response from the server to task
	var task todo.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return todo.Task{}, err
	}

	return task, nil
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/punnch/go-todo/internal/core/domains/todo"
)

type TodoClient struct {
	baseURL string
	client  *http.Client
}

func NewTodoClient(url string) *TodoClient {
	return &TodoClient{
		baseURL: url,
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
		return todo.Task{}, fmt.Errorf("failed to create task")
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
		return todo.Task{}, fmt.Errorf("failed to get task")
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
		return fmt.Errorf("failed to delete task")
	}

	return nil
}

func (c *TodoClient) CompleteTask(ctx context.Context, id int, completed *bool) (todo.Task, error) {
	url := c.baseURL + fmt.Sprintf("/tasks/%d", id)

	payload := map[string]*bool{"Completed": completed}

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
		return todo.Task{}, fmt.Errorf("failed to complete task")
	}

	// Decode response from the server to task
	var task todo.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return todo.Task{}, err
	}

	return task, nil
}

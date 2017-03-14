package squarescale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	StatusWaiting   = "waiting"
	StatusRunning   = "running"
	StatusCancelled = "cancelled"
	StatusDone      = "done"
)

type Task struct {
	Status string `json:"status"`
}

func (c *Client) GetTask(id int) (task Task, err error) {
	code, body, err := c.get(fmt.Sprintf("/tasks/%d", id))
	if err != nil {
		return
	}

	if code != http.StatusOK {
		err = unexpectedHTTPError(code, body)
		return
	}

	err = json.Unmarshal(body, &task)
	return
}

func (c *Client) WaitTask(id int, interval time.Duration) (status string, err error) {
	for status != StatusDone && status != StatusCancelled {
		task, err := c.GetTask(id)
		if err != nil {
			return status, err
		}

		status = task.Status
		time.Sleep(interval)
	}

	return
}

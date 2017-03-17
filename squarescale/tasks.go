package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	StatusWaiting   = "waiting"
	StatusRunning   = "running"
	StatusCancelled = "cancelled"
	StatusDone      = "done"
)

var ErrInProgress error = errors.New("Task still in progress")

type Task struct {
	Id            int             `json:"id"`
	ProjectId     int             `json:"project_id"`
	WaitingEvents []string        `json:"waiting_events"`
	Params        json.RawMessage `json:"params"`
	Status        string          `json:"status"`
	CompletedBy   string          `json:"completed_by"`
	CompletedAt   string          `json:"completed_at"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	Hold          bool            `json:"hold"`
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

func (c *Client) WaitTask(id int) (*Task, error) {
	if id == 0 {
		return nil, ErrInProgress
	}

	task, err := c.GetTask(id)
	if err != nil {
		return &task, err
	}

	for task.Status == StatusDone || task.Status == StatusCancelled {
		return &task, nil
	}

	var subscription struct {
		Channel string `json:"channel"`
		TaskId  int    `json:"task_id"`
	}

	subscription.Channel = "TaskChannel"
	subscription.TaskId = id
	subsc, err := json.Marshal(subscription)
	if err != nil {
		return nil, err
	}

	ch, err := c.cableClient().SubscribeWith(subsc)
	if err != nil {
		return nil, err
	}
	defer c.cableClient().Unsubscribe("TaskChannel")

	for task.Status != StatusDone && task.Status != StatusCancelled {
		ev := <-ch
		if ev.Err != nil {
			return nil, err
		} else if len(ev.Event.Message) > 0 {
			err := json.Unmarshal(ev.Event.Message, &task)
			if err != nil {
				return nil, fmt.Errorf("Could not unmarshal JSON %s: %v", string(ev.Event.Message))
			}
		}
	}

	return &task, nil
}

package squarescale

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/squarescale/logger"
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
	TaskType      string          `json:"type"`
	ProjectId     int             `json:"project_id"`
	WaitingEvents []string        `json:"waiting_events"`
	Params        json.RawMessage `json:"params"`
	Status        string          `json:"status"`
	CompletedBy   string          `json:"completed_by"`
	CompletedAt   string          `json:"completed_at"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	Hold          bool            `json:"hold"`
	CableToken    string          `json:"table_token"`
}

func (t Task) String() string {
	return fmt.Sprintf("task of type %s for project %d with id: %d and status %s", t.TaskType, t.ProjectId, t.Id, t.Status)
}

func (t *Task) CreatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, t.CreatedAt)
}

func (t *Task) UpdatedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, t.UpdatedAt)
}

func (t *Task) CompletedTime() (time.Time, error) {
	return time.Parse(time.RFC3339, t.CompletedAt)
}

func (t *Task) LatestTime(format string) string {
	kind := "created"
	createdTime, errCreated := t.CreatedTime()
	if errCreated != nil {
		return "-"
	}
	updatedTime, errUpdated := t.UpdatedTime()
	completedTime, errCompleted := t.CompletedTime()
	time := createdTime
	if errCompleted == nil && completedTime.After(time) {
		kind = "completed"
		time = completedTime
	}
	if errUpdated == nil && updatedTime.After(time) {
		kind = "updated"
		time = updatedTime
	}
	return kind + " " + time.Format(format)
}

func (t *Task) StatusWithHold() string {
	res := t.Status
	if t.Hold {
		res += " (hold)"
	}
	return res
}

func (c *Client) GetTasks(projectName string) (tasks []Task, err error) {
	code, body, err := c.get(fmt.Sprintf("/projects/%s/tasks", projectName))
	if err != nil {
		return
	}

	if code != http.StatusOK {
		err = unexpectedHTTPError(code, body)
		return
	}

	err = json.Unmarshal(body, &tasks)
	return
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

	logger.Info.Println("wait for task : ", task)

	for task.Status != StatusDone && task.Status != StatusCancelled && err == nil {
		time.Sleep(5 * time.Second)
		task, err = c.GetTask(id)
		logger.Debug.Println("task status update: ", task)
	}

	return &task, err
}

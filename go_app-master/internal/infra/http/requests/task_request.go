package requests

import (
	"net/http"
	"strconv"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/go-chi/chi/v5"
)

type TaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Deadline    *uint64 `json:"deadline"`
}

type UpdateTaskRequest struct {
	Title    *string            `json:"title"`
	Status   *domain.TaskStatus `json:"status"`
	Deadline *uint64            `json:"deadline"`
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {

	var deadline *time.Time
	if r.Deadline != nil {
		ddl := time.Unix(int64(*r.Deadline), 0)
		deadline = &ddl
	}

	/*descr := ""
	if r.Description != nil {
		descr = *r.Description
	}

	var deadline uint64
	if r.Deadline != nil {
		deadline = *r.Deadline
	}*/

	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Deadline:    deadline,
	}, nil
}

func ParseTaskId(r *http.Request) (uint64, error) {
	taskIdStr := chi.URLParam(r, "taskId")
	taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return taskId, nil
}

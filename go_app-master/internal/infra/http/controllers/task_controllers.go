package controllers

import (
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.NewTakStatus
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) FindByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		tasks, err := c.taskService.FindByUserId(user.Id)
		if err != nil {
			log.Printf("TaskController -> FindBuUserId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tsDto resources.TasksDto
		tsDto = tsDto.DomainToDtoCollection(tasks)
		Success(w, tsDto)
	}
}

func (c TaskController) FindByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := requests.ParseTaskId(r)
		if err != nil {
			log.Printf("TaskController -> FindByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		task, err := c.taskService.FindByTaskId(taskId)
		if err != nil {
			log.Printf("TaskController -> FindByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}

func (c TaskController) DeleteByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := requests.ParseTaskId(r)
		if err != nil {
			log.Printf("TaskController -> DeleteByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.taskService.DeleteByTaskId(taskId)
		if err != nil {
			log.Printf("TaskController -> DeleteByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (c TaskController) UpdateByTaskId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId, err := requests.ParseTaskId(r)
		if err != nil {
			log.Printf("TaskController -> UpdateByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		taskReq, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> UpdateByTaskId: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		taskReq.UserId = user.Id
		taskReq.Id = taskId

		updatedTask, err := c.taskService.UpdateByTaskId(taskReq)
		if err != nil {
			log.Printf("TaskController -> UpdateByTaskId: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(updatedTask)
		Success(w, tDto)
	}
}

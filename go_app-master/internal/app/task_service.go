package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	FindByUserId(uId uint64) ([]domain.Task, error)
	DeleteByTaskId(taskId uint64) error
	FindByTaskId(tId uint64) (domain.Task, error)
	UpdateByTaskId(t domain.Task) (domain.Task, error)
}

type taskService struct {
	taskRepo database.TaskRepositiry
}

func NewTaskService(tr database.TaskRepositiry) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("TaskService -> Save %s", err)
		return domain.Task{}, err
	}
	return task, nil
}

func (r taskService) FindByUserId(uId uint64) ([]domain.Task, error) {
	tasks, err := r.taskRepo.FindByUserId(uId)

	if err != nil {
		log.Printf("TaskService -> FindByUserId %s", err)
		return nil, err
	}

	return tasks, nil
	//return r.mapModelToDomain(tsk), nil
}

func (s taskService) DeleteByTaskId(taskId uint64) error {
	err := s.taskRepo.DeleteByTaskId(taskId)
	if err != nil {
		log.Printf("TaskService -> DeleteByTaskId %s", err)
		return err
	}
	return nil
}

func (s taskService) FindByTaskId(taskId uint64) (domain.Task, error) {
	task, err := s.taskRepo.FindByTaskId(taskId)
	if err != nil {
		log.Printf("TaskService -> FindByTaskId %s", err)
		return domain.Task{}, err
	}
	return task, nil
}

func (s taskService) UpdateByTaskId(t domain.Task) (domain.Task, error) {
	updatedTask, err := s.taskRepo.UpdateByTaskId(t)
	if err != nil {
		log.Printf("TaskService -> UpdateByTaskId %s", err)
		return domain.Task{}, err
	}
	return updatedTask, nil
}

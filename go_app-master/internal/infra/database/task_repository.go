package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description *string           `db:"description"`
	Status      domain.TaskStatus `db:"status"`
	Deadline    *time.Time        `db:"deadline"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepositiry interface {
	Save(t domain.Task) (domain.Task, error)
	FindByUserId(uId uint64) ([]domain.Task, error)
	DeleteByTaskId(taskId uint64) error
	FindByTaskId(taskId uint64) (domain.Task, error)
	UpdateByTaskId(t domain.Task) (domain.Task, error)
}

type taskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(session db.Session) TaskRepositiry {
	return taskRepository{
		coll: session.Collection(TasksTableName),
		sess: session,
	}
}

func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
	//return r.mapModelToDomain(tsk), nil
}

func (r taskRepository) FindByUserId(uId uint64) ([]domain.Task, error) {
	var tasks []task
	err := r.coll.
		Find(db.Cond{"user_id": uId, "deleted_date": nil}).
		OrderBy("-deadline").
		All(&tasks)
	if err != nil {
		return nil, err
	}

	result := r.mapModelToDomainCollection(tasks)
	return result, nil
}

/*-------------------*/
func (r taskRepository) DeleteByTaskId(taskId uint64) error {
	err := r.coll.Find(db.Cond{"id": taskId, "deleted_date": nil}).Update(map[string]interface{}{
		"deleted_date": time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r taskRepository) FindByTaskId(taskId uint64) (domain.Task, error) {
	var tsk task
	err := r.coll.Find(db.Cond{"id": taskId, "deleted_date": nil}).One(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(tsk), nil
}

// перетаорення в домайн та з домейну
func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	tasks := make([]domain.Task, len(ts))
	for i, t := range ts {
		tasks[i] = r.mapModelToDomain(t)
	}
	return tasks
}

func (r taskRepository) UpdateByTaskId(t domain.Task) (domain.Task, error) {
	var existingTask task
	err := r.coll.Find(db.Cond{"id": t.Id, "deleted_date": nil}).One(&existingTask)
	if err != nil {
		return domain.Task{}, err
	}

	if t.Title != "" {
		existingTask.Title = t.Title
	}
	if t.Status != "" {
		existingTask.Status = t.Status
	}
	if t.Deadline != nil {
		existingTask.Deadline = t.Deadline
	}

	existingTask.UpdatedDate = time.Now()
	err = r.coll.UpdateReturning(&existingTask)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(existingTask), nil
}

package uc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/utsavgupta/go-todo/entities"
	"github.com/utsavgupta/go-todo/logger"
	"github.com/utsavgupta/go-todo/repos"
)

type GetTaskUC func(ctx context.Context, taskId int) (*entities.Task, error)
type ListTasksUC func(ctx context.Context) (entities.Tasks, error)
type CreateTaskUC func(ctx context.Context, task entities.Task) (*entities.Task, error)
type UpdateTaskUC func(ctx context.Context, task entities.Task) (*entities.Task, error)
type DeleteTaskUC func(ctx context.Context, taskId int) error

func NewGetTaskUC(repo repos.TaskRepo) GetTaskUC {

	return func(ctx context.Context, taskId int) (*entities.Task, error) {

		entity, err := repo.Get(ctx, taskId)

		if err != nil {

			logger.Instance().Info(ctx, fmt.Sprintf("could not get task %v", taskId))
		}

		return entity, err
	}
}

func NewListTasksUC(repo repos.TaskRepo) ListTasksUC {

	return func(ctx context.Context) (entities.Tasks, error) {

		entities, err := repo.List(ctx)

		if err != nil {

			logger.Instance().Info(ctx, "could not list tasks")
		}

		return entities, err
	}
}

func NewCreateTaskUC(repo repos.TaskRepo) CreateTaskUC {

	return func(ctx context.Context, task entities.Task) (*entities.Task, error) {

		if validationError := validateTask(task); validationError != nil {

			logger.Instance().Info(ctx, fmt.Sprintf("validation failed for task %v", task))
			return nil, validationError
		}

		setCreationDate(&task)
		markIncomplete(&task)

		idx, err := repo.Create(ctx, task)

		if err != nil {

			logger.Instance().Info(ctx, fmt.Sprintf("could not create task %v", task))
			return nil, err
		}

		task.Id = idx

		return &task, nil
	}
}

func setCreationDate(task *entities.Task) {

	now := time.Now()
	task.CreatedAt = &now
}

func markIncomplete(task *entities.Task) {

	task.Completed = false
}

func NewUpdateTaskUC(repo repos.TaskRepo) UpdateTaskUC {

	return func(ctx context.Context, task entities.Task) (*entities.Task, error) {

		if validationError := validateTask(task); validationError != nil {

			logger.Instance().Info(ctx, fmt.Sprintf("validation failed for task %v", task))
			return nil, validationError
		}

		setUpdationDate(&task)
		setCompletionDate(&task)

		err := repo.Update(ctx, task)

		if err != nil {

			logger.Instance().Info(ctx, fmt.Sprintf("could not update task %v", task))
			return nil, err
		}

		return &task, nil
	}
}

func setUpdationDate(task *entities.Task) {

	now := time.Now()
	task.UpdatedAt = &now
}

func setCompletionDate(task *entities.Task) {

	if task.Completed {
		now := time.Now()
		task.CompletedAt = &now
	}
}

func validateTask(task entities.Task) error {

	err := validateDescription(task.Description)

	return err
}

func validateDescription(description string) error {

	if len(description) < 5 || len(description) > 50 {

		return errors.New("the task description must be between 5 and 50 characters")
	}

	return nil
}

func NewDeleteTaskUC(repo repos.TaskRepo) DeleteTaskUC {

	return func(ctx context.Context, taskId int) error {

		err := repo.Delete(ctx, taskId)

		if err != nil {

			logger.Instance().Info(ctx, fmt.Sprintf("could not delete task %v", taskId))
		}

		return err
	}
}

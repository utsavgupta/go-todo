package repos

import (
	"context"

	"github.com/utsavgupta/go-todo/entities"
)

type TaskRepo interface {
	Get(context.Context, int) (*entities.Task, error)
	List(context.Context) (entities.Tasks, error)
	Create(context.Context, entities.Task) (int, error)
	Update(context.Context, entities.Task) error
	Delete(context.Context, int) error
}

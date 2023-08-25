package datasources

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/utsavgupta/go-todo/entities"
)

type pgTaskRepo struct {
	conn *pgx.Conn
}

func NewPGTaskRepo(connStr string) (*pgTaskRepo, error) {

	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return nil, err
	}

	return &pgTaskRepo{conn}, nil
}

func (repo pgTaskRepo) Get(ctx context.Context, id int) (*entities.Task, error) {

	var task entities.Task

	row, err := repo.conn.Query(ctx, "SELECT id, description, completed, created_at, updated_at, completed_at FROM tasks WHERE id = $1", id)

	if err != nil {

		return nil, fmt.Errorf("could not fetch task with id %d: %w", id, err)
	}

	defer row.Close()

	if !row.Next() {

		return nil, nil
	}

	if err = row.Scan(&task.Id, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt, &task.CompletedAt); err != nil {

		return nil, fmt.Errorf("could not fetch task with id %d: %w", id, err)
	}

	return &task, nil
}

func (repo pgTaskRepo) List(ctx context.Context) (entities.Tasks, error) {

	var tasks entities.Tasks

	row, err := repo.conn.Query(ctx, "SELECT id, description, completed, created_at, updated_at, completed_at FROM tasks")

	if err != nil {

		return nil, fmt.Errorf("could not list tasks: %w", err)
	}

	defer row.Close()

	tasks = make(entities.Tasks, 0)

	for row.Next() {

		task := entities.Task{}

		if err = row.Scan(&task.Id, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt, &task.CompletedAt); err != nil {

			return nil, fmt.Errorf("could not read task: %w", err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (repo pgTaskRepo) Create(ctx context.Context, task entities.Task) (int, error) {

	row := repo.conn.QueryRow(ctx, "INSERT INTO tasks (description, created_at) VALUES ($1, $2) RETURNING id", task.Description, task.CreatedAt)
	idx := -1

	if err := row.Scan(&idx); err != nil {

		err = fmt.Errorf("could not create task %v: %w", task, err)
	}

	return idx, nil
}

func (repo pgTaskRepo) Update(ctx context.Context, task entities.Task) error {

	var err error

	if task.Completed {

		_, err = repo.conn.Exec(ctx, "UPDATE tasks SET description = $1, completed = $2, updated_at = $3, completed_at = $4 WHERE id = $5", task.Description, task.Completed, task.UpdatedAt, task.CompletedAt, task.Id)
	} else {

		_, err = repo.conn.Exec(ctx, "UPDATE tasks SET description = $1, completed = $2, updated_at = $3 WHERE id = $4", task.Description, task.Completed, task.UpdatedAt, task.Id)
	}

	if err != nil {

		err = fmt.Errorf("could not update task %v: %w", task, err)
	}

	return err
}

func (repo pgTaskRepo) Delete(ctx context.Context, id int) error {

	var err error

	if _, err = repo.conn.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id); err != nil {

		err = fmt.Errorf("could not delete task with id %d: %w", id, err)
	}

	return err
}

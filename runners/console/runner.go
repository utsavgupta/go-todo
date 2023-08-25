package console

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/utsavgupta/go-todo/adapters/datasources"
	"github.com/utsavgupta/go-todo/entities"
	"github.com/utsavgupta/go-todo/logger"
	"github.com/utsavgupta/go-todo/runners"
	"github.com/utsavgupta/go-todo/uc"
)

type consoleRunner struct {
	inStr      io.Reader
	outStr     io.Writer
	errStr     io.Writer
	listTasks  uc.ListTasksUC
	getTask    uc.GetTaskUC
	createTask uc.CreateTaskUC
	updateTask uc.UpdateTaskUC
	deleteTask uc.DeleteTaskUC
}

func NewConsoleRunner(inStr io.Reader, outStr, errStr io.Writer) (runners.Runner, error) {

	repo, err := datasources.NewPGTaskRepo("postgres://postgres:postgres@localhost:5432/go-todo")

	if err != nil {
		return nil, err
	}

	getTaskUC := uc.NewGetTaskUC(repo)
	listTasksUC := uc.NewListTasksUC(repo)
	createTaskUC := uc.NewCreateTaskUC(repo)
	updateTaskUC := uc.NewUpdateTaskUC(repo)
	deleteTaskUC := uc.NewDeleteTaskUC(repo)

	return &consoleRunner{inStr, outStr, errStr, listTasksUC, getTaskUC, createTaskUC, updateTaskUC, deleteTaskUC}, nil
}

func (c *consoleRunner) Run() error {

	initLogger()

	for {
		fmt.Fprintf(c.outStr, "Choose option:\n")
		fmt.Fprintf(c.outStr, "1. List\n")
		fmt.Fprintf(c.outStr, "2. Get\n")
		fmt.Fprintf(c.outStr, "3. Create new\n")
		fmt.Fprintf(c.outStr, "4. Update description\n")
		fmt.Fprintf(c.outStr, "5. Mark as complete\n")
		fmt.Fprintf(c.outStr, "6. Delete\n")
		fmt.Fprintf(c.outStr, "7. Exit\n")

		selection := c.readInt()

		switch selection {
		case 1:
			c.list()
		case 2:
			c.get()
		case 3:
			c.create()
		case 4:
			c.updateDescription()
		case 5:
			c.markAsCompleted()
		case 6:
			c.delete()
		case 7:
			fmt.Fprintf(c.outStr, "Exiting...\n")
			return nil
		default:
			fmt.Fprintf(c.outStr, "Invalid option\n")
		}
	}
}

func initLogger() {

	loggerInstance := logger.NewZeroLogger()
	logger.InitLogger(loggerInstance)
}

func (c *consoleRunner) list() {

	tasks, err := c.listTasks(context.Background())

	if err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	fmt.Fprintf(c.outStr, "%-3s %-50s %-5s\n", "Id", "Description", "Status")
	fmt.Fprintf(c.outStr, "%-3s %-50s %-5s\n", strings.Repeat("+", 3), strings.Repeat("+", 50), strings.Repeat("+", 6))

	for _, task := range tasks {

		fmt.Fprintf(c.outStr, "%-3d %-50s %-5s\n", task.Id, task.Description, c.statusString(task.Completed))
	}
}

func (c *consoleRunner) get() {

	fmt.Fprintf(c.outStr, "Enter Id: ")
	id := c.readInt()

	task, err := c.getTask(context.Background(), id)

	if err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	if task == nil {

		fmt.Fprintf(c.outStr, "Task with %d Id was not found\n", id)
		return
	}

	fmt.Fprintf(c.outStr, "Id: %d\n", task.Id)
	fmt.Fprintf(c.outStr, "Description: %s\n", task.Description)
	fmt.Fprintf(c.outStr, "Completed: %s\n", c.statusString(task.Completed))
	fmt.Fprintf(c.outStr, "Created At: %v\n", task.CreatedAt)
	fmt.Fprintf(c.outStr, "Updated At: %v\n", task.UpdatedAt)
	fmt.Fprintf(c.outStr, "Completed At: %v\n", task.CompletedAt)
}

func (c *consoleRunner) create() {

	var err error
	task := &entities.Task{}

	fmt.Fprintf(c.outStr, "Description: ")
	task.Description = c.readLine()

	if task, err = c.createTask(context.Background(), *task); err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	fmt.Fprintf(c.outStr, "Task created with id %d\n", task.Id)
}

func (c *consoleRunner) delete() {

	fmt.Fprintf(c.outStr, "Enter Id: ")
	id := c.readInt()

	if err := c.deleteTask(context.Background(), id); err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	fmt.Fprintf(c.outStr, "The task was deleted successfully\n")
}

func (c *consoleRunner) markAsCompleted() {

	fmt.Fprintf(c.outStr, "Enter Id: ")
	id := c.readInt()

	task, err := c.getTask(context.Background(), id)

	if err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	if task == nil {

		fmt.Fprintf(c.outStr, "Task with %d Id was not found\n", id)
		return
	}

	task.Completed = true

	if _, err := c.updateTask(context.Background(), *task); err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	fmt.Fprintf(c.outStr, "The task was successfully marked as completed\n")
}

func (c *consoleRunner) updateDescription() {

	fmt.Fprintf(c.outStr, "Enter Id: ")
	id := c.readInt()

	task, err := c.getTask(context.Background(), id)

	if err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	if task == nil {

		fmt.Fprintf(c.outStr, "Task with %d Id was not found\n", id)
		return
	}

	fmt.Fprintf(c.outStr, "New description: ")
	task.Description = c.readLine()

	if _, err := c.updateTask(context.Background(), *task); err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return
	}

	fmt.Fprintf(c.outStr, "The task description was successfully updated\n")
}

func (c *consoleRunner) statusString(status bool) string {

	if status {
		return "Yes"
	}

	return "No"
}

func (c *consoleRunner) readLine() string {

	br := bufio.NewReader(c.inStr)
	b, _, err := br.ReadLine()

	if err != nil {

		logger.Instance().Error(context.Background(), err.Error())
		return ""
	}

	return string(b)
}

func (c *consoleRunner) readInt() int {

	str := c.readLine()
	i, err := strconv.Atoi(str)

	if err != nil {

		i = -1
		fmt.Fprintf(c.outStr, "Invalid input\n")
	}

	return i
}

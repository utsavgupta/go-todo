package web

import (
	"net/http"

	"github.com/utsavgupta/go-todo/adapters/datasources"
	routers "github.com/utsavgupta/go-todo/adapters/transport/web"
	"github.com/utsavgupta/go-todo/runners"
	"github.com/utsavgupta/go-todo/uc"
)

type webRunner struct {
	handler http.Handler
}

func NewWebRunner() (runners.Runner, error) {

	repo, err := datasources.NewPGTaskRepo("postgres://postgres:postgres@localhost:5432/go-todo")

	if err != nil {
		return nil, err
	}

	getUC := uc.NewGetTaskUC(repo)
	listUC := uc.NewListTasksUC(repo)
	createUC := uc.NewCreateTaskUC(repo)
	updateUC := uc.NewUpdateTaskUC(repo)
	deleteUC := uc.NewDeleteTaskUC(repo)

	return &webRunner{routers.NewRouter(&routers.RouterDependencies{
		GetUC:    getUC,
		ListUC:   listUC,
		CreateUC: createUC,
		UpdateUC: updateUC,
		DeleteUC: deleteUC,
	})}, nil
}

func (runner *webRunner) Run() error {

	go func() { http.ListenAndServe(":8080", runner.handler) }()

	return nil
}

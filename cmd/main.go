package main

import (
	"os"
	"os/signal"

	"github.com/utsavgupta/go-todo/runners/web"
)

func main() {

	runner, err := web.NewWebRunner()

	if err != nil {
		panic(err)
	}

	if err := runner.Run(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		return
	}
}

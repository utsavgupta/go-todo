package main

import (
	"os"

	"github.com/utsavgupta/go-todo/runners/console"
)

func main() {

	runner, err := console.NewConsoleRunner(os.Stdin, os.Stdout, os.Stderr)

	if err != nil {
		panic(err)
	}

	if err := runner.Run(); err != nil {
		panic(err)
	}
}

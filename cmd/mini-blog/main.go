package main

import (
	miniblog "github.com/liaomars/mini-blog/internal/mini-blog"
	_ "go.uber.org/automaxprocs"
	"os"
)

func main() {
	command := miniblog.NewMiniBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

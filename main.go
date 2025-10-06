package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"

	"github.com/Parapheen/fullkek-starter/cmd"
)

func main() {
	if err := fang.Execute(context.Background(), cmd.RootCommand()); err != nil {
		os.Exit(1)
	}
}

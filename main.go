package main

import (
	"fmt"
	"os"

	"github.com/Old-Goggles/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		return
	}

	programState := &state{
		Cfg: &cfg,
	}

	cmds := &commands{
		Handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/migueldor/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	var newState state
	newState.config = &cfg
	var myCommands commands
	myCommands.commandHandlers = make(map[string]func(*state, command) error)
	myCommands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments passed")
		os.Exit(1)
	}
	cmdName := args[1]
	cmdArgs := args[2:]
	newCommand := command{
		name: cmdName,
		args: cmdArgs,
	}
	err = myCommands.run(&newState, newCommand)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/migueldor/gator/internal/config"
	"github.com/migueldor/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	dbURL := cfg.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(db)
	newState := state{
		db:     dbQueries,
		config: &cfg,
	}
	var myCommands commands
	myCommands.commandHandlers = make(map[string]func(*state, command) error)
	myCommands.register("login", handlerLogin)
	myCommands.register("register", handlerRegister)
	myCommands.register("reset", handlerReset)
	myCommands.register("users", handlerUsers)
	myCommands.register("agg", handlerAgg)
	myCommands.register("addfeed", handlerAddFeed)
	myCommands.register("feeds", handlerGetFeeds)
	myCommands.register("follow", handlerFollow)
	myCommands.register("following", handlerFollowing)
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

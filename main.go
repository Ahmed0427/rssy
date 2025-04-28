package main

import (
	"fmt"
	"os"
	"log"
	"database/sql"

	"github.com/Ahmed0427/rssy/internal/config"
	"github.com/Ahmed0427/rssy/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

type Command struct {
	name string
	args []string
}

type Commands struct {
	handlersMap map[string]func(*State, Command) error
}

func (cmds *Commands) register(name string, handler func(*State, Command) error) {
	cmds.handlersMap[name] = handler
}

func (cmds *Commands) registerAll() {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("help", handlerHelp)
	cmds.register("aggregate", handlerAggregate)
	cmds.register("addfeed", handlerAddfeed)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)
}

func (cmds *Commands) run(s *State, cmd Command) error {
	handler, found := cmds.handlersMap[cmd.name]	
	if !found {
		return fmt.Errorf("Error: command not found")	
	}
	err := handler(s, cmd);
	if err != nil {
		return err
	}

	return nil
}

func main() {
	godotenv.Load()
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not in the environment")	
	}
	
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.ConnStr)
	dbQueries := database.New(db)

	state := State {
		db: dbQueries,
		cfg: &cfg,	
	}

	cmds := Commands {
		handlersMap: make(map[string]func(*State, Command) error),
	}
	cmds.registerAll()
	
	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments were provided.")
		os.Exit(1)
	}

	cmd := Command {
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&state, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


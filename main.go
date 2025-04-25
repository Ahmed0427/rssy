package main

import (
	"fmt"
	"os"
	"log"
	"time"
	"context"
	"database/sql"

	"github.com/Ahmed0427/rssy/internal/config"
	"github.com/Ahmed0427/rssy/internal/database"

	"github.com/google/uuid"
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

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: login command expects <username> argument")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error: %s is not in the database", cmd.args[0])
	}
		
	s.cfg.Username = cmd.args[0]
	s.cfg.Write()

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: register command expects <username> argument")
	}

	userParams := database.CreateUserParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.args[0],
	}
	
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' has been successfully registered.\n", user.Name)

	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("All users have been successfully deleted from the database.")

	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("- %s", user.Name)
		if (s.cfg.Username == user.Name) {
			fmt.Printf(" (current)")
		}
		fmt.Println()
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

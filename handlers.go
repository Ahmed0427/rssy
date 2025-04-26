package main

import (
	"fmt"
	"time"
	"context"

	"github.com/Ahmed0427/rssy/internal/database"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

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

func handlerAggregate(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: aggregate command expects <URL> argument")
	}

	feed, err := fetchFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("Channel Title:", feed.Channel.Title)
	fmt.Println("Channel Description:", feed.Channel.Description)
	fmt.Println("Items:")

	for i, item := range feed.Channel.Item {
		fmt.Printf("\nItem #%d\n", i+1)
		fmt.Println("Title:", item.Title)
		fmt.Println("Description:", item.Description)
	}

	return nil
}

func handlerHelp(s *State, cmd Command) error {
	fmt.Println()
	fmt.Println("help                  -- Display this help message")
	fmt.Println("login <username>      -- Log in as an existing user")
	fmt.Println("register <username>   -- Register a new user")
	fmt.Println("users                 -- List all registered users")
	fmt.Println("aggregate <URL>       -- fetches updates from the site's RSS feed")

	return nil
}


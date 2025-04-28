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
		return fmt.Errorf("Error: login command expects <username>")
	}

	_, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error: %s is not in the database", cmd.args[0])
	}
		
	s.cfg.Username = cmd.args[0]
	s.cfg.Write()

	fmt.Printf("User '%s' has been successfully loged in.\n", cmd.args[0])

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: register command expects <username>")
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

	err = handlerLogin(s, Command {
		name: "login",
		args: []string{cmd.args[0]},
	})

	if err != nil {
		return err
	}

	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	err = s.db.DeleteAllFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("The database has been reseted.")

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

	if len(users) == 0 {
		fmt.Println("No users in the database")
	}

	return nil
}

func handlerAggregate(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: aggregate command expects <URL>")
	}

	feed, err := fetchFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("===================================")
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Println("===================================\n")

	for i, item := range feed.Channel.Item {
		fmt.Printf("Item #%d\n", i+1)
		fmt.Printf("Title       : %s\n", item.Title)
		fmt.Printf("Description : %s\n\n", item.Description)
	}

	return nil
}

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: follow command expects <URL>")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first", s.cfg.Username)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error: '%s' is not in the database", cmd.args[0])
	}

	userFeedParams := database.CreateUserFeedParams {
		ID       : uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID   : user.ID,
		FeedID   : feed.ID,
	}

	userFeed, err := s.db.CreateUserFeed(context.Background(), userFeedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been followed by '%s'.\n",
		userFeed.FeedName, userFeed.UserName)

	return nil
}

func handlerAddfeed(s *State, cmd Command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Error: addfeed command expects <username> <URL>")
	}

	_, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first", s.cfg.Username)
	}

	feedParams := database.CreateFeedParams {
		ID       : uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name     : cmd.args[0],
		Url      : cmd.args[1],
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' has been successfully added by '%s'.\n",
		feed.Name, s.cfg.Username)

	err = handlerFollow(s, Command {
		name: "follow",
		args: []string{cmd.args[1]},
	})

	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}

	for i, feed := range feeds {
		fmt.Printf("Feed #%d\n", i+1)
		fmt.Printf("  Name: %s\n  URL: %s\n",
			feed.Name, feed.Url)

		fmt.Println()
	}

	if len(feeds) == 0 {
		fmt.Println("No Feeds in the database")
	}

	return nil
}

func handlerFollowing (s *State, cmd Command) error {
	userFeeds, err := s.db.GetUserFeedsForUser(
		context.Background(),
		s.cfg.Username,
	)

	if err != nil {
		return fmt.Errorf("Error: '%s' is not in the database", s.cfg.Username)
	}

	for _, userFeed := range userFeeds {
		fmt.Printf(" - %s\n",
			userFeed.FeedName)
	}

	if len(userFeeds) == 0 {
		fmt.Println("No Feeds in the database")
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
	fmt.Println("addfeed <name> <URL>  -- Add a new feed")
	fmt.Println("feeds                 -- List all added feeds")
	fmt.Println("follow <URL>          -- Follow the site's RSS feed")
	fmt.Println("following             -- List all the feeds the current user follows")
	return nil
}


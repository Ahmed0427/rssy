package main

import (
	"fmt"
	"time"
	"context"
	"strconv"

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
		return fmt.Errorf("Error: aggregate command expects <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	for {
		rssFeed, feed, err := scrapeFeeds(context.Background(), s)
		if err != nil {
			return err
		}

		for _, post := range rssFeed.Channel.Item {
			if post.Title == "" || post.Description == "" {
				continue
			}
			params := database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title: post.Title,
				Url: post.Link,
				Description: post.Description, 
				PublishedAt: parseRSSTimeFromat(post.PubDate),
				FeedID: feed.ID,
			}
			post, err := s.db.CreatePost(context.Background(), params)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(post.Title)
		}

		time.Sleep(timeBetweenRequests)
	}

	return nil
}

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Error: follow command expects <URL>")
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first")
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

func handlerUnfollow (s *State, cmd Command) error {
	_, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first")
	}

	params := database.DeleteUserFeedByUserAndURLParams{
		Name: s.cfg.Username,
		Url: cmd.args[0],
	}

	err = s.db.DeleteUserFeedByUserAndURL(context.Background(), params)
	if err != nil {
		return err
	}

	return nil
}

func handlerFollowing (s *State, cmd Command) error {
	_, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first")
	}

	userFeeds, err := s.db.GetUserFeedsForUser(
		context.Background(),
		s.cfg.Username,
	)

	if err != nil {
		return err
	}

	for _, userFeed := range userFeeds {
		fmt.Printf(" - %s\n",
			userFeed.FeedName)
	}

	if len(userFeeds) == 0 {
		fmt.Println("You are not following anyone yet.")
	}

	return nil
}

func handlerAddfeed(s *State, cmd Command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Error: addfeed command expects <username> <URL>")
	}

	_, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first")
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

func handlerBrowse(s *State, cmd Command) error {
	limit := 2
	var err error
	if len(cmd.args) >= 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}

	user, err := s.db.GetUserByName(context.Background(), s.cfg.Username)
	if err != nil {
		return fmt.Errorf("Error: You have to login first")
	}

	posts, err := s.db.GetPostsForUser(context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit: int32(limit),
		},
	)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println("Title:", post.Title)
		fmt.Println(post.Description)
		fmt.Println()
	}	

	return nil
}

func handlerHelp(s *State, cmd Command) error {
	fmt.Println()
	fmt.Println("help                           -- Display this help message")
	fmt.Println("login <username>               -- Log in as an existing user")
	fmt.Println("register <username>            -- Register a new user")
	fmt.Println("users                          -- List all registered users")
	fmt.Println("aggregate <time_between_reqs>  -- Fetch updates from feeds")
	fmt.Println("addfeed <name> <URL>           -- Add a new RSS feed by name and URL")
	fmt.Println("feeds                          -- List all added feeds")
	fmt.Println("follow <URL>                   -- Follow an RSS feed")
	fmt.Println("unfollow <URL>                 -- Unfollow an RSS feed")
	fmt.Println("following                      -- List all feeds the current user is following")
	fmt.Println("browse [limit]                 -- List recent posts (default limit is 2)")
	return nil
}

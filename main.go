package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hamstimusprime/blog_aggregator/internal/config"
	"github.com/hamstimusprime/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to database, %v", err)
	}
	dbQueries := database.New(db)
	defer db.Close()

	s := &state{
		config: &cfg,
		db:     dbQueries,
	}

	cmds := commands{
		handlersMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerDisplayFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

	commandLineInput := os.Args
	commandName := commandLineInput[1]
	commandArgs := commandLineInput[2:]
	cmd := command{Name: commandName, Args: commandArgs}

	if err = cmds.run(s, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if err := containArgs(cmd); err != nil {
		return err
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Println("couldn't find user")
		return err
	}

	if err := s.config.SetUser(cmd.Args[0]); err != nil {
		fmt.Printf("unable to set user name, %v", err)
		return err
	}
	fmt.Printf("username: %v is logged in\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {

	if err := containArgs(cmd); err != nil {
		return err
	}

	now := time.Now()
	id := uuid.New()
	newUserParams := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
	}

	newUser, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Error creating user: %v\n", err)
		return fmt.Errorf("failed to create new user. User might already exist %v", err)
	}

	if err := s.config.SetUser(newUser.Name); err != nil {
		return fmt.Errorf("could not set current user %v", err)
	}

	fmt.Println("registered new user successfully")
	return nil
}

func handlerReset(s *state, cmd command) error {

	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("unable to reset users table %v", err)
	}
	fmt.Println("database successfully reset")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users %v", err)
	}
	if len(users) == 0 {
		fmt.Println("no user has been added")
		return nil
	}
	for _, username := range users {
		if s.config.CurrentUserName == username {
			fmt.Printf("%v (current)\n", username)
		} else {
			fmt.Println(username)
		}

	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {

	if err := containArgs(cmd); err != nil {
		return err
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if err := containArgs(cmd); err != nil {
		return err
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	loggedUser := user
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    loggedUser.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("unable to add feed %v", err)
	}

	cmd.Args[0] = cmd.Args[1]
	if err := handlerFollow(s, cmd, user); err != nil {
		fmt.Println("unable to create entry to feedfollow ")
		return err
	}

	fmt.Println("Feed created successfully:")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("Url: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)

	return nil
}

func handlerDisplayFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get feeds %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("feed Name: %v\n", feed.Feedname)
		fmt.Printf("url: %v\n", feed.Url)
		fmt.Printf("username: %v\n", feed.Username)

	}
	return nil
}
func handlerFollow(s *state, cmd command, user database.User) error {

	if err := containArgs(cmd); err != nil {
		return err
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Printf("unable to fetch feed from database %v\n", err)
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Printf("unable to create entry into feedFollow table %v\n", err)
		return err
	}

	fmt.Printf("current user: %v just followed %v\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	followsList, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("could not fetch list of followed feeds")
		return err
	}
	fmt.Println("feeds you follow: ")
	for _, followsRow := range followsList {
		fmt.Printf("	- %v\n", followsRow.FeedName)
	}
	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	url := cmd.Args[0]
	deleteParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	}
	if err := s.db.DeleteFeedFollow(context.Background(), deleteParams); err != nil {
		fmt.Println("unfollow feed operation failed")
		return err
	}
	return nil
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlersMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlersMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlersMap[cmd.Name]
	if !ok {
		return fmt.Errorf("invalid command")
	}
	return handler(s, cmd)
}

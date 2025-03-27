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
	//read config file from disk and store in cfg
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config file: %v", err)
		return
	}

	//connect to database
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to database, %v", err)
	}
	dbQueries := database.New(db)
	defer db.Close()

	//store config in a state struct
	s := &state{
		config: &cfg,
		db:     dbQueries,
	}

	//create a new instance of the commands map and store handler functions inside of it
	cmds := commands{
		handlersMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAggregate)

	//get the input from the command line when program runs
	/*os.Args would have as its first value the address of the program,
	its second item would be the name of the command(i.e register, login)
	anything after that would be the extra arguments you provide*/
	commandLineInput := os.Args
	commandName := commandLineInput[1]
	commandArgs := commandLineInput[2:]

	//check if arguments are passed and if command used is valid
	if len(commandLineInput) < 3 {
		cmdsWithoutArgs := []string{"reset", "users", "agg"}
		isValidCmd := false
		for _, cmd := range cmdsWithoutArgs {
			if commandName == cmd {
				isValidCmd = true
			}
		}
		if !isValidCmd {
			fmt.Println("error, no arguments provided")
			os.Exit(1)
		}

	}
	if len(commandLineInput) > 3 {
		fmt.Println("warning! multiple arguments provided, only first argument would be used")
	}
	cmd := command{Name: commandName, Args: commandArgs}

	//call command with arguments. run checks if command passed is a valid one.
	if err = cmds.run(s, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// --------------------------------- handlers -------------------------------------------------//
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expected arguments. No arguments provided")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't find user: %v", err)
	}

	if err := s.config.SetUser(cmd.Args[0]); err != nil {
		fmt.Printf("unable to set user name, %v", err)
		return err
	}
	fmt.Println("username has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	/*The register function registers a new user into the database. we first check if any arguments were passed.
	then we check if the user provided is one that already exists
	*/
	if len(cmd.Args) == 0 {
		return fmt.Errorf("expected arguments. No arguments provided")
	}
	now := time.Now()
	id := uuid.New()
	newUserParams := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.Args[0],
	}
	/* We create a new user and also use it as a way to check if a user already exist. if a
	user already exists(duplicate records), we would get an error value
	*/
	newUser, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DEBUG: Error creating user: %v\n", err)
		return fmt.Errorf("failed to create new user. User might already exist %v", err)
	}
	/* We then set the username of the config to the name field of the newUser using the SetUser
	function of the state passed to this handler
	*/

	if err := s.config.SetUser(newUser.Name); err != nil {
		return fmt.Errorf("could not set current user %v", err)
	}
	// If all the checks have been successful, we print to the console that it was a success
	fmt.Println("registered new user successfully")
	return nil
}

func handlerReset(s *state, cmd command) error {
	/*your functions can call a query to the database using the queryname in the .sql file
	e.g s.db.ResetUsers()
	*/
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
	for _, username := range users {
		if s.config.CurrentUserName == username {
			fmt.Printf("%v (current)", username)
		} else {
			fmt.Println(username)
		}

	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("unable to run agg command %v", err)
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

//---------------------------Handlers End---------------------------------------------------//

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

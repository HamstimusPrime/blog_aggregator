package main

import (
	"context"

	"github.com/hamstimusprime/blog_aggregator/internal/database"
)

/*
this middleware accepts gets info from the database about the currently
loggedIn user and then passes that user info(users) to a handler function which it
then returns. the returned handler function contaning data about the user is used
in handler functions that need to use data of the currently logged in user
*/
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

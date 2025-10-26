package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/componhead/gator/internal/config"
	"github.com/componhead/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	var username string
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username = cmd.Args[0]
	ctx := context.Background()
	u, err := s.db.GetUserByName(ctx, username)
	if (err != nil || u == database.User{}) {
		return fmt.Errorf("user %s doesn't exist", username)
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Printf("Username %s setted\n", username)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	var name string
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name = cmd.Args[0]
	ctx := context.Background()
	u, err := s.db.GetUserByName(ctx, name)
	if u.ID.Valid {
		return fmt.Errorf("user %s already exists", name)
	}
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("couldn't get user %s: %w", name, err)
	}
	userParams := database.CreateUserParams{
		ID:        uuid.NullUUID{UUID: uuid.New(), Valid: true},
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdUser, err := s.db.CreateUser(ctx, userParams)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	s.cfg.SetUser(createdUser.Name)
	fmt.Printf("User created: %+v\n", createdUser)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't truncate user table: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	usrs, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}
	var userRows []string
	for _, u := range usrs {
		var userRow []string
		userRow = append(userRow, fmt.Sprintf("* %s", u.Name))
		if u.Name == s.cfg.CurrentUserName {
			userRow = append(userRow, " (current)")
		}
		userRows = append(userRows, strings.Join(userRow, ""))
	}
	for _, userRow := range userRows {
		fmt.Println(userRow)
	}
	return nil
}

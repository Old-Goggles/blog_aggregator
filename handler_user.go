package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Old-Goggles/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Args[0]
	ctx := context.Background()
	id := uuid.New()
	now := time.Now().UTC()
	params := database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      username,
	}

	result, err := s.Db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("user name already exists: %w", err)
	}

	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("User Created: %+v", result)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}

	ctx := context.Background()
	username := cmd.Args[0]

	_, err := s.Db.GetUser(ctx, username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user name does not exist")
	} else if err != nil {
		return fmt.Errorf("error finding user in database %w", err)
	}

	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Println("user has been set to", username)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()

	result, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("unable to get users")
	}

	for _, user := range result {
		if user.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

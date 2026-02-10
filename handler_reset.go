package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.Db.DeleteUsers(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete: %w", err)
	}

	fmt.Printf("users reset")
	return nil
}

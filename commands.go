package main

import (
	"fmt"

	"github.com/Old-Goggles/blog_aggregator/internal/config"
	"github.com/Old-Goggles/blog_aggregator/internal/database"
)

type state struct {
	Db  *database.Queries
	Cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*state, command) error)
	}
	c.Handlers[name] = f
}

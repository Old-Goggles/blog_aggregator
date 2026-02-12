package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Old-Goggles/blog_aggregator/internal/config"
	"github.com/Old-Goggles/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		Db:  dbQueries,
		Cfg: &cfg,
	}

	cmds := &commands{
		Handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", handlerBrowse)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		username := s.Cfg.CurrentUserName
		user, err := s.Db.GetUser(ctx, username)
		if err != nil {
			return fmt.Errorf("error finding user in database %w", err)
		}
		return handler(s, cmd, user)
	}
}

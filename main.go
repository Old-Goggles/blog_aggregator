package main

import (
	"fmt"

	"github.com/Old-Goggles/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error reading config:", err)
		return
	}

	err = cfg.SetUser("old_goggles")
	if err != nil {
		fmt.Println("error setting user:", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("error reading config again:", err)
		return
	}

	fmt.Println(cfg)
}

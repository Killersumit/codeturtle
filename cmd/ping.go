package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/killersumit/codeturtle/internal/config"
	"github.com/killersumit/codeturtle/internal/provider"
)

// Ping verifies provider connectivity by sending a test message.
func Ping(args []string) {
	fs := flag.NewFlagSet("ping", flag.ExitOnError)
	configPath := fs.String("config", ".codeturtle/config.toml", "Path to config file")
	_ = fs.Parse(args)

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	p, err := provider.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing provider: %v\n", err)
		os.Exit(1)
	}

	messages := []provider.Message{
		{
			Role:    "user",
			Content: "Say hello and tell me what model you are in one sentence.",
		},
	}

	reply, err := p.Chat(context.Background(), messages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error contacting provider: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(reply)
	fmt.Println("🐢 CodeTurtle ping OK")
}

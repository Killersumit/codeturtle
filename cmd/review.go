package cmd

import (
	"flag"
	"fmt"
)

func Review(args []string) {
	fs := flag.NewFlagSet("review", flag.ExitOnError)
	base := fs.String("base", "main", "Base branch to compare against")
	config := fs.String("config", ".codeturtle/config.toml", "Path to config file")
	_ = fs.Parse(args)

	fmt.Println("🐢 CodeTurtle review")
	fmt.Printf("   Base branch: %s\n", *base)
	fmt.Printf("   Config:      %s\n", *config)
	fmt.Println("   (Provider not connected yet — Phase 1)")
}
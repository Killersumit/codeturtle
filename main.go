package main

import (
	"fmt"
	"os"

	"github.com/killersumit/codeturtle/cmd"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "--version", "-V", "version":
		fmt.Printf("codeturtle %s\n", version)
	case "--help", "-h", "help":
		printHelp()
	case "review":
		cmd.Review(os.Args[2:])
	case "init":
		cmd.Init(os.Args[2:])
	case "config":
		cmd.Config(os.Args[2:])
	case "ping":
		cmd.Ping(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("🐢 CodeTurtle — a shell for your code reviews")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("    codeturtle <COMMAND> [options]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("    review     Run a code review on the current git diff")
	fmt.Println("    init       Initialize a new .codeturtle config")
	fmt.Println("    config     Show current configuration")
	fmt.Println("    ping       Test connection to configured LLM provider")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("    -h, --help       Print help")
	fmt.Println("    -V, --version    Print version")
}
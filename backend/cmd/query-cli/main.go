package main

import (
	"fmt"
	"os"
	"stafind-backend/internal/queries"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Query Manager CLI")
		fmt.Println("Usage: go run cmd/query-cli/main.go <command> [args]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  list                    - List all queries")
		fmt.Println("  show <name>             - Show detailed query information")
		fmt.Println("  category                - List queries by category")
		fmt.Println("  tag <tag>               - List queries by tag")
		fmt.Println("  validate                - Validate all queries")
		fmt.Println("  export                  - Export all queries")
		fmt.Println("  stats                   - Show configuration statistics")
		fmt.Println("  interactive             - Start interactive mode")
		os.Exit(1)
	}

	cli, err := queries.NewQueryCLI()
	if err != nil {
		fmt.Printf("Error initializing CLI: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "list":
		cli.ListQueries()
	case "show":
		if len(args) == 0 {
			fmt.Println("Usage: show <query_name>")
			os.Exit(1)
		}
		if err := cli.ShowQuery(args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "category":
		cli.ListByCategory()
	case "tag":
		if len(args) == 0 {
			fmt.Println("Usage: tag <tag_name>")
			os.Exit(1)
		}
		cli.ListByTag(args[0])
	case "validate":
		if err := cli.ValidateQueries(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "export":
		cli.ExportQueries()
	case "stats":
		cli.ShowStats()
	case "interactive":
		cli.RunCLI()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

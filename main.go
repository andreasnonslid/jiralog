package main

import (
	"fmt"
	"os"
)

// main is the application entry point, handling CLI arguments and routing to commands.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./jiralog <command> [date-range] [options]")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	handleCommand(command, args)
}

// handleCommand processes CLI input to determine the date range and invoke the correct function.
func handleCommand(command string, args []string) {
	// Default range is last 6 months if not specified
	start, end := parseDateRange("6m")
	if len(args) > 0 {
		start, end = parseDateRange(args[0])
	}

	switch command {
	case "recent":
		listRecentIssues(start, end)
	case "view":
		viewDetailedEntries(start, end)
	case "summary":
		listSummary(start, end)
	case "viewraw":
		viewRawEntries(start, end)
	default:
		fmt.Println("Unknown command. Available commands: recent, view, summary.")
	}
}

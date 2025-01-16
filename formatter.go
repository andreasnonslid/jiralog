package main

import (
	"fmt"
)

func formatEntriesGrouped(entries []map[string]interface{}) {
	grouped := make(map[string][]map[string]interface{})

	// Group entries by issue key
	for _, e := range entries {
		issueKey, _ := e["issue"].(string)
		grouped[issueKey] = append(grouped[issueKey], e)
	}

	fmt.Println("Detailed Worklog Entries (Grouped by Issue):")
	for issueKey, issueEntries := range grouped {
		fmt.Printf("Issue: %s\n", issueKey)
		for _, e := range issueEntries {
			started := e["started"]
			hours := e["hours"]

			fmt.Printf("  - %s | %s\n", started, hours)
		}
		fmt.Println()
	}
}

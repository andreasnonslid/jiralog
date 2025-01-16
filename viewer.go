package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// listRecentIssues fetches and prints the issues a user has recently worked on in the given date range.
func listRecentIssues(start, end time.Time) {
	config, apiToken := loadConfigAndToken()

	url := fmt.Sprintf("https://%s/rest/api/3/search", config.CloudSpace)
	query := map[string]interface{}{
		"jql": fmt.Sprintf(
			"worklogAuthor = \"%s\" AND worklogDate >= %s AND worklogDate <= %s",
			config.Email, start.Format("2006-01-02"), end.Format("2006-01-02"),
		),
		"fields": []string{"summary"},
	}

	payload, err := json.Marshal(query)
	handleError(err, "Failed to marshal listRecentIssues query")

	body := makeAPIRequest("POST", url, payload, config, apiToken)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	handleError(err, "Failed to parse issues")

	fmt.Println("Recent Issues:")
	if issues, ok := result["issues"].([]interface{}); ok {
		for _, issue := range issues {
			issueData := issue.(map[string]interface{})
			key := issueData["key"]
			summary := issueData["fields"].(map[string]interface{})["summary"]
			fmt.Printf("%s - %s\n", key, summary)
		}
	}
}

// viewDetailedEntries fetches detailed worklogs and prints them either chronologically or grouped by issue.
func viewDetailedEntries(start, end time.Time) {
	config, apiToken := loadConfigAndToken()

	url := fmt.Sprintf("https://%s/rest/api/3/search", config.CloudSpace)
	query := map[string]interface{}{
		"jql": fmt.Sprintf("worklogAuthor = \"%s\" AND worklogDate >= %s AND worklogDate <= %s",
			config.Email, start.Format("2006-01-02"), end.Format("2006-01-02")),
		"fields": []string{"summary", "worklog"},
		"expand": []string{"worklog"},
	}

	payload, err := json.Marshal(query)
	handleError(err, "Failed to marshal viewDetailedEntries query")

	body := makeAPIRequest("POST", url, payload, config, apiToken)

	entries, err := parseWorklogs(body)
	handleError(err, "Failed to parse worklogs")

	formatEntriesGrouped(entries)
}

// viewRawEntries fetches and prints the raw JSON response from Jira's API.
func viewRawEntries(start, end time.Time) {
	config, apiToken := loadConfigAndToken()

	// Construct API URL and query
	url := fmt.Sprintf("https://%s/rest/api/3/search", config.CloudSpace)
	query := map[string]interface{}{
		"jql": fmt.Sprintf("worklogAuthor = \"%s\" AND worklogDate >= %s AND worklogDate <= %s",
			config.Email, start.Format("2006-01-02"), end.Format("2006-01-02")),
		"fields": []string{"summary", "worklog"},
		"expand": []string{"worklog"},
	}

	// Marshal query into JSON payload
	payload, err := json.Marshal(query)
	handleError(err, "Failed to marshal viewRawEntries query")

	// Make API request
	body := makeAPIRequest("POST", url, payload, config, apiToken)

	// Print the entire raw response
	fmt.Println("Raw Response:")
	fmt.Println(string(body))
fmt.Printf("Constructed JQL: worklogAuthor = \"%s\" AND worklogDate >= %s AND worklogDate <= %s\n", config.Email, start.Format("2006-01-02"), end.Format("2006-01-02"))
}

// listSummary sums time spent per issue over a date range.
func listSummary(start, end time.Time) {
	config, apiToken := loadConfigAndToken()

	url := fmt.Sprintf("https://%s/rest/api/3/search", config.CloudSpace)
	query := map[string]interface{}{
		"jql": fmt.Sprintf(
			"worklogAuthor = \"%s\" AND worklogDate >= %s AND worklogDate <= %s",
			config.Email, start.Format("2006-01-02"), end.Format("2006-01-02"),
		),
		"fields": []string{"summary", "worklog"},
		"expand": []string{"worklog"},
	}

	payload, err := json.Marshal(query)
	handleError(err, "Failed to marshal listSummary query")

	body := makeAPIRequest("POST", url, payload, config, apiToken)

	entries, err := parseWorklogs(body)
	handleError(err, "Failed to parse worklog summary")

	fmt.Println("Summary:")
	issueMap := make(map[string]float64)
	for _, entry := range entries {
		issue := entry["issue"].(string)
		timeSpent := parseTimeSpent(entry["hours"].(string))
		issueMap[issue] += timeSpent
	}

	for issue, totalHours := range issueMap {
		fmt.Printf("%s: %.2f hours\n", issue, totalHours)
	}
}

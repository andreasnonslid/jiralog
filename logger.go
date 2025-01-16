package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// logWork adds a worklog entry to a Jira issue.
func logWork(issueKey, timeSpent, startDate, comment string) {
	config, apiToken := loadConfigAndToken()

	url := fmt.Sprintf("https://%s/rest/api/3/issue/%s/worklog", config.CloudSpace, issueKey)
	worklog := map[string]interface{}{
		"timeSpent": timeSpent,
		"started":   fmt.Sprintf("%sT09:00:00.000+0000", startDate),
		"comment":   comment,
	}

	payload, err := json.Marshal(worklog)
	handleError(err, "Failed to create worklog payload")

	req := createRequest("POST", url, payload, config.Email, apiToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	handleResponse(resp, err, "Failed to log work")

	fmt.Println("Worklog added successfully!")
}

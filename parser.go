package main

import (
	"encoding/json"
	"fmt"
)

// parseWorklogs takes the raw JSON from a Jira search response and extracts a list of worklog entries.
func parseWorklogs(body []byte) ([]map[string]interface{}, error) {
	// Mirror Jira's "search" response structure
	var response struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary string `json:"summary"`
				Worklog struct {
					Worklogs []struct {
						TimeSpent string      `json:"timeSpent"`
						Started   string      `json:"started"`
						Comment   interface{} `json:"comment"`
					} `json:"worklogs"`
				} `json:"worklog"`
			} `json:"fields"`
		} `json:"issues"`
	}

	err := json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var entries []map[string]interface{}
	for _, issue := range response.Issues {
		for _, wl := range issue.Fields.Worklog.Worklogs {
			// Convert comment to string, even if it might be structured
			commentStr := fmt.Sprintf("%v", wl.Comment)

			entry := map[string]interface{}{
				"issue":   issue.Key,
				"hours":   wl.TimeSpent,
				"started": wl.Started,
				"comment": commentStr,
				"summary": issue.Fields.Summary,
			}
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

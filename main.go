package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type SearchRequest struct {
	JQL    string   `json:"jql"`
	Fields []string `json:"fields"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

type SearchResponse struct {
	Issues []Issue `json:"issues"`
}

func getAuthHeader() (string, error) {
	email := os.Getenv("JIRA_EMAIL")
	token := os.Getenv("JIRA_TOKEN")
	if email == "" || token == "" {
		return "", fmt.Errorf("JIRA_EMAIL or JIRA_TOKEN not set")
	}
	creds := fmt.Sprintf("%s:%s", email, token)
	enc := base64.StdEncoding.EncodeToString([]byte(creds))
	return "Basic " + enc, nil
}

func jiraRequest(method, path string, body io.Reader) ([]byte, error) {
	host := os.Getenv("JIRA_HOST")
	if host == "" {
		host = "autostore.atlassian.net"
	}
	url := fmt.Sprintf("https://%s%s", host, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	header, err := getAuthHeader()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", header)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("jira http error: %s", resp.Status)
	}
	return data, nil
}

func searchIssues(jql string, fields []string) error {
	reqBody := SearchRequest{JQL: jql, Fields: fields}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	data, err := jiraRequest(http.MethodPost, "/rest/api/3/search", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	var resp SearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	for _, issue := range resp.Issues {
		fmt.Printf("%s : %s\n", issue.Key, issue.Fields.Summary)
	}
	return nil
}

func logWork(issue, started string, seconds int, comment string) error {
	worklog := map[string]any{
		"comment": map[string]any{
			"type":    "doc",
			"version": 1,
			"content": []map[string]any{
				{
					"type":    "paragraph",
					"content": []map[string]any{{"type": "text", "text": comment}},
				},
			},
		},
		"started":          started,
		"timeSpentSeconds": seconds,
		"adjustEstimate":   "auto",
	}
	bodyBytes, err := json.Marshal(worklog)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/rest/api/3/issue/%s/worklog?notifyUsers=false", issue)
	_, err = jiraRequest(http.MethodPost, path, bytes.NewReader(bodyBytes))
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: jiralog [search|log] [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "search":
		searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
		jql := searchCmd.String("jql", "", "JQL query")
		fieldsStr := searchCmd.String("fields", "summary", "comma separated fields")
		searchCmd.Parse(os.Args[2:])
		if *jql == "" {
			log.Fatal("-jql is required")
		}
		fields := []string{"summary"}
		if *fieldsStr != "" {
			fields = []string{}
			for _, f := range splitComma(*fieldsStr) {
				fields = append(fields, f)
			}
		}
		if err := searchIssues(*jql, fields); err != nil {
			log.Fatal(err)
		}
	case "log":
		logCmd := flag.NewFlagSet("log", flag.ExitOnError)
		issue := logCmd.String("issue", "", "Issue key")
		started := logCmd.String("started", "", "start time RFC3339")
		seconds := logCmd.Int("seconds", 0, "time spent in seconds")
		comment := logCmd.String("comment", "", "comment text")
		logCmd.Parse(os.Args[2:])
		if *issue == "" || *started == "" || *seconds == 0 {
			log.Fatal("issue, started and seconds are required")
		}
		if err := logWork(*issue, *started, *seconds, *comment); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}

func splitComma(s string) []string {
	var parts []string
	for _, p := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

// handleError prints an error and exits if err is not nil.
func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
}

// handleResponse checks HTTP status and reads the response body if an error is detected.
func handleResponse(resp *http.Response, err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s: %s\n", message, body)
		os.Exit(1)
	}
}

// createRequest constructs an HTTP request with basic auth.
func createRequest(method, url string, payload []byte, email, apiToken string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	handleError(err, "Failed to create HTTP request")
	req.SetBasicAuth(email, apiToken)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// parseDateRange converts strings like "6m", "1w" into a start and end time.
func parseDateRange(rangeStr string) (time.Time, time.Time) {
	now := time.Now()
	re := regexp.MustCompile(`(\d+)([wdhm])`)
	matches := re.FindAllStringSubmatch(rangeStr, -1)

	var duration time.Duration
	for _, match := range matches {
		value, _ := strconv.Atoi(match[1])
		unit := match[2]

		switch unit {
		case "w":
			duration += time.Hour * 24 * 7 * time.Duration(value)
		case "d":
			duration += time.Hour * 24 * time.Duration(value)
		case "h":
			duration += time.Hour * time.Duration(value)
		case "m": // Approximate 1m = 30 days
			duration += time.Hour * 24 * 30 * time.Duration(value)
		}
	}

	start := now.Add(-duration)
	return start, now
}

// parseTimeSpent converts a string like "1h 30m" into a float64 in hours.
func parseTimeSpent(timeStr string) float64 {
	var totalMinutes float64

	re := regexp.MustCompile(`(\d+)([hm])`)
	matches := re.FindAllStringSubmatch(timeStr, -1)

	for _, match := range matches {
		value, _ := strconv.Atoi(match[1])
		unit := match[2]
		switch unit {
		case "h":
			totalMinutes += float64(value * 60)
		case "m":
			totalMinutes += float64(value)
		}
	}
	return totalMinutes / 60.0
}

// makeAPIRequest wraps an HTTP request to the Jira API.
func makeAPIRequest(method, url string, payload []byte, config Config, apiToken string) []byte {
	req := createRequest(method, url, payload, config.Email, apiToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	handleResponse(resp, err, "Request error")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	handleError(err, "Failed to read response body")
	return body
}

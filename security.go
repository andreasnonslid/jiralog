package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

// Config holds the necessary domain and email info to connect to Jira.
type Config struct {
	CloudSpace string `json:"cloud_space"`
	Email      string `json:"email"`
}

const service = "jira-session"

// setup is called to configure Jira domain, email, and secure token.
func setup() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter your Jira Cloud domain (e.g., company.atlassian.net):")
	fmt.Print("Cloud Space: ")
	cloudSpace, _ := reader.ReadString('\n')
	cloudSpace = strings.TrimSpace(cloudSpace)

	fmt.Println("Enter your Atlassian account email:")
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Println("Generate an API token from your Atlassian account here:")
	fmt.Println("https://id.atlassian.com/manage-profile/security/api-tokens")
	fmt.Print("Paste your API token here (input will be hidden): ")
	apiToken := readHiddenInput()

	config := Config{
		CloudSpace: cloudSpace,
		Email:      email,
	}
	saveConfig(config)

	err := keyring.Set(service, email, apiToken)
	handleError(err, "Failed to securely store API token")

	fmt.Println("Setup complete!")
}

// loadConfigAndToken returns the Config and API token, or exits on error.
func loadConfigAndToken() (Config, string) {
	config, err := loadConfig()
	handleError(err, "Failed to load configuration. Please run 'setup' first.")

	apiToken, err := keyring.Get(service, config.Email)
	handleError(err, "Failed to retrieve API token. Please run 'setup' again.")

	return config, apiToken
}

// saveConfig writes config to a JSON file in the user’s home directory.
func saveConfig(config Config) {
	usr, _ := user.Current()
	configPath := filepath.Join(usr.HomeDir, ".jira_config.json")
	data, _ := json.MarshalIndent(config, "", "  ")
	_ = os.WriteFile(configPath, data, 0600)
}

// loadConfig reads config from the home directory JSON file.
func loadConfig() (Config, error) {
	usr, _ := user.Current()
	configPath := filepath.Join(usr.HomeDir, ".jira_config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	json.Unmarshal(data, &config)
	return config, nil
}

// readHiddenInput reads keyboard input without echoing characters to the terminal.
func readHiddenInput() string {
	input, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return strings.TrimSpace(string(input))
}

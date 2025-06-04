# jiralog

A small CLI tool to search issues and log work hours in Jira using Go.

## Setup

1. Generate an API token from [Atlassian](https://id.atlassian.com/manage-profile/security/api-tokens).
2. Set the following environment variables:
   - `JIRA_EMAIL` – your Jira account email.
   - `JIRA_TOKEN` – the API token.
   - `JIRA_HOST`  – the Jira host (default: `autostore.atlassian.net`).

## Usage

Run the commands from the repository root:

### Search issues

```
go run . search -jql "project = TIME" -fields summary
```

### Log work

```
go run . log -issue TIME-25 -started 2025-01-23T12:30:00Z -seconds 3600 -comment "worked on stuff"
```

## Development

Run the Go tests with:

```
go test
```

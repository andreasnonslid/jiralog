# jiralog

An OCaml CLI to search Jira issues and log work hours.

## Setup

1. Generate an API token from [Atlassian](https://id.atlassian.com/manage-profile/security/api-tokens).
2. Set environment variables:
   - `JIRA_EMAIL` – your Jira account email.
   - `JIRA_TOKEN` – the API token.
   - `JIRA_HOST`  – Jira host (default `autostore.atlassian.net`).
3. Install dependencies using `opam`:

```bash
opam install cohttp-lwt-unix yojson cmdliner base64 dune
```

## Building

Use `dune` to build the executable:

```bash
dune build
```

The binary will appear at `_build/default/src/jiralog.exe` or can be installed with `dune install`.

## Usage

Run commands from the repo root:

### Search issues

```bash
./_build/default/src/jiralog.exe search --jql "project = TIME" --fields summary
```

### Log work

```bash
./_build/default/src/jiralog.exe log --issue TIME-25 --started 2025-01-23T12:30:00Z --seconds 3600 --comment "worked on stuff"
```

## Development

You can run the OCaml code directly with `dune exec`:

```bash
dune exec jiralog -- search --jql "project = TIME" --fields summary
```

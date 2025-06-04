# jiralog
A minimal OCaml CLI for searching Jira issues and logging work.  It performs the
same requests as the example `curl` commands in the Atlassian API.

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

The binary will appear at `_build/default/bin/jiralog.exe` or can be installed with
`dune install`.

## Usage

Run commands from the repo root:

### Search issues

```bash
./_build/default/bin/jiralog.exe search --jql "project = TIME" --fields summary
# Pipe through `fzf` and copy the selected issue key to the clipboard
# ./_build/default/bin/jiralog.exe search --jql "project = TIME" --fields summary | fzf | awk '{print $1}' | xclip -sel clip
```

### Log work

```bash
./_build/default/bin/jiralog.exe log --issue TIME-25 --started 2025-01-23T12:30:00Z --seconds 3600 --comment "worked on stuff"
```

## Development

You can run the OCaml code directly with `dune exec`:

```bash
dune exec jiralog -- search --jql "project = TIME" --fields summary
```

Run the test suite with:

```bash
dune runtest
```

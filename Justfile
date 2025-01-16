build:
    go build -o jiralog main.go logger.go viewer.go security.go utility.go parser.go formatter.go

setup:
    just build && ./jiralog

log *args:
    just build && ./jiralog log {{ args }}

view *args:
    just build && ./jiralog view {{ args }}

recent *args:
    just build && ./jiralog recent {{ args }}

summary *args:
    just build && ./jiralog summary {{ args }}

raw *args:
    just build && ./jiralog viewraw {{ args }}

_cat_files:
    cat main.go
    cat logger.go
    cat viewer.go
    cat security.go
    cat utility.go
    cat parser.go
    cat formatter.go

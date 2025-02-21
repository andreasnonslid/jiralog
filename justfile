email email="":
    source env/bin/activate && ./app view {{ email }}

domain domain="":
    source env/bin/activate && ./app domain {{ domain }}

oauth:
    source env/bin/activate && ./app oauth

log *args:
    source env/bin/activate && ./app log {{args}}

view:
    source env/bin/activate && ./app view

recent:
    source env/bin/activate && ./app recent

summary:
    source env/bin/activate && ./app summary

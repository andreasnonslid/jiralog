#!/usr/bin/env python3

import sys
import subprocess
import requests
import json
import keyring
import getpass
from datetime import datetime, timedelta, date
import pyperclip

def usage():
    return """Usage:
  ./app oauth                  # Prompt for and securely store OAuth token
  ./app email <email>          # Save your email
  ./app domain <domain>        # Save the Jira domain (e.g. https://company.atlassian.net)
  ./app view                   # View the TIME project issues in fzf and copy the selected key
  ./app log <timecode> <date> <hours> [comment]

Examples:
  ./app oauth
  ./app email you@example.com
  ./app domain https://company.atlassian.net
  ./app view
  ./app log TIME-25 2025-01-23 2 "Finished some tasks"
"""

# -------------------------------------------------------------------------
# Functions to store/retrieve settings from the keyring
# -------------------------------------------------------------------------
def set_oauth_token():
    print("Get API token at: https://id.atlassian.com/manage-profile/security/api-tokens")
    token = getpass.getpass("Enter your OAuth token: ")
    keyring.set_password("jira", "token", token)
    print("OAuth token stored securely.")

def get_oauth_token():
    return keyring.get_password("jira", "token")

def set_email(email):
    keyring.set_password("jira", "email", email)
    print(f"Email stored: {email}")

def get_email():
    return keyring.get_password("jira", "email")

def set_domain(domain):
    keyring.set_password("jira", "domain", domain)
    print(f"Domain stored: {domain}")

def get_domain():
    domain = keyring.get_password("jira", "domain")
    if domain and not domain.startswith(("http://", "https://")):
        domain = "https://" + domain
    return domain

# -------------------------------------------------------------------------
# View issues from "TIME" project and copy selected key to clipboard
# -------------------------------------------------------------------------
def view_issues():
    email = get_email()
    token = get_oauth_token()
    domain = get_domain()

    if not (email and token and domain):
        print("Error: Missing email/token/domain. Please set them before viewing issues.")
        return

    url = f"{domain}/rest/api/3/search"
    headers = {"Content-Type": "application/json"}
    jql = "project = TIME ORDER BY key ASC"
    max_results = 100
    start_at = 0
    all_issues = []

    while True:
        payload = {
            "jql": jql,
            "fields": ["summary"],
            "startAt": start_at,
            "maxResults": max_results
        }
        response = requests.post(url, headers=headers, json=payload, auth=(email, token))
        response.raise_for_status()
        data = response.json()
        issues = data.get("issues", [])
        all_issues.extend(issues)
        
        if len(issues) < max_results:
            break  # No more issues to fetch
        
        start_at += max_results

    if not all_issues:
        print("No issues found.")
        return

    lines = [f"{issue['key']} : {issue['fields'].get('summary', '')}" for issue in all_issues]

    # Use fzf to select an issue
    result = subprocess.run(["fzf"], input="\n".join(lines), text=True, capture_output=True)
    selected = result.stdout.strip()
    if not selected:
        print("No selection made.")
        return

    # Extract key and copy to clipboard using pyperclip
    selected_key = selected.split()[0]
    pyperclip.copy(selected_key)
    print(f"Copied {selected_key} to clipboard.")

# -------------------------------------------------------------------------
# Log time in a specific JIRA ticket
# -------------------------------------------------------------------------
def log_time(timecode, date_str, hours, comment=None):
    email = get_email()
    token = get_oauth_token()
    domain = get_domain()

    if not (email and token and domain):
        print("Error: Missing email/token/domain. Please set them before logging.")
        return

    # Handle special cases for date_str
    if date_str.lower() == "today":
        dt = date.today()
    elif date_str.lower() == "yesterday":
        dt = date.today() - timedelta(days=1)
    else:
        try:
            dt = datetime.strptime(date_str, "%Y-%m-%d").date()
        except ValueError:
            print("Error: Date must be 'YYYY-MM-DD', 'today', or 'yesterday'.")
            return

    # Use a default time of 09:00:00; adjust as needed.
    started = dt.strftime("%Y-%m-%dT09:00:00.000+0000")

    # Allow hours to be a float (e.g., 1.5 hours)
    try:
        hours_float = float(hours)
    except ValueError:
        print("Error: Hours must be a number (float or int).")
        return
    time_spent_seconds = int(hours_float * 3600)

    url = f"{domain}/rest/api/3/issue/{timecode}/worklog?notifyUsers=false"
    headers = {
        "Accept": "application/json",
        "Content-Type": "application/json"
    }

    # Build payload without comment if not provided.
    payload = {
        "started": started,
        "timeSpentSeconds": time_spent_seconds,
        "adjustEstimate": "auto"
    }
    if comment:
        payload["comment"] = {
            "type": "doc",
            "version": 1,
            "content": [
                {
                    "type": "paragraph",
                    "content": [
                        {
                            "type": "text",
                            "text": comment
                        }
                    ]
                }
            ]
        }

    response = requests.post(url, headers=headers, data=json.dumps(payload), auth=(email, token))
    try:
        response.raise_for_status()
    except Exception as e:
        print(f"Error creating worklog: {e}")
        return

    print(f"Worklog created on {timecode} with {hours_float} hour(s) at {dt.isoformat()} (started at {started}).")

# -------------------------------------------------------------------------
# Main entry
# -------------------------------------------------------------------------
def main():
    if len(sys.argv) < 2:
        print(usage())
        sys.exit(1)

    command = sys.argv[1].lower()

    if command == "oauth":
        set_oauth_token()

    elif command == "email":
        if len(sys.argv) < 3:
            print(usage())
            sys.exit(1)
        set_email(sys.argv[2])

    elif command == "domain":
        if len(sys.argv) < 3:
            print(usage())
            sys.exit(1)
        set_domain(sys.argv[2])

    elif command == "view":
        view_issues()

    elif command == "log":
        if len(sys.argv) < 5:
            print(usage())
            sys.exit(1)
        timecode = sys.argv[2]
        date_str = sys.argv[3]
        hours = sys.argv[4]
        comment = " ".join(sys.argv[5:]) if len(sys.argv) > 5 else None
        log_time(timecode, date_str, hours, comment)

    else:
        print(usage())

if __name__ == "__main__":
    main()

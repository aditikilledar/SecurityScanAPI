# Repository Scan API

## Overview

This is my submission for the take-home test. Authored by me, Aditi Killedar in February 2025.
The Security Scan API is a web service that allows users to query and retrieve information about security vulnerabilities. The API provides endpoints to filter and retrieve vulnerability data based on specific criteria, such as severity.

## Usage Instructions

### Using the service

#### Managing the container
##### Start the service with:
```
docker compose up
```

##### Stop the service with:
```
docker compose down
```

###### To verify or access the database for testing purposes, attach a shell to the container and then run:
```
sqlite3 ./scans.db
```

#### Sending requests to the service

On the terminal, use the ``` curl ``` command to make HTTP requests to the service, or alternatively use Postman.

Here's an example for a request to the endpoint ```\scan```:
```
curl -X POST \
  http://localhost:8080/scan \
  -H 'Content-Type: application/json' \
  -d '{
    "repo": "https://github.com/velancio/vulnerability_scans",
    "files": [
        "vulnscan15.json"
    ]
}'
```

## Testing Instructions

To test, I have written test cases 


## Features

This service lets you:
- Scan a GitHub repository for a set of requested files.
- Query scanned payloads, filter based on some attributes.

## Endpoints

### 1. POST /scan

Scans a (public, root) repository for all json files mentioned in the request and stores its payloads into an sqlite3 database.

Stores metadata like file name and time scanned along with the payloads in the file.

#### Request Format

``` json
curl -X POST \
  http://localhost:8080/scan \
  -H 'Content-Type: application/json' \
  -d '{
    "repo": "https://github.com/<repo_owner_username>/<repository_root>",
    "files": [
        "<filename1.json>",
        "<filename2.json>", ...
    ]
}'
```

#### Response Format

``` json
{
    "message":"Scan completed. Stored files successfully",
    "timestamp":<Datetime>
}
```

#### Example Request Body

```json
{
    "repo": "https://github.com/velancio/vulnerability_scans",
    "files": [
        "vulnscan15.json",
        "vulnscan16.json",
    ]
}
```


### 2. POST /query

Retrieves vulnerabilities based on the provided filters. This currently only implements filtering based on security. 

If no severity filter is provided, then all the records scanned will be returned.

#### Request Body Format

```
{
    "filters": {"severity": <severity>}
}
```

#### Example Request

```json
{
    "filters": {
        "severity": "HIGH"
    }
}
```

#### Example Response

```json
[
    {
        "id": "CVE-2024-1234",
        "severity": "HIGH",
        "cvss": 8.5,
        "status": "fixed",
        "package_name": "openssl",
        "current_version": "1.1.1t-r0",
        "fixed_version": "1.1.1u-r0",
        "description": "Buffer overflow vulnerability in OpenSSL",
        "published_date": "2024-01-15T00:00:00Z",
        "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-1234",
        "risk_factors": [
            "Remote Code Execution",
            "High CVSS Score",
            "Public Exploit Available"
        ]
    },
    {
        "id": "CVE-2024-8902",
        "severity": "HIGH",
        "cvss": 8.2,
        "status": "fixed",
        "package_name": "openldap",
        "current_version": "2.4.57",
        "fixed_version": "2.4.58",
        "description": "Authentication bypass vulnerability in OpenLDAP",
        "published_date": "2024-01-21T00:00:00Z",
        "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-8902",
        "risk_factors": [
            "Authentication Bypass",
            "High CVSS Score"
        ]
    }
]
```

## Assumptions Made

1. Assumes that the payload is the data inside 'vulnerabilities', as in the below snippet.

``` json
"vulnerabilities": [
        {
          "id": "CVE-2024-7701",
          "severity": "CRITICAL",
          "cvss": 9.1,
          "status": "active",
          "package_name": "postgresql",
          "current_version": "13.4",
          "fixed_version": "13.5",
          "description": "Buffer overflow in PostgreSQL database engine",
          "published_date": "2024-01-23T00:00:00Z",
          "link": "https://nvd.nist.gov/vuln/detail/CVE-2024-7701",
          "risk_factors": [
            "Buffer Overflow",
            "Remote Code Execution",
            "Critical CVSS Score"
          ]
        },
    ...]
```

2. Number of attempts for \scan assumed to be 2 attempts in total
```
"Error Handling: Retry failed GitHub API calls (2 attempts)."
I have implemented it as 1 attempt and 1 reattempt.
```
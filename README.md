# Seekr

is a powerful command-line tool written in Go for scanning and searching URLs. It is designed to handle large-scale URL lists efficiently, support concurrent processing, and provide flexible search capabilities with regex patterns. Ideal for penetration testers, bug bounty hunters, and developers looking to automate the search for sensitive information or vulnerabilities.

## Features

- **Regex Support:** Search for patterns in the response body using regular expressions.
- **Case-Insensitive Search:** Enable case-insensitive matching with the `-i` flag.
- **Concurrency:** Process multiple URLs simultaneously using the `-c` flag to specify the number of concurrent workers.
- **Verbose Logging:** Control log verbosity with the `-v` flag for debugging purposes.
- **File Input:** Load search patterns from a file for complex queries.

## Installation

`go install -v github.com/Vulnpire/seekr`

## Basic Usage

Scan URLs from stdin for a specific keyword or regex pattern:

`cat urls.txt | seekr -q "pattern"`

Search for regex patterns in the response body:

`cat urls.txt | gofind -q "user|order|product|api|invoice|account|profile/[0-9]+|[a-fA-F0-9-]{36}"`

Enable case-insensitive matching:

`cat urls.txt | gofind -q "pattern" -i`

Load multiple patterns from a file:

`cat urls.txt | gofind -qf query_file.txt`

Use concurrent workers to speed up processing:

`cat urls.txt | gofind -q "pattern" -c 20`

## Flags

    -q string : Keyword or regex pattern to search for.
    -qf string : File containing keywords or regex patterns to search for.
    -c int : Number of concurrent workers (default is 1).
    -i : Enable case-insensitive matching.
    -v : Enable verbose logging.

## Example PoC:

Finding Google Maps API keys using RegEx:

![image](https://github.com/user-attachments/assets/a5f296c9-e5c4-4300-8650-ef6f62d9fb2f)

## Example Use Cases:

### API Key Discovery:

Find exposed API keys such as Google API keys (`AIza[0-9A-Za-z-_]{35}`), AWS keys, and OAuth tokens in web responses.

### IDOR Vulnerability Detection:

Search for insecure direct object references using patterns like `user/[0-9]+`, `order/[0-9]+`, or UUIDs (`[a-fA-F0-9-]{36}`).

### Sensitive Data Exposure:

Detect email addresses (`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`), credit card numbers (`\b(?:\d[ -]*?){13,16}\b`), Phone numbers (`\+?[0-9]{1,4}?[-. \(\)]?([0-9]{1,3})?[-. \(\)]?[0-9]{1,4}[-. ]?[0-9]{1,4}[-. ]?[0-9]{1,9}
`), JWT tokens, and session IDs in exposed web data.

### SQL Injection Patterns:

Search for exposed SQL queries (`(select|union|update|delete|insert|drop|alter)\s+(from|into|table)\s+[a-zA-Z0-9_]+
`) that might indicate SQL injection vulnerabilities (e.g., `SELECT, UNION, INSERT, DELETE` statements).

### Directory Traversal and File Path Exposure:

Identify potential directory traversal paths (`\.\./|\.\.\\`), as well as exposed file paths like config.php, index.py, or .sh scripts (`\/([a-zA-Z0-9_\-\.]+\/)*[a-zA-Z0-9_\-\.]+\.(php|html|json|py|sh)
`).

## Disclaimer

This tool is intended for authorized security research, vulnerability testing, and educational purposes only. Unauthorized use of this tool on networks or systems where you do not have explicit permission is illegal and unethical. I do not accept any responsibility or liability for any misuse or damage caused by its use.

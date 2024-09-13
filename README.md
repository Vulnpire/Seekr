# GoFind

is a powerful command-line tool written in Go for scanning and searching URLs. It is designed to handle large-scale URL lists efficiently, support concurrent processing, and provide flexible search capabilities with regex patterns. Ideal for penetration testers, bug bounty hunters, and developers looking to automate the search for sensitive information or vulnerabilities.

## Features

- **Regex Support:** Search for patterns in the response body using regular expressions.
- **Case-Insensitive Search:** Enable case-insensitive matching with the `-i` flag.
- **Concurrency:** Process multiple URLs simultaneously using the `-c` flag to specify the number of concurrent workers.
- **Verbose Logging:** Control log verbosity with the `-v` flag for debugging purposes.
- **File Input:** Load search patterns from a file for complex queries.

## Installation

`go install -v github.com/Vulnpire/gofind`

## Basic Usage

Scan URLs from stdin for a specific keyword or regex pattern:

`cat urls.txt | gofind -q "pattern"`

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

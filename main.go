package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

// Search for regex patterns in the request body of a URL
func searchInRequest(url string, patterns []*regexp.Regexp, verbose bool) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		if verbose {
			fmt.Printf("Error creating request for URL %s: %v\n", url, err)
		}
		return
	}

	// Convert the request object to a string (including headers and body)
	reqDump := fmt.Sprintf("%s %s HTTP/1.1\nHost: %s\n", req.Method, req.URL.Path, req.Host)

	for _, pattern := range patterns {
		matches := pattern.FindAllString(reqDump, -1) // Find all matches in the request dump
		if len(matches) > 0 {
			uniqueMatches := make(map[string]int)
			for _, match := range matches {
				uniqueMatches[match]++
			}
			// Print the summary: number of matches for each unique pattern
			for match, count := range uniqueMatches {
				fmt.Printf("Found %d match(es) of \"%s\" in the request of %s\n", count, match, url)
			}
		}
	}
}

// Search for regex patterns in the response body of a URL
func searchInResponse(url string, patterns []*regexp.Regexp, verbose bool) {
	resp, err := http.Get(url)
	if err != nil {
		if verbose {
			fmt.Printf("Error fetching URL %s: %v\n", url, err)
		}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if verbose {
			fmt.Printf("Error reading response body for %s: %v\n", url, err)
		}
		return
	}

	bodyStr := string(body)

	for _, pattern := range patterns {
		matches := pattern.FindAllString(bodyStr, -1) // Find all matches for the pattern
		if len(matches) > 0 {
			uniqueMatches := make(map[string]int)
			for _, match := range matches {
				uniqueMatches[match]++
			}
			// Print the summary: number of matches for each unique pattern
			for match, count := range uniqueMatches {
				fmt.Printf("Found %d match(es) of \"%s\" in the response of %s\n", count, match, url)
			}
		}
	}
}

// Compile regex patterns from keywords
func compilePatterns(keywords []string, caseInsensitive bool) ([]*regexp.Regexp, error) {
	var patterns []*regexp.Regexp
	for _, keyword := range keywords {
		// Add case-insensitive flag to regex if -i is enabled
		if caseInsensitive {
			keyword = "(?i)" + keyword
		}
		pattern, err := regexp.Compile(keyword)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %s", keyword)
		}
		patterns = append(patterns, pattern)
	}
	return patterns, nil
}

// Load queries from a file
func loadQueriesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var queries []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		queries = append(queries, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return queries, nil
}

// Worker function for concurrency
func worker(urls <-chan string, patterns []*regexp.Regexp, wg *sync.WaitGroup, verbose bool, searchRequest bool) {
	defer wg.Done()
	for url := range urls {
		if searchRequest {
			searchInRequest(url, patterns, verbose)
		} else {
			searchInResponse(url, patterns, verbose)
		}
	}
}

func main() {
	// Define flags for keyword, query file, concurrency, case-insensitive, and verbose logging
	query := flag.String("q", "", "Keyword or regex pattern to search for in the response body")
	queryFile := flag.String("qf", "", "File containing keywords or regex patterns to search for in the response body")
	searchRequest := flag.String("req", "", "Keyword or regex pattern to search for in the request body")
	concurrency := flag.Int("c", 1, "Number of concurrent workers")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	caseInsensitive := flag.Bool("i", false, "Enable case-insensitive matching")
	flag.Parse()

	// Ensure -q and -req flags are not used together
	if *query != "" && *searchRequest != "" {
		fmt.Println("Error: -q and -req flags cannot be used together.")
		os.Exit(1)
	}

	var keywords []string

	// Handle search for response body (-q or -qf)
	if *query != "" || *queryFile != "" {
		// If -q is provided, use it as the search keyword or regex pattern
		if *query != "" {
			keywords = append(keywords, *query)
		}
		// If -qf is provided, load the keywords from the file
		if *queryFile != "" {
			queries, err := loadQueriesFromFile(*queryFile)
			if err != nil {
				fmt.Printf("Error reading query file: %v\n", err)
				os.Exit(1)
			}
			keywords = append(keywords, queries...)
		}
	}

	// Handle search for request body (-req)
	if *searchRequest != "" {
		keywords = append(keywords, *searchRequest)
	}

	// Compile the keywords into regex patterns with case-insensitive option if -i is enabled
	patterns, err := compilePatterns(keywords, *caseInsensitive)
	if err != nil {
		fmt.Printf("Error compiling regex patterns: %v\n", err)
		os.Exit(1)
	}

	// Channel for feeding URLs to workers
	urls := make(chan string)

	// WaitGroup to wait for all workers to complete
	var wg sync.WaitGroup

	// Start workers based on the concurrency flag
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(urls, patterns, &wg, *verbose, *searchRequest != "")
	}

	// Read URLs from stdin (cat urls.txt | gofind)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := scanner.Text()
		urls <- url
	}

	// Close the URL channel and wait for all workers to finish
	close(urls)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading URLs: %v\n", err)
	}
}

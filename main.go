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

// Search for regex patterns in the response body of a URL
func searchInURL(url string, patterns []*regexp.Regexp, verbose bool) {
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
				fmt.Printf("Found %d match(es) of \"%s\" in %s\n", count, match, url)
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
func worker(urls <-chan string, patterns []*regexp.Regexp, wg *sync.WaitGroup, verbose bool) {
	defer wg.Done()
	for url := range urls {
		searchInURL(url, patterns, verbose)
	}
}

func main() {
	// Define flags for keyword, query file, concurrency, case-insensitive, and verbose logging
	query := flag.String("q", "", "Keyword or regex pattern to search for")
	queryFile := flag.String("qf", "", "File containing keywords or regex patterns to search for")
	concurrency := flag.Int("c", 1, "Number of concurrent workers")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	caseInsensitive := flag.Bool("i", false, "Enable case-insensitive matching")
	flag.Parse()

	// Check if either -q or -qf is provided
	if *query == "" && *queryFile == "" {
		fmt.Println("You must specify a keyword/regex (-q) or a query file (-qf)")
		os.Exit(1)
	}

	var keywords []string

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
		go worker(urls, patterns, &wg, *verbose)
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


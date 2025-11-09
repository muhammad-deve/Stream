package parse

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// AnalyzeDomains is a helper function to analyze unique domains in the JSON file
// This can be called separately if needed for debugging
func AnalyzeDomains(filePath string) error {
	// Read JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var streams []StreamEntry
	if err := json.Unmarshal(data, &streams); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Collect unique domain extensions
	domainMap := make(map[string]int)
	totalValid := 0

	for _, stream := range streams {
		// Only analyze entries with all required fields
		if stream.Channel == nil || *stream.Channel == "" {
			continue
		}
		if stream.URL == nil || *stream.URL == "" {
			continue
		}
		if stream.Quality == nil || *stream.Quality == "" {
			continue
		}
		if stream.Title == "" {
			continue
		}

		totalValid++

		// Extract domain extension
		parts := strings.Split(*stream.Channel, ".")
		if len(parts) >= 2 {
			domain := strings.ToLower(parts[len(parts)-1])
			domainMap[domain]++
		}
	}

	// Sort domains by frequency (descending)
	type domainCount struct {
		domain string
		count  int
	}
	var domains []domainCount
	for domain, count := range domainMap {
		domains = append(domains, domainCount{domain, count})
	}
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].count > domains[j].count
	})

	fmt.Printf("Total valid entries: %d\n", totalValid)
	fmt.Printf("Unique domain extensions: %d\n\n", len(domains))
	fmt.Println("Domain | Count | Example Channel")
	fmt.Println("-------+-------+------------------")

	// Find example for each domain
	for _, dc := range domains {
		example := ""
		for _, stream := range streams {
			if stream.Channel != nil && *stream.Channel != "" {
				parts := strings.Split(*stream.Channel, ".")
				if len(parts) >= 2 {
					domain := strings.ToLower(parts[len(parts)-1])
					if domain == dc.domain {
						example = *stream.Channel
						break
					}
				}
			}
		}
	fmt.Printf(".%-5s | %-5d | %s\n", dc.domain, dc.count, example)
	}
	return nil
}

package scrape

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func ScrapeCommand(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "scrape",
		Short: "Scrape actual logo image URLs from webpage URLs and update logo_url fields",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runScrape(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func runScrape(app *pocketbase.PocketBase) error {
	// Get logos collection
	logosCollection, err := app.FindCollectionByNameOrId("logos")
	if err != nil {
		return fmt.Errorf("failed to find logos collection: %w", err)
	}

	// Get all logos (7112 as mentioned)
	allLogos, err := app.FindRecordsByFilter(
		logosCollection.Name,
		"",
		"-created",
		10000, // Increased limit to handle all 7112 logos
		0,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to load logos: %w", err)
	}

	log.Printf("Found %d logos to process\n", len(allLogos))

	updated := 0
	skipped := 0
	failed := 0
	alreadyDirect := 0

	for i, logoRecord := range allLogos {
		logoURL := logoRecord.GetString("logo_url")
		
		// Skip empty URLs
		if logoURL == "" {
			skipped++
			continue
		}

		// Skip if it already points to a direct image (has common image extensions)
		if isDirectImageURL(logoURL) {
			alreadyDirect++
			log.Printf("[%d/%d] Already direct image: %s\n", i+1, len(allLogos), logoURL)
			continue
		}

		log.Printf("[%d/%d] Processing: %s\n", i+1, len(allLogos), logoURL)

		// Scrape the actual image URL from the webpage
		actualImageURL, err := scrapeImageFromPage(logoURL)
		if err != nil {
			log.Printf("Failed to scrape %s: %v\n", logoURL, err)
			failed++
			continue
		}

		// Skip if we couldn't find an image
		if actualImageURL == "" {
			log.Printf("No image found at %s\n", logoURL)
			failed++
			continue
		}

		// Update the logo_url with the direct image URL
		logoRecord.Set("logo_url", actualImageURL)
		if err := app.Save(logoRecord); err != nil {
			log.Printf("Failed to save logo: %v\n", err)
			failed++
			continue
		}

		updated++
		log.Printf("✓ Updated: %s -> %s\n", logoURL, actualImageURL)

		// Rate limiting is handled by Colly
	}

	log.Printf("\n========== Scrape Complete ==========")
	log.Printf("Total processed: %d logos\n", len(allLogos))
	log.Printf("Successfully updated: %d logos\n", updated)
	log.Printf("Already direct images: %d logos\n", alreadyDirect)
	log.Printf("Skipped: %d logos\n", skipped)
	log.Printf("Failed: %d logos\n", failed)

	return nil
}

// isDirectImageURL checks if the URL already points directly to an image file
func isDirectImageURL(url string) bool {
	lowerURL := strings.ToLower(url)
	
	// Special handling for Imgur URLs - these need to be scraped even if they have image extensions
	if strings.Contains(lowerURL, "imgur.com") {
		// Only consider it a direct image if it has the _d.webp pattern (actual direct image)
		// or if it's from i.ibb.co or i.postimg.cc (these are direct)
		if strings.Contains(lowerURL, "_d.webp") {
			return true
		}
		// Basic imgur URLs like https://i.imgur.com/IcWtXCZ.png are NOT direct images
		// They are webpage URLs that need scraping
		return false
	}
	
	// For non-imgur URLs, check for image extensions
	imageExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", 
		".svg", ".tiff", ".ico", ".heic", ".heif",
	}
	
	for _, ext := range imageExtensions {
		if strings.Contains(lowerURL, ext) {
			// Check if it's actually a direct image URL (has extension before query params)
			if idx := strings.Index(lowerURL, "?"); idx != -1 {
				beforeQuery := lowerURL[:idx]
				if strings.HasSuffix(beforeQuery, ext) {
					return true
				}
			} else if strings.HasSuffix(lowerURL, ext) {
				return true
			}
		}
	}
	
	return false
}

// scrapeImageFromPage visits a webpage and extracts the direct image URL
func scrapeImageFromPage(url string) (string, error) {
	// Normalize URL to ensure it's a full URL
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	// Special handling for Imgur URLs
	if strings.Contains(strings.ToLower(url), "imgur.com") {
		return scrapeImgurImage(url)
	}

	var imageURL string
	var scrapeErr error

	// Create Colly collector without domain restrictions
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.Async(false),
	)

	// Add rate limiting: 1 request per second with random delay
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	// Priority 1: Look for img with class="image-placeholder"
	c.OnHTML("img.image-placeholder", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if src != "" && imageURL == "" {
			// Make sure it's a full URL
			if strings.HasPrefix(src, "//") {
				src = "https:" + src
			} else if strings.HasPrefix(src, "/") {
				// Relative URL, construct full URL
				baseURL := e.Request.URL.Scheme + "://" + e.Request.URL.Host
				src = baseURL + src
			}
			imageURL = src
			log.Printf("  Found image-placeholder: %s\n", src)
		}
	})

	// Priority 2: Look for any img with class containing "logo" 
	if imageURL == "" {
		c.OnHTML("img[class*='logo']", func(e *colly.HTMLElement) {
			src := e.Attr("src")
			if src != "" && imageURL == "" && isDirectImageURL(src) {
				// Make sure it's a full URL
				if strings.HasPrefix(src, "//") {
					src = "https:" + src
				} else if strings.HasPrefix(src, "/") {
					baseURL := e.Request.URL.Scheme + "://" + e.Request.URL.Host
					src = baseURL + src
				}
				imageURL = src
				log.Printf("  Found logo class image: %s\n", src)
			}
		})
	}

	// Priority 3: Try meta og:image tag as fallback
	if imageURL == "" {
		c.OnHTML("meta[property='og:image']", func(e *colly.HTMLElement) {
			content := e.Attr("content")
			if content != "" && imageURL == "" {
				// Make sure it's a full URL
				if strings.HasPrefix(content, "//") {
					content = "https:" + content
				} else if strings.HasPrefix(content, "/") {
					baseURL := e.Request.URL.Scheme + "://" + e.Request.URL.Host
					content = baseURL + content
				}
				imageURL = content
				log.Printf("  Found og:image: %s\n", content)
			}
		})
	}

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("scrape error (status %d): %w", r.StatusCode, err)
	})

	// Visit the page
	if err := c.Visit(url); err != nil {
		return "", fmt.Errorf("failed to visit URL: %w", err)
	}

	// Wait for async operations
	c.Wait()

	if scrapeErr != nil {
		return "", scrapeErr
	}

	if imageURL == "" {
		return "", fmt.Errorf("could not find image URL in page")
	}

	return imageURL, nil
}

// scrapeImgurImage handles Imgur URLs specifically
func scrapeImgurImage(url string) (string, error) {
	// Extract the image ID from the URL
	// URLs can be in formats like:
	// https://imgur.com/IcWtXCZ
	// https://i.imgur.com/IcWtXCZ.png
	// https://i.imgur.com/IcWtXCZ.jpg
	
	var imgurID string
	
	// Remove protocol and get the path
	urlLower := strings.ToLower(url)
	
	// Handle different Imgur URL patterns
	if strings.Contains(urlLower, "imgur.com/") {
		parts := strings.Split(url, "imgur.com/")
		if len(parts) >= 2 {
			pathPart := parts[1]
			// Remove any query parameters
			if idx := strings.Index(pathPart, "?"); idx != -1 {
				pathPart = pathPart[:idx]
			}
			// Remove any file extension
			if idx := strings.LastIndex(pathPart, "."); idx != -1 {
				pathPart = pathPart[:idx]
			}
			// Handle album links (skip them)
			if strings.HasPrefix(pathPart, "a/") || strings.HasPrefix(pathPart, "gallery/") {
				return "", fmt.Errorf("album or gallery links not supported")
			}
			imgurID = pathPart
		}
	}
	
	if imgurID == "" {
		return "", fmt.Errorf("could not extract Imgur ID from URL: %s", url)
	}
	
	// Construct the direct image URL using the pattern:
	// https://i.imgur.com/{ID}_d.webp?maxwidth=760&fidelity=grand
	directURL := fmt.Sprintf("https://i.imgur.com/%s_d.webp?maxwidth=760&fidelity=grand", imgurID)
	
	log.Printf("  Converted Imgur URL: %s -> %s\n", url, directURL)
	
	return directURL, nil
}

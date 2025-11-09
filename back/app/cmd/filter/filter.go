package filter

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

type ChannelURL struct {
	ID  string
	URL string
}

type Result struct {
	ChannelID string
	URL       string
	Works     bool
	Reason    string
}

func FilterCommand(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "filter",
		Short: "Validate stream URLs and update channel status",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runFilter(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func runFilter(app *pocketbase.PocketBase) error {
	fmt.Println("📡 Connecting to PocketBase")

	// Get channels collection
	channelsCollection, err := app.FindCollectionByNameOrId("channels")
	if err != nil {
		return fmt.Errorf("failed to find channels collection: %w", err)
	}

	// Fetch all channels
	records, err := app.FindAllRecords(channelsCollection.Name)
	if err != nil {
		return fmt.Errorf("failed to fetch channels: %w", err)
	}

	channels := make([]ChannelURL, 0)
	for _, record := range records {
		url := record.GetString("url")
		if url != "" {
			channels = append(channels, ChannelURL{
				ID:  record.Id,
				URL: url,
			})
		}
	}

	fmt.Printf("🚀 Starting validation of %d channels\n", len(channels))

	workers := 10
	timeoutSec := 8 * time.Second

	fmt.Printf("⚙️  Workers: %d, Timeout: %v\n\n", workers, timeoutSec)

	startTime := time.Now()

	results := processURLsConcurrently(channels, workers, timeoutSec)

	working := 0
	broken := 0

	// Update database with results
	for result := range results {
		// Find the record
		record, err := app.FindRecordById(channelsCollection.Name, result.ChannelID)
		if err != nil {
			fmt.Printf("⚠️  Failed to find channel %s: %v\n", result.ChannelID, err)
			continue
		}

		// Update is_working field
		record.Set("is_working", result.Works)

		if err := app.Save(record); err != nil {
			fmt.Printf("⚠️  Failed to update channel %s: %v\n", result.ChannelID, err)
			continue
		}

		if result.Works {
			fmt.Printf("✅ %s — %s\n", result.URL, result.Reason)
			working++
		} else {
			fmt.Printf("❌ %s — %s\n", result.URL, result.Reason)
			broken++
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n✨ Done in %s!\n", elapsed.Round(time.Second))
	fmt.Printf("📊 Results: ✅ %d working | ❌ %d broken\n", working, broken)

	return nil
}

func processURLsConcurrently(channelURLs []ChannelURL, workers int, timeout time.Duration) <-chan Result {
	results := make(chan Result, len(channelURLs))
	urlChan := make(chan ChannelURL, len(channelURLs))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: timeout,
				Transport: &http.Transport{
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 10,
					IdleConnTimeout:     30 * time.Second,
					DisableKeepAlives:   false,
				},
			}

			for ch := range urlChan {
				works, reason := checkURL(ch.URL, client, timeout)
				results <- Result{
					ChannelID: ch.ID,
					URL:       ch.URL,
					Works:     works,
					Reason:    reason,
				}
			}
		}()
	}

	go func() {
		for _, ch := range channelURLs {
			urlChan <- ch
		}
		close(urlChan)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func checkURL(urlStr string, client *http.Client, timeout time.Duration) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return false, "request error"
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return false, "connection error"
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, fmt.Sprintf("status %d", resp.StatusCode)
	}

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	isPlaylist := strings.Contains(contentType, "mpegurl") ||
		strings.Contains(contentType, "application/vnd.apple.mpegurl") ||
		strings.Contains(urlStr, ".m3u8")

	if isPlaylist {
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
		if err != nil {
			return false, "read error"
		}

		content := string(body)
		if !strings.Contains(content, "#EXTM3U") {
			return false, "invalid m3u8"
		}

		segmentURL := extractFirstSegment(content, urlStr)
		if segmentURL == "" {
			return false, "no segments"
		}

		segmentWorks, _ := checkSegment(segmentURL, client, timeout)
		if !segmentWorks {
			return false, "segments broken"
		}

		return true, "ok"
	}

	validTypes := []string{"video/", "audio/"}
	hasValidType := false
	for _, vt := range validTypes {
		if strings.Contains(contentType, vt) {
			hasValidType = true
			break
		}
	}

	if !hasValidType {
		return false, "invalid type"
	}

	buf := make([]byte, 1024)
	n, err := io.ReadAtLeast(resp.Body, buf, 10)
	if (err != nil && err != io.EOF && err != io.ErrUnexpectedEOF) || n < 10 {
		return false, "no data"
	}

	return true, "ok"
}

func extractFirstSegment(playlist string, baseURL string) string {
	lines := strings.Split(playlist, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasSuffix(line, ".ts") || strings.HasSuffix(line, ".m4s") ||
			strings.HasSuffix(line, ".mp4") || strings.Contains(line, "segment") ||
			strings.HasSuffix(line, ".m3u8") {

			if !strings.HasPrefix(line, "http") {
				baseURLParsed, err := url.Parse(baseURL)
				if err != nil {
					return ""
				}
				segmentURL, err := baseURLParsed.Parse(line)
				if err != nil {
					return ""
				}
				return segmentURL.String()
			}
			return line
		}
	}
	return ""
}

func checkSegment(segmentURL string, client *http.Client, timeout time.Duration) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "HEAD", segmentURL, nil)
	if err != nil {
		return false, "request error"
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return false, "unreachable"
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, fmt.Sprintf("status %d", resp.StatusCode)
	}

	return true, "ok"
}

package parse

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type StreamEntry struct {
	Channel   *string `json:"channel"`
	Feed      *string `json:"feed"`
	Title     string  `json:"title"`
	URL       *string `json:"url"`
	Quality   *string `json:"quality"`
	UserAgent *string `json:"user_agent"`
	Referrer  *string `json:"referrer"`
}

func ParseCommand(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "parse",
		Short: "Parse streams.json and import to PocketBase",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runParse(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func runParse(app *pocketbase.PocketBase) error {
	// Read JSON file
	jsonPath := filepath.Join("pkg", "json", "streams.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var streams []StreamEntry
	if err := json.Unmarshal(data, &streams); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("Found %d streams in JSON file\n", len(streams))

	// Get collections
	channelsCollection, err := app.FindCollectionByNameOrId("channels")
	if err != nil {
		return fmt.Errorf("failed to find channels collection: %w", err)
	}

	qualitiesCollection, err := app.FindCollectionByNameOrId("qualities")
	if err != nil {
		return fmt.Errorf("failed to find qualities collection: %w", err)
	}

	imported := 0
	skipped := 0

	for _, stream := range streams {
		// Skip if no URL
		if stream.URL == nil || *stream.URL == "" {
			skipped++
			continue
		}

		// Get or create quality record
		var qualityID string
		if stream.Quality != nil && *stream.Quality != "" {
			qualityID, err = getOrCreateQuality(app, qualitiesCollection, *stream.Quality)
			if err != nil {
				log.Printf("Warning: failed to get/create quality for %s: %v\n", stream.Title, err)
				continue
			}
		}

		// Create channel record
		channel := core.NewRecord(channelsCollection)
		
		if stream.Channel != nil {
			channel.Set("channel", *stream.Channel)
		}
		channel.Set("title", stream.Title)
		channel.Set("url", *stream.URL)
		if qualityID != "" {
			channel.Set("quality", qualityID)
		}

		if err := app.Save(channel); err != nil {
			log.Printf("Warning: failed to save channel %s: %v\n", stream.Title, err)
			continue
		}

		imported++
		if imported%100 == 0 {
			log.Printf("Imported %d channels...\n", imported)
		}
	}

	log.Printf("\nParse complete!")
	log.Printf("Imported: %d channels\n", imported)
	log.Printf("Skipped: %d entries (no URL)\n", skipped)

	return nil
}

func getOrCreateQuality(app *pocketbase.PocketBase, collection *core.Collection, qualityValue string) (string, error) {
	// Try to find existing quality
	records, err := app.FindRecordsByFilter(
		collection.Name,
		"quality = {:quality}",
		"-created",
		1,
		0,
		map[string]any{"quality": qualityValue},
	)
	if err != nil {
		return "", err
	}

	if len(records) > 0 {
		return records[0].Id, nil
	}

	// Create new quality record
	quality := core.NewRecord(collection)
	quality.Set("quality", qualityValue)

	if err := app.Save(quality); err != nil {
		return "", err
	}

	return quality.Id, nil
}

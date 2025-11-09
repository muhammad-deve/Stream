package logo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

type LogoEntry struct {
	Channel *string `json:"channel"`
	Feed    *string `json:"feed"`
	Tags    []string `json:"tags"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Format  string  `json:"format"`
	URL     string  `json:"url"`
}

func LogoCommand(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "logo",
		Short: "Parse logos.json and import to PocketBase",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runLogo(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func runLogo(app *pocketbase.PocketBase) error {
	// Read JSON file
	jsonPath := filepath.Join("pkg", "json", "logos.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	var logos []LogoEntry
	if err := json.Unmarshal(data, &logos); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	log.Printf("Found %d logos in JSON file\n", len(logos))

	// Get collections
	channelsCollection, err := app.FindCollectionByNameOrId("channels")
	if err != nil {
		return fmt.Errorf("failed to find channels collection: %w", err)
	}

	logosCollection, err := app.FindCollectionByNameOrId("logos")
	if err != nil {
		return fmt.Errorf("failed to find logos collection: %w", err)
	}

	// Build a map of all channels (lowercase channel name -> record)
	log.Printf("Loading all channels from database...\n")
	channelMap := make(map[string]*core.Record)
	
	allChannels, err := app.FindRecordsByFilter(
		channelsCollection.Name,
		"",
		"-created",
		10000, // Get first 10k channels
		0,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to load channels: %w", err)
	}
	
	for _, ch := range allChannels {
		channelName := ch.GetString("channel")
		channelMap[strings.ToLower(channelName)] = ch
	}
	
	log.Printf("Loaded %d channels from database\n", len(channelMap))

	imported := 0
	skipped := 0
	notFound := 0

	for _, logo := range logos {
		// Skip if channel is missing
		if logo.Channel == nil || *logo.Channel == "" {
			skipped++
			continue
		}

		// Skip if URL is missing
		if logo.URL == "" {
			skipped++
			continue
		}

		// Find channel by name (case-insensitive lookup)
		channelName := *logo.Channel
		channelRecord, found := channelMap[strings.ToLower(channelName)]
		
		if !found {
			notFound++
			skipped++
			continue
		}

		channelID := channelRecord.Id

		// Check if logo already exists for this channel
		existingLogos, err := app.FindRecordsByFilter(
			logosCollection.Name,
			"channel = {:channel}",
			"-created",
			1,
			0,
			map[string]any{"channel": channelID},
		)
		if err == nil && len(existingLogos) > 0 {
			skipped++
			continue
		}

		// Create logo record
		logoRecord := core.NewRecord(logosCollection)
		logoRecord.Set("channel", channelID)
		logoRecord.Set("logo_url", logo.URL)
		logoRecord.Set("width", logo.Width)
		logoRecord.Set("height", logo.Height)

		if err := app.Save(logoRecord); err != nil {
			log.Printf("Warning: failed to save logo for channel %s: %v\n", channelName, err)
			skipped++
			continue
		}

		imported++
		if imported%100 == 0 {
			log.Printf("Imported %d logos...\n", imported)
		}
	}

	log.Printf("\nLogo import complete!")
	log.Printf("Imported: %d logos\n", imported)
	log.Printf("Not found in database: %d channels\n", notFound)
	log.Printf("Skipped: %d entries (total)\n", skipped)

	return nil
}

// normalizeChannelName normalizes the channel name for comparison
func normalizeChannelName(channel string) string {
	// Convert to lowercase and trim whitespace
	return strings.TrimSpace(strings.ToLower(channel))
}

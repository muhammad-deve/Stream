package delete

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/spf13/cobra"
)

func DeleteCommand(app *pocketbase.PocketBase) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete all broken channels (is_working = false)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runDelete(app); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func runDelete(app *pocketbase.PocketBase) error {
	fmt.Println("🗑️  Starting deletion of broken channels...")

	// Get channels collection
	channelsCollection, err := app.FindCollectionByNameOrId("channels")
	if err != nil {
		return fmt.Errorf("failed to find channels collection: %w", err)
	}

	// Fetch all broken channels (is_working = false)
	records, err := app.FindRecordsByFilter(
		channelsCollection.Name,
		"is_working = false",
		"",
		0,
		0,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch broken channels: %w", err)
	}

	totalBroken := len(records)
	fmt.Printf("📊 Found %d broken channels to delete\n\n", totalBroken)

	if totalBroken == 0 {
		fmt.Println("✨ No broken channels to delete!")
		return nil
	}

	// Delete each broken channel
	deleted := 0
	failed := 0

	for _, record := range records {
		title := record.GetString("title")
		url := record.GetString("url")

		if err := app.Delete(record); err != nil {
			fmt.Printf("❌ Failed to delete: %s (%s) - %v\n", title, url, err)
			failed++
		} else {
			fmt.Printf("🗑️  Deleted: %s (%s)\n", title, url)
			deleted++
		}
	}

	fmt.Printf("\n✨ Deletion complete!\n")
	fmt.Printf("📊 Results: 🗑️  %d deleted | ❌ %d failed\n", deleted, failed)

	return nil
}

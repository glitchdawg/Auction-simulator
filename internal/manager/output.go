package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"auction-simulator/pkg/models"
)

// OutputGenerator handles the generation of output files
type OutputGenerator struct {
	outputDir string
}

// NewOutputGenerator creates a new output generator
func NewOutputGenerator(outputDir string) *OutputGenerator {
	return &OutputGenerator{
		outputDir: outputDir,
	}
}

// WriteAuctionResults writes individual auction result files
func (og *OutputGenerator) WriteAuctionResults(auctions []*models.Auction) error {
	// Ensure output directory exists
	if err := os.MkdirAll(og.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, auction := range auctions {
		filename := filepath.Join(og.outputDir, fmt.Sprintf("auction_%d_result.json", auction.ID))

		data, err := json.MarshalIndent(auction, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal auction %d: %w", auction.ID, err)
		}

		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write auction %d result: %w", auction.ID, err)
		}
	}

	return nil
}

// WriteSummary writes the execution summary file
func (og *OutputGenerator) WriteSummary(
	auctions []*models.Auction,
	firstStart, lastEnd time.Time,
	maxCPUs int,
	peakMemoryMB float64,
	avgGoroutines int,
) error {
	// Calculate statistics
	totalBids := 0
	auctionsWithNoBids := 0

	for _, auction := range auctions {
		totalBids += auction.TotalBids
		if auction.TotalBids == 0 {
			auctionsWithNoBids++
		}
	}

	avgBidsPerAuction := 0.0
	if len(auctions) > 0 {
		avgBidsPerAuction = float64(totalBids) / float64(len(auctions))
	}

	summary := models.ExecutionSummary{
		TotalAuctions:        len(auctions),
		FirstAuctionStart:    firstStart,
		LastAuctionEnd:       lastEnd,
		TotalExecutionTimeMs: lastEnd.Sub(firstStart).Milliseconds(),
		ResourceProfile: models.ResourceProfile{
			MaxCPUs:       maxCPUs,
			PeakMemoryMB:  peakMemoryMB,
			AvgGoroutines: avgGoroutines,
		},
		Statistics: models.Statistics{
			TotalBids:          totalBids,
			AvgBidsPerAuction:  avgBidsPerAuction,
			AuctionsWithNoBids: auctionsWithNoBids,
		},
	}

	filename := filepath.Join(og.outputDir, "execution_summary.json")

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}

	return nil
}

// PrintSummary prints a summary to the console
func (og *OutputGenerator) PrintSummary(
	auctions []*models.Auction,
	firstStart, lastEnd time.Time,
	maxCPUs int,
	peakMemoryMB float64,
	avgGoroutines int,
) {
	totalBids := 0
	auctionsWithNoBids := 0

	for _, auction := range auctions {
		totalBids += auction.TotalBids
		if auction.TotalBids == 0 {
			auctionsWithNoBids++
		}
	}

	avgBidsPerAuction := 0.0
	if len(auctions) > 0 {
		avgBidsPerAuction = float64(totalBids) / float64(len(auctions))
	}

	executionTime := lastEnd.Sub(firstStart)

	fmt.Println()
	for range 60 {
		fmt.Print("=")
	}
	fmt.Println()
	fmt.Println("AUCTION SIMULATOR - EXECUTION SUMMARY")
	for range 60 {
		fmt.Print("=")
	}
	fmt.Println()

	fmt.Printf("\nTotal Auctions:           %d\n", len(auctions))
	fmt.Printf("Total Execution Time:     %v (%.2f seconds)\n", executionTime, executionTime.Seconds())
	fmt.Printf("First Auction Start:      %s\n", firstStart.Format(time.RFC3339))
	fmt.Printf("Last Auction End:         %s\n", lastEnd.Format(time.RFC3339))

	fmt.Println("\nBid Statistics:")
	fmt.Printf("  Total Bids:             %d\n", totalBids)
	fmt.Printf("  Avg Bids per Auction:   %.2f\n", avgBidsPerAuction)
	fmt.Printf("  Auctions with No Bids:  %d\n", auctionsWithNoBids)

	fmt.Println("\nResource Usage:")
	fmt.Printf("  Max CPUs:               %d\n", maxCPUs)
	fmt.Printf("  Peak Memory:            %.2f MB\n", peakMemoryMB)
	fmt.Printf("  Avg Goroutines:         %d\n", avgGoroutines)

	for range 60 {
		fmt.Print("=")
	}
	fmt.Println()
}

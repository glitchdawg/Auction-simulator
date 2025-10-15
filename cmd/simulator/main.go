package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"auction-simulator/internal/manager"
	"auction-simulator/internal/resource"
	"auction-simulator/pkg/models"
)

func main() {
	// Parse command-line flags
	maxCPUs := flag.Int("cpus", runtime.NumCPU(), "Maximum number of CPUs to use")
	outputDir := flag.String("output", "output", "Output directory for results")
	seed := flag.Int64("seed", time.Now().UnixNano(), "Random seed for reproducibility")
	flag.Parse()

	// Set random seed for reproducibility
	rand.Seed(*seed)

	// Configure resource constraints
	runtime.GOMAXPROCS(*maxCPUs)

	config := models.ResourceConfig{
		MaxCPUs:     *maxCPUs,
		MaxMemoryMB: 0, // No hard limit, just monitoring
	}

	fmt.Println("===================================================")
	fmt.Println("        AUCTION SIMULATOR - STARTING")
	fmt.Println("===================================================")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Max CPUs:        %d\n", config.MaxCPUs)
	fmt.Printf("  Output Dir:      %s\n", *outputDir)
	fmt.Printf("  Random Seed:     %d\n", *seed)
	fmt.Printf("  Auctions:        %d\n", manager.NumAuctions)
	fmt.Printf("  Bidders:         %d\n", manager.NumBidders)
	fmt.Println("===================================================\n")

	// Create resource monitor
	monitor := resource.NewMonitor()
	monitor.Start(100 * time.Millisecond) // Sample every 100ms

	// Create auction manager
	mgr := manager.NewManager(config)

	// Run auctions
	ctx := context.Background()
	fmt.Println("Running auctions...")

	auctions, firstStart, lastEnd, err := mgr.Run(ctx)
	if err != nil {
		log.Fatalf("Error running auctions: %v", err)
	}

	// Stop monitoring
	monitor.Stop()

	// Get resource statistics
	maxCPUsUsed := monitor.GetMaxCPUs()
	peakMemoryMB := monitor.GetPeakMemoryMB()
	avgGoroutines := monitor.GetAvgGoroutines()

	fmt.Println("\nAll auctions completed!")
	fmt.Println("Generating output files...")

	// Generate output files
	outputGen := manager.NewOutputGenerator(*outputDir)

	if err := outputGen.WriteAuctionResults(auctions); err != nil {
		log.Fatalf("Error writing auction results: %v", err)
	}

	if err := outputGen.WriteSummary(
		auctions,
		firstStart,
		lastEnd,
		maxCPUsUsed,
		peakMemoryMB,
		avgGoroutines,
	); err != nil {
		log.Fatalf("Error writing summary: %v", err)
	}

	// Print summary to console
	outputGen.PrintSummary(
		auctions,
		firstStart,
		lastEnd,
		maxCPUsUsed,
		peakMemoryMB,
		avgGoroutines,
	)

	fmt.Printf("\nOutput files written to: %s\n", *outputDir)
	fmt.Println("  - 40 individual auction result files (auction_N_result.json)")
	fmt.Println("  - 1 execution summary file (execution_summary.json)")
	fmt.Println("\nSimulation completed successfully!")
}

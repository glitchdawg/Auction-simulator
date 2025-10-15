package manager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"auction-simulator/internal/auction"
	"auction-simulator/internal/bidder"
	"auction-simulator/pkg/models"
)

const (
	NumAuctions = 40
	NumBidders  = 100
)

// Manager orchestrates the execution of multiple concurrent auctions
type Manager struct {
	config  models.ResourceConfig
	bidders []*bidder.Bidder
}

// NewManager creates a new auction manager
func NewManager(config models.ResourceConfig) *Manager {
	// Create 100 bidders
	bidders := make([]*bidder.Bidder, NumBidders)
	for i := 0; i < NumBidders; i++ {
		bidders[i] = bidder.NewBidder(i + 1)
	}

	return &Manager{
		config:  config,
		bidders: bidders,
	}
}

// Run executes all auctions concurrently and returns the results
func (m *Manager) Run(ctx context.Context) ([]*models.Auction, time.Time, time.Time, error) {
	// Create channel for results
	results := make(chan *models.Auction, NumAuctions)

	var wg sync.WaitGroup

	// Create a function to notify all bidders about an auction
	notifyBidders := func(auction *models.Auction, bidChan chan<- models.Bid) {
		// Notify all 100 bidders about this auction
		for _, b := range m.bidders {
			b.ConsiderBid(auction, bidChan)
		}
	}

	// Launch all 40 auctions concurrently
	for i := 1; i <= NumAuctions; i++ {
		wg.Add(1)
		go func(auctionID int) {
			defer wg.Done()

			// Run auction with timeout (5 seconds)
			timeout := 5 * time.Second
			auction.Run(ctx, auctionID, timeout, notifyBidders, results)
		}(i)
	}

	// Wait for all auctions to complete in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	var auctionResults []*models.Auction
	for result := range results {
		auctionResults = append(auctionResults, result)
		fmt.Printf("Auction %d completed with %d bids\n", result.ID, result.TotalBids)
	}

	// Record actual first start time and last end time from results
	var firstStart, lastEnd time.Time
	if len(auctionResults) > 0 {
		firstStart = auctionResults[0].StartTime
		lastEnd = auctionResults[0].EndTime

		for _, a := range auctionResults {
			if a.StartTime.Before(firstStart) {
				firstStart = a.StartTime
			}
			if a.EndTime.After(lastEnd) {
				lastEnd = a.EndTime
			}
		}
	}

	return auctionResults, firstStart, lastEnd, nil
}

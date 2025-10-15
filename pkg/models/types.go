package models

import (
	"sync"
	"time"
)

// Bid represents a single bid in an auction
type Bid struct {
	BidderID  int       `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// Auction represents a single auction with its attributes and state
type Auction struct {
	ID         int         `json:"auction_id"`
	Attributes [20]float64 `json:"attributes"`
	Timeout    time.Duration `json:"-"`
	TimeoutMs  int64       `json:"timeout_ms"`
	StartTime  time.Time   `json:"start_time"`
	EndTime    time.Time   `json:"end_time"`
	Bids       []Bid       `json:"bids"`
	Winner     *Bid        `json:"winner"`
	TotalBids  int         `json:"total_bids"`
	mu         sync.Mutex
}

// NewAuction creates a new auction with random attributes
func NewAuction(id int, timeout time.Duration) *Auction {
	return &Auction{
		ID:        id,
		Timeout:   timeout,
		TimeoutMs: timeout.Milliseconds(),
		Bids:      make([]Bid, 0),
	}
}

// AddBid adds a bid to the auction in a thread-safe manner
func (a *Auction) AddBid(bid Bid) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Bids = append(a.Bids, bid)
}

// DetermineWinner finds the highest bid and sets it as the winner
func (a *Auction) DetermineWinner() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.TotalBids = len(a.Bids)

	if len(a.Bids) == 0 {
		a.Winner = nil
		return
	}

	// Find the highest bid (first one in case of tie)
	winner := &a.Bids[0]
	for i := 1; i < len(a.Bids); i++ {
		if a.Bids[i].Amount > winner.Amount {
			winner = &a.Bids[i]
		} else if a.Bids[i].Amount == winner.Amount && a.Bids[i].Timestamp.Before(winner.Timestamp) {
			// In case of tie, earlier timestamp wins
			winner = &a.Bids[i]
		}
	}
	a.Winner = winner
}

// AuctionResult represents the result of a single auction
type AuctionResult struct {
	AuctionID  int           `json:"auction_id"`
	Attributes [20]float64   `json:"attributes"`
	TotalBids  int           `json:"total_bids"`
	Winner     *Bid          `json:"winner"`
	Duration   time.Duration `json:"-"`
	DurationMs int64         `json:"duration_ms"`
}

// ExecutionSummary represents the overall execution summary
type ExecutionSummary struct {
	TotalAuctions        int              `json:"total_auctions"`
	FirstAuctionStart    time.Time        `json:"first_auction_start"`
	LastAuctionEnd       time.Time        `json:"last_auction_end"`
	TotalExecutionTimeMs int64            `json:"total_execution_time_ms"`
	ResourceProfile      ResourceProfile  `json:"resource_profile"`
	Statistics           Statistics       `json:"statistics"`
}

// ResourceProfile contains resource usage information
type ResourceProfile struct {
	MaxCPUs        int     `json:"max_cpus"`
	PeakMemoryMB   float64 `json:"peak_memory_mb"`
	AvgGoroutines  int     `json:"avg_goroutines"`
}

// Statistics contains aggregate statistics
type Statistics struct {
	TotalBids            int     `json:"total_bids"`
	AvgBidsPerAuction    float64 `json:"avg_bids_per_auction"`
	AuctionsWithNoBids   int     `json:"auctions_with_no_bids"`
}

// ResourceConfig defines resource constraints
type ResourceConfig struct {
	MaxCPUs     int
	MaxMemoryMB int64
}

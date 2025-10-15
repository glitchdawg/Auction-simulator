package bidder

import (
	"math/rand"
	"time"

	"auction-simulator/pkg/models"
)

// Bidder represents a bidder that participates in auctions
type Bidder struct {
	ID                int
	ParticipationRate float64 // Probability of participating (0.6-0.8)
}

// NewBidder creates a new bidder with given ID
func NewBidder(id int) *Bidder {
	return &Bidder{
		ID:                id,
		ParticipationRate: 0.6 + rand.Float64()*0.2, // 60-80% participation rate
	}
}

// ConsiderBid decides whether to bid and places a bid if decided to participate
func (b *Bidder) ConsiderBid(auction *models.Auction, bidChan chan<- models.Bid) {
	// Decide whether to participate
	if rand.Float64() > b.ParticipationRate {
		return // Not participating in this auction
	}

	go b.placeBid(auction, bidChan)
}

// placeBid calculates and places a bid for the given auction
func (b *Bidder) placeBid(auction *models.Auction, bidChan chan<- models.Bid) {
	// Simulate processing delay (10-500ms)
	processingDelay := time.Duration(10+rand.Intn(490)) * time.Millisecond
	time.Sleep(processingDelay)

	// Calculate bid amount based on weighted attribute scoring
	bidAmount := b.calculateBid(auction.Attributes)

	bid := models.Bid{
		BidderID:  b.ID,
		Amount:    bidAmount,
		Timestamp: time.Now(),
	}

	// Try to submit bid (may fail if auction has already closed)
	select {
	case bidChan <- bid:
		// Bid submitted successfully
	default:
		// Channel closed or full, auction likely ended
	}
}

// calculateBid calculates bid amount based on auction attributes
func (b *Bidder) calculateBid(attributes [20]float64) float64 {
	// Generate random weights for this bidder's preferences
	var score float64
	for i := 0; i < 20; i++ {
		weight := rand.Float64()
		score += attributes[i] * weight
	}

	// Normalize and scale to a reasonable bid range (e.g., 100-10000)
	bidAmount := 100 + (score/20)*9900

	// Add some randomness (±20%)
	randomFactor := 0.8 + rand.Float64()*0.4
	bidAmount *= randomFactor

	return bidAmount
}

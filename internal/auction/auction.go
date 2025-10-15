package auction

import (
	"context"
	"math/rand"
	"time"

	"auction-simulator/pkg/models"
)

// Run executes a single auction with the given timeout and bidder notifier
func Run(ctx context.Context, auctionID int, timeout time.Duration, notifyBidders func(*models.Auction, chan<- models.Bid), results chan<- *models.Auction) {
	auction := models.NewAuction(auctionID, timeout)

	// Generate random attributes for this auction (values between 0 and 1)
	for i := 0; i < 20; i++ {
		auction.Attributes[i] = rand.Float64()
	}

	auction.StartTime = time.Now()

	// Create a channel to receive bids (buffered to handle concurrent submissions)
	bidChan := make(chan models.Bid, 200)

	// Create context with timeout
	auctionCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Notify all bidders about this auction
	notifyBidders(auction, bidChan)

	// Collect bids until timeout
	done := make(chan struct{})
	go func() {
		for {
			select {
			case bid := <-bidChan:
				auction.AddBid(bid)
			case <-auctionCtx.Done():
				close(done)
				return
			}
		}
	}()

	// Wait for timeout
	<-auctionCtx.Done()
	<-done
	close(bidChan)

	auction.EndTime = time.Now()

	// Determine winner
	auction.DetermineWinner()

	// Send result
	results <- auction
}

// AuctionBroadcast contains auction information broadcasted to bidders
type AuctionBroadcast struct {
	Auction *models.Auction
	BidChan chan<- models.Bid
}

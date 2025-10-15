# Auction Simulator

A high-performance, concurrent auction simulation system written in Go that runs 40 auctions simultaneously with 100 bidders, implementing timeout mechanisms, resource monitoring, and comprehensive performance tracking.

## Overview

This project simulates a real-world auction system where:
- **40 auctions** run concurrently
- **100 bidders** participate across all auctions
- Each auction has **20 attributes** that influence bidding decisions
- **5-second timeout** per auction
- Bidders have **60-80% participation rate**
- Processing delays simulate real-world bid submission (10-500ms)
<!--
## Features

- ✅ Concurrent auction execution using goroutines
- ✅ Thread-safe bid collection with mutexes
- ✅ Context-based timeout management
- ✅ Real-time resource monitoring (CPU, Memory, Goroutines)
- ✅ Resource standardization for reproducible results
- ✅ JSON output for each auction and summary statistics
- ✅ Configurable CPU limits and random seed
- ✅ Comprehensive logging and console output
-->
<!--
## Project Structure

```
auction-simulator/
├── cmd/
│   └── simulator/
│       └── main.go              # Entry point
├── internal/
│   ├── auction/
│   │   └── auction.go           # Auction logic and timeout handling
│   ├── bidder/
│   │   └── bidder.go            # Bidder simulation and bid calculation
│   ├── manager/
│   │   ├── manager.go           # Orchestrates concurrent auctions
│   │   └── output.go            # Output generation (JSON, console)
│   └── resource/
│       └── monitor.go           # Resource usage monitoring
├── pkg/
│   └── models/
│       └── types.go             # Core data structures
├── output/                       # Generated result files (40 + 1)
├── go.mod                        # Go module definition
├── DESIGN.md                     # Detailed design document
└── README.md                     # This file
```
-->
 ## Architecture

### Concurrency Model

The system uses a **Fan-Out, Fan-In** pattern:

1. **Fan-Out**: Launch 40 auction goroutines simultaneously
2. **Parallel Processing**: Each auction notifies all 100 bidders
3. **Fan-In**: Collect results through a shared channel

### Key Components

#### 1. Auction Entity
- Manages 20 random attributes (0.0-1.0)
- Implements timeout using `context.WithTimeout`
- Thread-safe bid collection with `sync.Mutex`
- Winner determination (highest bid, earliest timestamp breaks ties)

#### 2. Bidder Entity
- Unique ID (1-100)
- Randomized participation rate (60-80%)
- Weighted attribute scoring for bid calculation
- Simulated processing delay (10-500ms)

#### 3. Auction Manager
- Orchestrates 40 concurrent auctions
- Tracks global start/end times
- Coordinates goroutines with `sync.WaitGroup`
- Collects and aggregates results

#### 4. Resource Monitor
- Samples CPU, memory, and goroutine count every 100ms
- Tracks peak memory usage
- Calculates average goroutine count
- Reports standardized resource profile

## Installation

### Prerequisites

- Go 1.21 or later
- Git (for cloning)

### Clone and Build

```bash
git clone <repository-url>
cd auction-simulator
go build -o auction-simulator.exe ./cmd/simulator
```

## Usage

### Basic Execution

```bash
./auction-simulator.exe
```

### Command-Line Options

```bash
./auction-simulator.exe [options]

Options:
  -cpus int
        Maximum number of CPUs to use (default: all available cores)
  -output string
        Output directory for results (default: "output")
  -seed int
        Random seed for reproducibility (default: current timestamp)
```
<!--
### Examples

```bash
# Run with 4 CPUs for standardized results
./auction-simulator.exe -cpus 4

# Use custom output directory
./auction-simulator.exe -output ./results

# Reproducible run with specific seed
./auction-simulator.exe -seed 12345 -cpus 4
```
-->

## Output Files

### Individual Auction Results

**Files**: `output/auction_{1-40}_result.json`

Each file contains:
- Auction ID and attributes
- All bids with timestamps
- Winner information
- Timeout and duration

**Example**: `auction_1_result.json`

```json
{
  "auction_id": 1,
  "attributes": [0.45, 0.89, ...],
  "timeout_ms": 5000,
  "start_time": "2025-10-15T23:40:53+05:30",
  "end_time": "2025-10-15T23:40:58+05:30",
  "total_bids": 66,
  "bids": [
    {
      "bidder_id": 42,
      "amount": 2340.23,
      "timestamp": "2025-10-15T23:40:53.311633+05:30"
    }
  ],
  "winner": {
    "bidder_id": 80,
    "amount": 3152.34,
    "timestamp": "2025-10-15T23:40:53.35534+05:30"
  }
}
```

### Execution Summary

**File**: `output/execution_summary.json`

Contains aggregate statistics:
- Total execution time
- Bid distribution statistics
- Resource usage profile

```json
{
  "total_auctions": 40,
  "first_auction_start": "2025-10-15T23:40:53+05:30",
  "last_auction_end": "2025-10-15T23:40:58+05:30",
  "total_execution_time_ms": 5007,
  "resource_profile": {
    "max_cpus": 4,
    "peak_memory_mb": 2.60,
    "avg_goroutines": 195
  },
  "statistics": {
    "total_bids": 2770,
    "avg_bids_per_auction": 69.25,
    "auctions_with_no_bids": 0
  }
}
```

## Performance Characteristics

### Expected Results

With typical hardware (4 CPUs, 8GB RAM):
- **Total Execution Time**: ~5 seconds (auctions run in parallel)
- **Total Bids**: 2,500-3,000 (60-80% participation × 100 bidders × 40 auctions)
- **Avg Bids per Auction**: 65-75
- **Peak Memory**: 2-5 MB
- **Goroutines**: ~200 concurrent

### Resource Standardization

To ensure reproducible results across different machines:

1. **CPU Limitation**: Use `-cpus` flag to limit GOMAXPROCS
2. **Fixed Seed**: Use `-seed` flag for deterministic randomness
3. **Consistent Environment**: Run on similar OS/architecture
<!--
**Example for standardized benchmarking**:

```bash
# Profile 1: 2 vCPUs
./auction-simulator.exe -cpus 2 -seed 12345

# Profile 2: 4 vCPUs
./auction-simulator.exe -cpus 4 -seed 12345

# Profile 3: 8 vCPUs
./auction-simulator.exe -cpus 8 -seed 12345
```
-->
## Testing

### Run Unit Tests

```bash
go test ./...
```

### Run with Race Detection

```bash
go test -race ./...
```

### Run Benchmarks

```bash
go test -bench=. ./...
```

## Implementation Details

### Concurrency Guarantees

- **Thread Safety**: All bid additions use `sync.Mutex`
- **No Race Conditions**: Verified with `go test -race`
- **Deadlock Prevention**: Proper channel closing and context cancellation
- **Timeout Compliance**: Context-based timeouts ensure auctions don't run forever

### Edge Cases Handled

1. **No Bids**: Auction completes with no winner
2. **Identical Bids**: First bid (by timestamp) wins
3. **Late Bids**: Rejected if submitted after timeout
4. **Channel Closures**: Graceful handling of closed channels

### Algorithm Complexity

- **Time Complexity**: O(n) where n = number of auctions (parallel execution)
- **Space Complexity**: O(a × b) where a = auctions, b = bids per auction
- **Goroutines**: O(a + b) = O(40 + 100) ≈ 140 + dynamic bid goroutines

## Assignment Requirements Checklist

- ✅ 40 auctions run concurrently
- ✅ 100 bidders participate across all auctions
- ✅ Each auction has 20 attributes
- ✅ Timeout mechanism implemented (5 seconds)
- ✅ Not all bidders respond (60-80% participation)
- ✅ Time measurement (first start → last end)
- ✅ Resource standardization (CPU limits, monitoring)
- ✅ 40 individual auction output files
- ✅ 1 execution summary file
- ✅ Well-documented, readable code
- ✅ Design document included
<!--
## Troubleshooting

### Issue: Low bid counts

**Solution**: Ensure the timeout is long enough for bidders to process and submit bids. The default 5 seconds accommodates delays up to 500ms per bidder.

### Issue: High memory usage

**Solution**: Reduce concurrent goroutines by adjusting buffer sizes or implementing goroutine pooling.

### Issue: Inconsistent results

**Solution**: Use `-seed` flag for deterministic randomness and `-cpus` flag to limit CPU variance.

## Future Enhancements

- [ ] Web UI for real-time auction visualization
- [ ] Distributed auction system across multiple nodes
- [ ] Database persistence for historical analysis
- [ ] Prometheus metrics export
- [ ] gRPC API for remote bidder integration
- [ ] Machine learning for bid prediction

## License

This project is created for educational purposes as part of an assignment.

## Author

Assignment implementation - 2025

## References

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Context Package](https://pkg.go.dev/context)
- [Effective Go](https://go.dev/doc/effective_go)
-->
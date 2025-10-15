package resource

import (
	"runtime"
	"sync"
	"time"
)

// Monitor tracks resource usage during execution
type Monitor struct {
	startTime     time.Time
	samples       []Sample
	mu            sync.Mutex
	stopChan      chan struct{}
	sampleTicker  *time.Ticker
}

// Sample represents a single resource measurement
type Sample struct {
	Timestamp    time.Time
	MemoryMB     float64
	NumGoroutines int
}

// NewMonitor creates a new resource monitor
func NewMonitor() *Monitor {
	return &Monitor{
		samples:  make([]Sample, 0),
		stopChan: make(chan struct{}),
	}
}

// Start begins monitoring resource usage
func (m *Monitor) Start(interval time.Duration) {
	m.startTime = time.Now()
	m.sampleTicker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-m.sampleTicker.C:
				m.takeSample()
			case <-m.stopChan:
				return
			}
		}
	}()
}

// Stop stops monitoring
func (m *Monitor) Stop() {
	if m.sampleTicker != nil {
		m.sampleTicker.Stop()
	}
	close(m.stopChan)

	// Take one final sample
	m.takeSample()
}

// takeSample captures current resource usage
func (m *Monitor) takeSample() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	sample := Sample{
		Timestamp:    time.Now(),
		MemoryMB:     float64(memStats.Alloc) / 1024 / 1024,
		NumGoroutines: runtime.NumGoroutine(),
	}

	m.mu.Lock()
	m.samples = append(m.samples, sample)
	m.mu.Unlock()
}

// GetPeakMemoryMB returns the peak memory usage in MB
func (m *Monitor) GetPeakMemoryMB() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.samples) == 0 {
		return 0
	}

	peak := m.samples[0].MemoryMB
	for _, s := range m.samples {
		if s.MemoryMB > peak {
			peak = s.MemoryMB
		}
	}

	return peak
}

// GetAvgGoroutines returns the average number of goroutines
func (m *Monitor) GetAvgGoroutines() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.samples) == 0 {
		return 0
	}

	total := 0
	for _, s := range m.samples {
		total += s.NumGoroutines
	}

	return total / len(m.samples)
}

// GetMaxCPUs returns the maximum number of CPUs being used
func (m *Monitor) GetMaxCPUs() int {
	return runtime.GOMAXPROCS(0)
}

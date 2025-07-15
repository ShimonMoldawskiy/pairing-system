package pairing

import (
	"sync"
)

// scoringContext is the set of metrics created in the first pass after filtering
// and passed to scorers for normalization support
type scoringContext struct {
	MaxStake        int64
	MinStake        int64
	MaxFeatureCount int
}

// map to store scoring contexts for passing between functions within the same goroutine
// without changing the function signatures
var scoringContextMap sync.Map // map[goroutineID]*ScoringContext

type scoringContextBuilder struct {
	mu              sync.Mutex
	initialized     bool
	minStake        int64
	maxStake        int64
	maxFeatureCount int
}

func newScoringContextBuilder() *scoringContextBuilder {
	return &scoringContextBuilder{}
}

func (b *scoringContextBuilder) updateFrom(p *Provider) {
	if p == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.initialized {
		b.minStake = p.Stake
		b.maxStake = p.Stake
		b.maxFeatureCount = len(p.Features)
		b.initialized = true
		return
	}

	if p.Stake < b.minStake {
		b.minStake = p.Stake
	}
	if p.Stake > b.maxStake {
		b.maxStake = p.Stake
	}
	if len(p.Features) > b.maxFeatureCount {
		b.maxFeatureCount = len(p.Features)
	}
}

func (b *scoringContextBuilder) build() *scoringContext {
	return &scoringContext{
		MinStake:        b.minStake,
		MaxStake:        b.maxStake,
		MaxFeatureCount: b.maxFeatureCount,
	}
}

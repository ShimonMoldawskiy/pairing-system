package pairing

import (
	"math"
)

var scoreWeights = map[string]float64{
	// keys must match Scorer.Name()
	"stake":    0.4,
	"feature":  0.4,
	"location": 0.2,
}

// StakeScorer assigns a score based on provider's stake relative to max stake
type StakeScorer struct{}

func (s StakeScorer) Name() string { return "stake" }

func (s StakeScorer) Score(p *Provider, policy *ConsumerPolicy, ctx *ScoringContext) float64 {
	if ctx.MaxStake == ctx.MinStake {
		return 1.0 // All have equal stake
	}
	return float64(p.Stake-ctx.MinStake) / float64(ctx.MaxStake-ctx.MinStake)
}

// FeatureScorer gives a higher score for more features than required
type FeatureScorer struct{}

func (s FeatureScorer) Name() string { return "feature" }

func (s FeatureScorer) Score(p *Provider, policy *ConsumerPolicy, ctx *ScoringContext) float64 {
	extra := len(p.Features) - len(policy.RequiredFeatures)
	if ctx.MaxFeatureCount == 0 {
		return 0
	}
	return math.Min(1.0, float64(extra+len(policy.RequiredFeatures))/float64(ctx.MaxFeatureCount))
}

// LocationScorer assigns 1.0 for exact match, else less (proximity logic later)
type LocationScorer struct{}

func (s LocationScorer) Name() string { return "location" }

func (s LocationScorer) Score(p *Provider, policy *ConsumerPolicy, ctx *ScoringContext) float64 {
	if policy.RequiredLocation == "" {
		return 1.0 // No restriction
	}
	if p.Location == policy.RequiredLocation {
		return 1.0
	}
	// Placeholder for future proximity calculation
	return 0.5
}

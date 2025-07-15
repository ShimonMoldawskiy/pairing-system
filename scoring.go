package pairing

import (
	"log"
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
	res := 0.0
	if ctx.MaxStake == ctx.MinStake {
		res = 1.0 // All have equal stake
	} else {
		res = float64(p.Stake-ctx.MinStake) / float64(ctx.MaxStake-ctx.MinStake)
	}
	if verbose {
		log.Printf("Provider %s %s score: %.3f\n", p.Address, s.Name(), res)
	}
	return res
}

// FeatureScorer gives a higher score for more features than required
type FeatureScorer struct{}

func (s FeatureScorer) Name() string { return "feature" }

func (s FeatureScorer) Score(p *Provider, policy *ConsumerPolicy, ctx *ScoringContext) float64 {
	extra := len(p.Features) - len(policy.RequiredFeatures)
	res := 0.0
	if ctx.MaxFeatureCount != 0 {
		res = math.Min(1.0, float64(extra+len(policy.RequiredFeatures))/float64(ctx.MaxFeatureCount))
	}
	if verbose {
		log.Printf("Provider %s %s score: %.3f\n", p.Address, s.Name(), res)
	}
	return res
}

// LocationScorer assigns 1.0 for exact match, else less (proximity logic later)
type LocationScorer struct{}

func (s LocationScorer) Name() string { return "location" }

func (s LocationScorer) Score(p *Provider, policy *ConsumerPolicy, ctx *ScoringContext) float64 {
	res := 0.0
	if policy.RequiredLocation == "" {
		res = 1.0 // No restriction
	} else if p.Location == policy.RequiredLocation {
		res = 1.0
	} else {
		// Placeholder for future proximity calculation
		res = 0.5
	}
	if verbose {
		log.Printf("Provider %s %s score: %.3f\n", p.Address, s.Name(), res)
	}
	return res
}

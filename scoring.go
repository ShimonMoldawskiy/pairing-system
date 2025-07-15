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

// linearStakeScorer assigns a score based on provider's stake relative to max stake
type linearStakeScorer struct{}

func (s linearStakeScorer) name() string { return "stake" }

func (s linearStakeScorer) score(p *Provider, policy *ConsumerPolicy, ctx *scoringContext) float64 {
	res := 0.0
	if ctx.MaxStake == ctx.MinStake {
		res = 1.0 // All have equal stake
	} else {
		res = float64(p.Stake-ctx.MinStake) / float64(ctx.MaxStake-ctx.MinStake)
	}
	if verbose {
		log.Printf("Provider %s %s score: %.3f\n", p.Address, s.name(), res)
	}
	return res
}

// linearFeatureCountScorer gives a higher score for more features than required
type linearFeatureCountScorer struct{}

func (s linearFeatureCountScorer) name() string { return "feature" }

func (s linearFeatureCountScorer) score(p *Provider, policy *ConsumerPolicy, ctx *scoringContext) float64 {
	extra := len(p.Features) - len(policy.RequiredFeatures)
	res := 0.0
	if ctx.MaxFeatureCount != 0 {
		res = math.Min(1.0, float64(extra+len(policy.RequiredFeatures))/float64(ctx.MaxFeatureCount))
	}
	if verbose {
		log.Printf("Provider %s %s score: %.3f\n", p.Address, s.name(), res)
	}
	return res
}

// locationProximityTableScorer assigns 1.0 for exact match, else less
type locationProximityTableScorer struct{}

func (s locationProximityTableScorer) name() string { return "location" }

func (s locationProximityTableScorer) score(p *Provider, policy *ConsumerPolicy, ctx *scoringContext) float64 {
	res := locationProximity(p.Location, policy.RequiredLocation)
	if verbose {
		log.Printf("Provider %s %s score: %.3f\n", p.Address, s.name(), res)
	}
	return res
}

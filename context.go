package pairing

// scoringContext is the set of metrics created in the first pass after filtering
// and passed to scorers for normalization support
type scoringContext struct {
	MaxStake        int64
	MinStake        int64
	MaxFeatureCount int
}

// buildScoringContext analyzes filtered providers to extract global metrics
func buildScoringContext(providers []*Provider) *scoringContext {

	maxStake := providers[0].Stake
	minStake := providers[0].Stake
	maxFeatures := len(providers[0].Features)

	for _, p := range providers[1:] {
		if p.Stake > maxStake {
			maxStake = p.Stake
		}
		if p.Stake < minStake {
			minStake = p.Stake
		}
		if len(p.Features) > maxFeatures {
			maxFeatures = len(p.Features)
		}
	}

	return &scoringContext{
		MaxStake:        maxStake,
		MinStake:        minStake,
		MaxFeatureCount: maxFeatures,
	}
}

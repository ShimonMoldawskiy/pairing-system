package pairing

// BuildScoringContext analyzes filtered providers to extract global max/min data
func BuildScoringContext(providers []*Provider) *ScoringContext {

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

	return &ScoringContext{
		MaxStake:        maxStake,
		MinStake:        minStake,
		MaxFeatureCount: maxFeatures,
	}
}

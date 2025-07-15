package pairing

// ScoringContext is created in the first pass after filtering
// and passed to scorers for normalization support
type ScoringContext struct {
	MaxStake        int64
	MinStake        int64
	MaxFeatureCount int
}

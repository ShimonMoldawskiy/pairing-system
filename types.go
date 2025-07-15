package pairing

type Provider struct {
	Address  string
	Stake    int64
	Location string
	Features []string
}

type ConsumerPolicy struct {
	RequiredLocation string
	RequiredFeatures []string
	MinStake         int64
}

type PairingScore struct {
	Provider   *Provider
	Score      float64
	Components map[string]float64
}

type PairingSystem interface {
	FilterProviders(providers []*Provider, policy *ConsumerPolicy) []*Provider
	RankProviders(providers []*Provider, policy *ConsumerPolicy) []*PairingScore
	GetPairingList(providers []*Provider, policy *ConsumerPolicy) ([]*Provider, error)
}

// Filter interface for pipeline filtering
type Filter interface {
	Name() string
	Apply(*Provider, *ConsumerPolicy) bool
}

// Scorer interface for score computation
type Scorer interface {
	Name() string
	Score(*Provider, *ConsumerPolicy, *ScoringContext) float64
}

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

// filter interface for pipeline filtering
type filter interface {
	name() string
	apply(*Provider, *ConsumerPolicy) bool
}

// scorer interface for score computation
type scorer interface {
	name() string
	score(*Provider, *ConsumerPolicy, *scoringContext) float64
}

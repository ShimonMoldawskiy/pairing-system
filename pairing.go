package pairing

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

type MainPairingSystem struct{}

func (ps *MainPairingSystem) GetPairingList(providers []*Provider, policy *ConsumerPolicy) (result []*Provider, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in GetPairingList: %v\n", r)
			err = fmt.Errorf("internal error occurred: %v", r)
			result = nil
		}
	}()

	if len(providers) == 0 {
		return nil, errors.New("provider list is empty")
	}
	if policy == nil {
		return nil, errors.New("policy is nil")
	}
	if policy.MinStake < 0 {
		return nil, errors.New("consumer MinStake is negative")
	}

	// Normalize Consumer required features
	normalizedPolicy := &ConsumerPolicy{
		// Copying the policy to avoid modifying the original
		RequiredLocation: policy.RequiredLocation,
		RequiredFeatures: normalizeFeatures(policy.RequiredFeatures),
		MinStake:         policy.MinStake,
	}

	filteredProviders := ps.FilterProviders(providers, normalizedPolicy)

	if len(filteredProviders) == 0 {
		return nil, errors.New("no matching providers after filtering")
	}

	scores := ps.RankProviders(filteredProviders, normalizedPolicy)

	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].Score == scores[j].Score {
			return scores[i].Provider.Address < scores[j].Provider.Address
		}
		return scores[i].Score > scores[j].Score
	})

	top := []*Provider{}
	for i := 0; i < len(scores) && i < 5; i++ {
		top = append(top, scores[i].Provider)
	}

	return top, nil
}

func (ps *MainPairingSystem) FilterProviders(providers []*Provider, policy *ConsumerPolicy) []*Provider {
	// Build filter pipeline
	filters := []filter{
		rejectEmptyAddressFilter{},
		normalizedFeaturesFilter{},
		locationProximityTableFilter{},
		stakeMinFilter{},
	}

	return concurrentFilterPipeline(providers, policy, filters)
}

func (ps *MainPairingSystem) RankProviders(providers []*Provider, policy *ConsumerPolicy) []*PairingScore {
	ctx := buildScoringContext(providers)

	// Build scorer pipeline
	scorers := []scorer{
		linearStakeScorer{},
		linearFeatureCountScorer{},
		locationProximityTableScorer{},
	}

	return concurrentScoringPipeline(providers, policy, ctx, scorers)
}

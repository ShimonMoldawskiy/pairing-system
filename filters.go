package pairing

import "strings"

const proximityThreshold = 0.0

// RejectEmptyAddressFilter filters out providers with empty addresses
type RejectEmptyAddressFilter struct{}

func (f RejectEmptyAddressFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
	return strings.TrimSpace(p.Address) != ""
}

// RequiredFeaturesFilter ensures the provider supports all required features
type RequiredFeaturesFilter struct{}

func (f RequiredFeaturesFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {

	if len(policy.RequiredFeatures) == 0 {
		return true
	}
	if len(p.Features) == 0 {
		return false
	}

	// Use the fact that the features are already normalized - sorted and deduplicated
	i, j := 0, 0
	for i < len(policy.RequiredFeatures) && j < len(p.Features) {
		if policy.RequiredFeatures[i] == p.Features[j] {
			i++
			j++
		} else if policy.RequiredFeatures[i] > p.Features[j] {
			j++
		} else {
			// Required feature is missing
			return false
		}
	}

	return i == len(policy.RequiredFeatures)
}

// LocationProximityFilter applies filtering based on unimplemented proximity logic
type LocationProximityFilter struct {
	ProximityThreshold float64
}

func (f LocationProximityFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
	if policy.RequiredLocation == "" {
		return true
	}
	if p.Location == "" {
		return false
	}
	// Proximity logic will be implemented later
	return true
}

// StakeMinFilter filters out providers that do not meet minimum stake
type StakeMinFilter struct{}

func (f StakeMinFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
	return p.Stake >= policy.MinStake
}

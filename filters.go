package pairing

import (
	"log"
	"strings"
)

const proximityThreshold = 0.0

// RejectEmptyAddressFilter filters out providers with empty addresses
type RejectEmptyAddressFilter struct{}

func (f RejectEmptyAddressFilter) Name() string { return "empty-address" }

func (f RejectEmptyAddressFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
	res := strings.TrimSpace(p.Address) != ""
	if !res && verbose {
		log.Printf("Provider %s filtered out by %s filter: %v\n", p.Address, f.Name(), p)
	}
	return res
}

// RequiredFeaturesFilter ensures the provider supports all required features
type RequiredFeaturesFilter struct{}

func (f RequiredFeaturesFilter) Name() string { return "feature" }

func (f RequiredFeaturesFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {

	if len(policy.RequiredFeatures) == 0 {
		return true
	}
	if len(p.Features) == 0 {
		if verbose {
			log.Printf("Provider %s filtered out by %s filter\n", p.Address, f.Name())
		}
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
			if verbose {
				log.Printf("Provider %s filtered out by %s filter\n", p.Address, f.Name())
			}
			return false
		}
	}

	res := i == len(policy.RequiredFeatures)
	if !res && verbose {
		log.Printf("Provider %s filtered out by %s filter", p.Address, f.Name())
	}
	return res
}

// LocationProximityFilter applies filtering based on unimplemented proximity logic
type LocationProximityFilter struct {
	ProximityThreshold float64
}

func (f LocationProximityFilter) Name() string { return "location" }

func (f LocationProximityFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
	if policy.RequiredLocation == "" {
		return true
	}
	if p.Location == "" {
		if verbose {
			log.Printf("Provider %s filtered out by %s filter", p.Address, f.Name())
		}
		return false
	}
	// Proximity logic will be implemented later
	return true
}

// StakeMinFilter filters out providers that do not meet minimum stake
type StakeMinFilter struct{}

func (f StakeMinFilter) Name() string { return "stake" }

func (f StakeMinFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
	res := p.Stake >= policy.MinStake
	if !res && verbose {
		log.Printf("Provider %s filtered out by %s filter\n", p.Address, f.Name())
	}
	return res
}

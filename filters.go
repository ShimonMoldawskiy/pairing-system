package pairing

import (
	"log"
	"strings"
)

const proximityThreshold = 0.0

// rejectEmptyAddressFilter filters out providers with empty addresses
type rejectEmptyAddressFilter struct{}

func (f rejectEmptyAddressFilter) name() string { return "empty-address" }

func (f rejectEmptyAddressFilter) apply(p *Provider, policy *ConsumerPolicy) bool {
	res := strings.TrimSpace(p.Address) != ""
	if !res && verbose {
		log.Printf("Provider %s filtered out by %s filter: %v\n", p.Address, f.name(), p)
	}
	return res
}

// normalizedFeaturesFilter ensures the provider supports all required features
type normalizedFeaturesFilter struct{}

func (f normalizedFeaturesFilter) name() string { return "feature" }

func (f normalizedFeaturesFilter) apply(p *Provider, policy *ConsumerPolicy) bool {

	if len(policy.RequiredFeatures) == 0 {
		return true
	}
	if len(p.Features) == 0 {
		if verbose {
			log.Printf("Provider %s filtered out by %s filter\n", p.Address, f.name())
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
				log.Printf("Provider %s filtered out by %s filter\n", p.Address, f.name())
			}
			return false
		}
	}

	res := i == len(policy.RequiredFeatures)
	if !res && verbose {
		log.Printf("Provider %s filtered out by %s filter", p.Address, f.name())
	}
	return res
}

// locationProximityFilter applies filtering based on unimplemented proximity logic
type locationProximityFilter struct {
	ProximityThreshold float64
}

func (f locationProximityFilter) name() string { return "location" }

func (f locationProximityFilter) apply(p *Provider, policy *ConsumerPolicy) bool {
	if policy.RequiredLocation == "" {
		return true
	}
	if p.Location == "" {
		if verbose {
			log.Printf("Provider %s filtered out by %s filter", p.Address, f.name())
		}
		return false
	}
	// Proximity logic will be implemented later
	return true
}

// stakeMinFilter filters out providers that do not meet minimum stake
type stakeMinFilter struct{}

func (f stakeMinFilter) name() string { return "stake" }

func (f stakeMinFilter) apply(p *Provider, policy *ConsumerPolicy) bool {
	res := p.Stake >= policy.MinStake
	if !res && verbose {
		log.Printf("Provider %s filtered out by %s filter\n", p.Address, f.name())
	}
	return res
}

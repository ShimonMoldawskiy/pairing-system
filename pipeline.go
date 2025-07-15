package pairing

import (
	"log"
	"runtime"
	"sync"
)

var maxConcurrency int

func init() {
	n := runtime.NumCPU() * 2
	if n <= 0 {
		maxConcurrency = 10
	} else {
		maxConcurrency = n
	}
}

func concurrentFilterPipeline(providers []*Provider, policy *ConsumerPolicy, filters []filter) []*Provider {
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
	result := []*Provider{}

	addressSeen := make(map[string]bool)

	for _, provider := range providers {
		p := provider
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			if p == nil {
				log.Println("Skipping nil provider")
				return
			}

			if _, exists := addressSeen[p.Address]; exists {
				log.Printf("Duplicate provider address detected: %s\n", p.Address)
			} else {
				mu.Lock()
				addressSeen[p.Address] = true
				mu.Unlock()
			}

			// Normalize Provider features
			normalizedP := &Provider{
				// Copying the provider to avoid modifying the original
				Address:  p.Address,
				Stake:    p.Stake,
				Location: p.Location,
				Features: normalizeFeatures(p.Features),
			}

			for _, f := range filters {
				if !f.apply(normalizedP, policy) {
					log.Printf("Provider %s filtered out by %T\n", p.Address, f)
					return
				}
			}

			mu.Lock()
			result = append(result, normalizedP)
			mu.Unlock()
		}()
	}

	wg.Wait()
	return result
}

func concurrentScoringPipeline(providers []*Provider, policy *ConsumerPolicy, ctx *scoringContext, scorers []scorer) []*PairingScore {
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
	scoreChan := make(chan *PairingScore, len(providers))

	for _, p := range providers {
		p := p
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in scoring provider %s: %v\n", p.Address, r)
				}
				<-sem
			}()

			components := map[string]float64{}
			sum := 0.0
			for _, scorer := range scorers {
				v := scorer.score(p, policy, ctx)
				components[scorer.name()] = v
				sum += scoreWeights[scorer.name()] * v
			}

			if verbose {
				log.Printf("Provider %s final score: %.3f\n", p.Address, sum)
			}

			scoreChan <- &PairingScore{
				Provider:   p,
				Score:      sum,
				Components: components,
			}
		}()
	}

	wg.Wait()
	close(scoreChan)

	scores := []*PairingScore{}
	for s := range scoreChan {
		scores = append(scores, s)
	}
	return scores
}

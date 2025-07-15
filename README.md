# Lava Network Provider Pairing System

## Overview
This Go module implements the provider pairing logic for the PRC Gateway Network. It takes a list of providers and a consumer policy, filters out incompatible providers, scores the remaining ones, and returns the top 5 best-matching providers based on multiple criteria.

## Design Decisions

### Two-Pass Architecture
- **First Pass (Filtering)**: Applies filters concurrently to exclude invalid providers and collects normalization metrics (max/min stake, max features count).
- **Second Pass (Scoring)**: Computes normalized scores concurrently using previously gathered metrics.

### Concurrency
- Limits the number of concurrent goroutines via `maxConcurrency` parameter.
- Uses `sync.WaitGroup` and semaphore channels for controlled parallelism.

### Scoring Strategy
- StakeScore: Normalized based on min-max stake.
- FeatureScore: Normalized based on number of features.
- LocationScore: 1.0 for exact match.
- Total score is a weighted sum defined by `scoreWeights` parameter.

### Composable Filters and Scorers
A pipeline pattern is used for Filters and Scorers which allows for:
- Clean separation of logic
- Easy testing and extension
- Plug-and-play filtering/scoring strategies

### Feature Normalization
- Consumer required features are trimmed, deduplicated, and sorted before filtering.
- Provider features are also normalized before comparison.
- Sorted feartures allow applying feature comparison algorithm with O(n + m) complexity.

### Preserving Input Immutability
To ensure that input data remains unchanged, internal copies of the providers and policy are created for normalization and scoring.

### Observability
A global `verbose` flag is added to log the reason why a provider was filtered out or received a certain score. 
The edge cases are also logged.

## Assumptions
- Feature and Location field values are case-sensitive.
- Providers with empty addresses are rejected.
- Duplicate provider addresses are allowed but logged.
- Score weights are configurable per scorer.

## Edge Case Handling
- Empty or nil providers list → error.
- Nil policy or negative stake → error.
- Provider stake exactly equals min → pass.
- Stake = 0 → StakeScore = 0 unless all are 0.
- Equal stake → StakeScore = 1.0.
- Missing location or features → filtered is not applied.
- Tie-breaking by provider address.
- Less than 5 valid providers → return as many as possible.

## Scoring/Filtering Flexibility
For the future it's possible to:
- Add a new struct that implements `Filter` or `Scorer` interfaces.
- Register it in the corresponding pipeline.
- Adjust the scoring weights.

```go
// Example:
type MyCustomFilter struct{}
func (f MyCustomFilter) Apply(p *Provider, policy *ConsumerPolicy) bool {
    return p.SomeField != "bad"
}

filters := []Filter{
    MyCustomFilter{},
    // ... existing filters
}
```

## License
This code is intended for evaluation and demonstration purposes.

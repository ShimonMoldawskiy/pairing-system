package pairing

import (
	"testing"
)

func TestNormalizeFeatures(t *testing.T) {
	input := []string{"  rpc", "REST", "rpc", "rest ", "", "REST"}
	expected := []string{"REST", "rest", "rpc"}
	output := normalizeFeatures(input)
	if len(output) != len(expected) {
		t.Fatalf("Expected %v, got %v", expected, output)
	}
	for i := range expected {
		if output[i] != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], output[i])
		}
	}
}

func TestRejectEmptyAddressFilter(t *testing.T) {
	f := rejectEmptyAddressFilter{}
	p1 := &Provider{Address: ""}
	p2 := &Provider{Address: "abc"}
	if f.apply(p1, &ConsumerPolicy{}) {
		t.Error("Expected false for empty address")
	}
	if !f.apply(p2, &ConsumerPolicy{}) {
		t.Error("Expected true for valid address")
	}
}

func TestRequiredFeaturesFilter(t *testing.T) {
	f := normalizedFeaturesFilter{}
	policy := &ConsumerPolicy{RequiredFeatures: []string{"rest", "rpc"}}
	provider := &Provider{Features: []string{"grpc", "rest", "rpc"}}
	if !f.apply(provider, policy) {
		t.Error("Expected provider to pass required features filter")
	}

	provider.Features = []string{"rest"}
	if f.apply(provider, policy) {
		t.Error("Expected provider to fail required features filter")
	}
}

func TestStakeMinFilter(t *testing.T) {
	f := stakeMinFilter{}
	policy := &ConsumerPolicy{MinStake: 100}
	p1 := &Provider{Stake: 100}
	p2 := &Provider{Stake: 50}
	if !f.apply(p1, policy) {
		t.Error("Expected provider with exact stake to pass")
	}
	if f.apply(p2, policy) {
		t.Error("Expected provider with low stake to fail")
	}
}

func TestLocationProximityFilter(t *testing.T) {
	f := locationProximityFilter{ProximityThreshold: 0.0}
	policy := &ConsumerPolicy{RequiredLocation: "2"}
	p1 := &Provider{Location: "2"}
	p2 := &Provider{Location: ""}
	if !f.apply(p1, policy) {
		t.Error("Expected matching location to pass")
	}
	if f.apply(p2, policy) {
		t.Error("Expected empty location to fail")
	}
}

func TestScorers(t *testing.T) {
	ctx := &scoringContext{MaxStake: 200, MinStake: 100, MaxFeatureCount: 6}
	policy := &ConsumerPolicy{RequiredFeatures: []string{"rpc", "rest"}, RequiredLocation: "2"}
	provider := &Provider{Stake: 150, Features: []string{"rpc", "rest", "grpc"}, Location: "2"}

	stakeScore := linearStakeScorer{}.score(provider, policy, ctx)
	if stakeScore <= 0 || stakeScore >= 1 {
		t.Errorf("Unexpected stake score: %f", stakeScore)
	}

	featureScore := linearFeatureCountScorer{}.score(provider, policy, ctx)
	if featureScore <= 0 {
		t.Errorf("Unexpected feature score: %f", featureScore)
	}

	locationScore := locationProximityScorer{}.score(provider, policy, ctx)
	if locationScore != 1.0 {
		t.Errorf("Expected exact location score to be 1.0, got %f", locationScore)
	}
}

func TestMainPairingSystem_GetPairingList(t *testing.T) {
	ps := &MainPairingSystem{}
	providers := []*Provider{
		{Address: "1", Stake: 150, Location: "2", Features: []string{"rpc", "rest"}},
		{Address: "2", Stake: 50, Location: "3", Features: []string{"grpc"}},
		{Address: "3", Stake: 200, Location: "2", Features: []string{"rpc", "rest", "grpc"}},
	}
	policy := &ConsumerPolicy{
		RequiredLocation: "2",
		RequiredFeatures: []string{"rpc", "rest"},
		MinStake:         100,
	}
	topProviders, err := ps.GetPairingList(providers, policy)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(topProviders) != 2 {
		t.Fatalf("Expected 2 top providers, got %d", len(topProviders))
	}
	if topProviders[0].Address != "3" || topProviders[1].Address != "1" {
		t.Errorf("Expected top providers to be 3 and 1, got %s and %s", topProviders[0].Address, topProviders[1].Address)
	}
	if topProviders[0].Stake < topProviders[1].Stake {
		t.Errorf("Expected provider 3 to have higher stake than provider 1, got %d and %d", topProviders[0].Stake, topProviders[1].Stake)
	}
}

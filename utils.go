package pairing

import (
	"sort"
	"strings"
)

var verbose = false

func EnableVerboseLogging() {
	verbose = true
}

// normalizeFeatures trims whitespace and sorts + deduplicates features
func normalizeFeatures(features []string) []string {
	featureSet := make(map[string]struct{})
	for _, f := range features {
		trimmed := strings.TrimSpace(f)
		if trimmed != "" {
			featureSet[trimmed] = struct{}{}
		}
	}
	result := make([]string, 0, len(featureSet))
	for f := range featureSet {
		result = append(result, f)
	}
	sort.Strings(result)
	return result
}

var locationProximityMap = map[string]map[string]float64{
	"EU": {
		"EU":   1.0,
		"US":   0.6,
		"ASIA": 0.3,
	},
	"US": {
		"EU":   0.6,
		"US":   1.0,
		"ASIA": 0.4,
	},
	"ASIA": {
		"EU":   0.3,
		"US":   0.4,
		"ASIA": 1.0,
	},
}

func locationProximity(loc1, loc2 string) float64 {
	if loc1 == loc2 {
		return 1.0 // Exact match
	}
	if m, ok := locationProximityMap[loc1]; ok {
		if score, ok := m[loc2]; ok {
			return score
		}
	}
	return 0.0 // unknown or unmatched
}

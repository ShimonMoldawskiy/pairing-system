package pairing

import (
	"sort"
	"strings"
)

// NormalizeFeatures trims whitespace and sorts + deduplicates features
func NormalizeFeatures(features []string) []string {
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

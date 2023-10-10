package handlers

import "github.com/matthisholleville/mapsyncproxy/pkg/haproxy"

func jsonResponse(message string) *map[string]string {
	response := map[string]string{"status": message}
	return &response
}

func hasDuplicateKeys(objects []haproxy.MapEntrie) bool {
	seen := make(map[string]bool)

	for _, obj := range objects {
		if seen[obj.Key] {
			return true
		}
		seen[obj.Key] = true
	}

	return false
}

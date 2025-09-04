package myipld

import (
	"sort"
)

func sortMapKeys(v interface{}) interface{} {
	switch x := v.(type) {
	case map[string]interface{}:
		sorted := make(map[string]interface{})
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sorted[k] = sortMapKeys(x[k])
		}
		return sorted
	case []interface{}:
		for i, v := range x {
			x[i] = sortMapKeys(v)
		}
		return x
	default:
		return v
	}
}
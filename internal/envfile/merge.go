package envfile

// MergeStrategy controls how conflicting keys are handled during merge.
type MergeStrategy int

const (
	StrategyKeepBase MergeStrategy = iota
	StrategyOverwrite
	StrategyPrompt
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Merged  map[string]string
	Skipped []string
	Added   []string
	Updated []string
}

// HasChanges reports whether the merge resulted in any additions or updates.
func (r MergeResult) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Updated) > 0
}

// Merge combines base and overlay env maps according to the given strategy.
// Base entries are always preserved; overlay entries are added or merged per strategy.
func Merge(base, overlay map[string]string, strategy MergeStrategy) MergeResult {
	result := MergeResult{
		Merged: make(map[string]string),
	}

	// Copy base into merged.
	for k, v := range base {
		result.Merged[k] = v
	}

	for k, v := range overlay {
		if existing, exists := base[k]; exists {
			if existing == v {
				continue
			}
			switch strategy {
			case StrategyOverwrite:
				result.Merged[k] = v
				result.Updated = append(result.Updated, k)
			case StrategyKeepBase:
				result.Skipped = append(result.Skipped, k)
			case StrategyPrompt:
				// Caller must handle prompt; default to keep base.
				result.Skipped = append(result.Skipped, k)
			}
		} else {
			result.Merged[k] = v
			result.Added = append(result.Added, k)
		}
	}

	return result
}

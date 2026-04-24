package envfile

import "fmt"

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Added     int
	Skipped   int
	Overwrite int
}

// PromoteOptions controls promotion behaviour.
type PromoteOptions struct {
	// Overwrite allows values in target to be replaced by source values.
	Overwrite bool
	// OnlyKeys, when non-empty, restricts promotion to these keys only.
	OnlyKeys []string
}

// Promote copies entries from src into dst according to opts.
// It returns the merged map and a result summary.
func Promote(src, dst map[string]string, opts PromoteOptions) (map[string]string, PromoteResult, error) {
	if src == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: source map must not be nil")
	}
	if dst == nil {
		dst = make(map[string]string)
	}

	filter := make(map[string]bool, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		filter[k] = true
	}

	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}

	var res PromoteResult
	for key, srcVal := range src {
		if len(filter) > 0 && !filter[key] {
			continue
		}
		existing, exists := result[key]
		switch {
		case !exists:
			result[key] = srcVal
			res.Added++
		case exists && srcVal == existing:
			res.Skipped++
		case opts.Overwrite:
			result[key] = srcVal
			res.Overwrite++
		default:
			res.Skipped++
		}
	}
	return result, res, nil
}

package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references
var interpolatePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateError describes a variable reference that could not be resolved.
type InterpolateError struct {
	Key string
	Ref string
}

func (e *InterpolateError) Error() string {
	return fmt.Sprintf("key %q references undefined variable %q", e.Key, e.Ref)
}

// Interpolate resolves variable references within env values using the
// provided env map. References to undefined variables are returned as errors.
// The env map is updated in-place with resolved values.
func Interpolate(env map[string]string) []error {
	var errs []error

	for key, value := range env {
		resolved, err := resolveValue(key, value, env)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		env[key] = resolved
	}

	return errs
}

func resolveValue(key, value string, env map[string]string) (string, error) {
	var resolveErr error

	result := interpolatePattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}

		varName := extractVarName(match)
		if resolved, ok := env[varName]; ok {
			return resolved
		}

		resolveErr = &InterpolateError{Key: key, Ref: varName}
		return match
	})

	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractVarName(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}

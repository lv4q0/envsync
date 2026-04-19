package diff

import (
	"fmt"
	"io"
	"sort"
)

// Summary holds counts of each diff status.
type Summary struct {
	Added   int
	Removed int
	Changed int
	Same    int
}

// Report writes a formatted diff report to w.
func Report(w io.Writer, entries []Entry) Summary {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	var s Summary
	for _, e := range sorted {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusChanged:
			s.Changed++
		case StatusSame:
			s.Same++
		}
		fmt.Fprintln(w, e.String())
	}

	fmt.Fprintf(w, "\nSummary: +%d added, -%d removed, ~%d changed, %d unchanged\n",
		s.Added, s.Removed, s.Changed, s.Same)
	return s
}

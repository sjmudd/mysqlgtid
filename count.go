package mysqlgtid

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// uuidRE matches a MySQL server UUID: 8-4-4-4-12 hex digits.
	uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	// tagRE matches a GTID tag: starts with a letter or underscore, followed
	// by letters, digits or underscores. An empty tag is handled separately.
	tagRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	// intervalRE matches a single interval: "n" or "n-n".
	intervalRE = regexp.MustCompile(`^[0-9]+(-[0-9]+)?$`)
)

// intervalTransactionCount converts a single interval ("n" or "n-n") into
// the number of transactions it represents.
func intervalTransactionCount(interval string) (int64, error) {
	if !intervalRE.MatchString(interval) {
		return 0, fmt.Errorf("invalid interval %q", interval)
	}

	parts := strings.SplitN(interval, "-", 2)

	minValue, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid transaction id %q: %w", parts[0], err)
	}
	if minValue < 1 {
		return 0, fmt.Errorf("invalid transaction id %d: sequence numbers start at 1", minValue)
	}

	if len(parts) == 1 {
		return 1, nil
	}

	maxValue, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid transaction id %q: %w", parts[1], err)
	}
	if maxValue < minValue {
		return 0, fmt.Errorf("invalid interval %q: end is before start", interval)
	}

	return maxValue - minValue + 1, nil
}

// uuidSetTransactionCount converts a single "uuid[:tag]:interval[:interval]..."
// entry into the number of transactions it represents.
func uuidSetTransactionCount(uuidSet string) (int64, error) {
	segments := strings.Split(uuidSet, ":")
	if len(segments) < 2 {
		return 0, fmt.Errorf("invalid gtid entry %q: missing uuid", uuidSet)
	}

	uuid := segments[0]
	if !uuidRE.MatchString(uuid) {
		return 0, fmt.Errorf("invalid uuid %q", uuid)
	}

	intervals := segments[1:]

	// An optional tag may follow the uuid. It's distinguished from an
	// interval by not matching the interval pattern (an empty tag, e.g.
	// "uuid::1-10", is also permitted by the GTID spec).
	if !intervalRE.MatchString(intervals[0]) {
		tag := intervals[0]
		if tag != "" && !tagRE.MatchString(tag) {
			return 0, fmt.Errorf("invalid tag %q", tag)
		}
		intervals = intervals[1:]
	}

	if len(intervals) == 0 {
		return 0, fmt.Errorf("invalid gtid entry %q: no intervals", uuidSet)
	}

	var count int64
	for _, interval := range intervals {
		n, err := intervalTransactionCount(interval)
		if err != nil {
			return 0, err
		}
		count += n
	}

	return count, nil
}

// TransactionCount takes a GTID_EXECUTED string and calculates the number of
// transactions it contains.
// - GTID gaps are also handled.
func TransactionCount(input string) (int64, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return 0, nil
	}

	// uuid_sets are comma separated; MySQL also pretty-prints with a
	// newline after each comma, so normalise whitespace away per entry
	// rather than splitting on newlines.
	entries := strings.Split(trimmed, ",")

	var count int64
	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		n, err := uuidSetTransactionCount(entry)
		if err != nil {
			return 0, err
		}
		count += n
	}

	return count, nil
}

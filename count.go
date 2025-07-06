package mysqlgtid

import (
	"strconv"
	"strings"
)

// convert <uuid>:<range> or <uuid>:<tag>:<range> into <range>.
func transactionPart(uuid_range string) string {
	parts := strings.Split(uuid_range, ":")
	if len(parts) == 2 {
		return parts[1]
	} else if len(parts) == 3 {
		return parts[2]
	}
	return "" // invalid so return nothing
}

// convert "1-100" -> 100, "1-10,100-110" -> 20, etc
func rangesToTxTransactionCount(ranges string) (int64, error) {
	var counter int64

	subRanges := strings.Split(ranges, ",")
	for _, r := range subRanges {
		if r == "" {
			// no change to counter as no data
		} else if !strings.Contains(r, "-") {
			// single value
			counter++
		} else {
			// find min and max
			minMax := strings.Split(r, "-")
			minValue, err := strconv.ParseInt(minMax[0], 10, 64)
			if err != nil {
				return 0, err
			}
			maxValue, err := strconv.ParseInt(minMax[1], 10, 64)
			if err != nil {
				return 0, err
			}
			counter = counter + maxValue - minValue + 1
		}
	}
	return counter, nil
}

// TransactionCount takes a GTID_EXECUTED string and calculates the number of transactions it contains.
// - GTID gaps are also handled.
func TransactionCount(input string) (int64, error) {
	var count int64

	// split by \n, handling windows usage too, just in case
	entries := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	// split by : keeping only 2nd part
	for _, entry := range entries {
		txRanges := transactionPart(entry)

		// handle gaps if present, e.g. split by ,
		txTransactionCount, err := rangesToTxTransactionCount(txRanges)
		if err != nil {
			return 0, err
		}
		count = count + txTransactionCount
	}

	return count, nil
}

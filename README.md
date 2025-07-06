## mysqlgtid - provide a GTID transaction counter

MySQL does not expose a transaction count very explicitly.  This routine
provides an easy way to calculate the number of transactions executed
on a system running with MySQL GTIDs.

Currently the code only works for the traditional GTID set and does not
yet handle the new GTID format introduced in MySQL 8.1.  This will be
handled later.

### Usage

```
package main

import (
	"fmt"
	"github.com/sjmudd/mysqlgtid"
)

func main() {
	gtidSet := "0a0bd206-4750-11ee-9a88-246e96822b80:1-100492614"

	count, err := mysqlgtid.TransactionCount(gtidSet)
	if err != nil {
		fmt.Errorf("TransactionCount(%q) failed: %v",
			gtidSet,
			err,
		)
		return
	}

	fmt.Printf("gtidSet: %q, count: %d\n", gtidSet, count)
}
```

### Licensing

BSD 2-Clause License

### Feedback

Feedback or patches welcome.

Simon J Mudd
<sjmudd@pobox.com>

### Code Documentation
[godoc.org/github.com/sjmudd/mysqlgtid](http://godoc.org/github.com/sjmudd/mysqlgtid)

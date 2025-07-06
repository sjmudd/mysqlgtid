## mysqlgtid - provide a GTID transaction counter

MySQL does not expose a transaction count very explicitly.  This routine
provides an easy way to calculate the number of transactions executed
on a system running with MySQL GTIDs.

### GTID Spec

This is the MySQL GTID Spec:
- https://dev.mysql.com/doc/refman/8.4/en/replication-gtids-concepts.html
- https://dev.mysql.com/doc/refman/9.3/en/replication-gtids-concepts.html (latest 9.X at time of writing this)

The current specs seem to imply that an *empty tag* is possible.
This would allow a format of line such as `<uuid>::<range>` which feels
a bit weird when you can just as well use `<uuid>:<range>` instead.

The code does not attempt to parse the GTID set rigourously, but parse
it as simply as possible. If there's a reason to be stricter in parsing
let me know.

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

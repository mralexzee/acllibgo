# acllibgo

Provides a method to scrub based on acl groups, zero, or keep property values of a struct. Library uses reflection and caching whenever possible. Intended for low frequency calls to ensure security compliance. Avoid using in high call frequency such as logging or tight loops.
- Good use cases: 
    - Scrubbing before returning object model by API
    - Scrubbing before persisting model to database
- Bad use cases: 
    - Scrubbing all log structures before writing to debug logs

### Functions

- Scrub(item, groups) -> zero out fields based on struct fields tag "acl" and provided groups
- Keep(item, fields) -> keep only the fields define in fields array, other fields get zero'd out
- Zero(item, fields) -> zero out all specified fields, leave others alone
- Parse(text) -> parses string to StructField array to pass into Keep and Zero

### Performance

Mid 2014 15" Macbook Pro i7 2.5GHz/16GB macOS 10.15
``` 
2020-06-02

Benchmark_Keep_Basic-8        	  465370	      2317 ns/op	     544 B/op	      18 allocs/op
Benchmark_ParseComplex-8      	  352665	      3202 ns/op	    1640 B/op	      50 allocs/op
Benchmark_ParseSimple-8       	 1409509	       851 ns/op	     448 B/op	      13 allocs/op
Benchmark_Scrub_NilAcl-8      	37911655	        32 ns/op	      16 B/op	       1 allocs/op
Benchmark_Scrub_EmptyAcl-8    	  285002	      4111 ns/op	    1568 B/op	      37 allocs/op
Benchmark_Scrub_SingleAcl-8   	  286686	      4131 ns/op	    1568 B/op	      37 allocs/op
Benchmark_Scrub_MultiAcl-8    	  264146	      4513 ns/op	    1568 B/op	      37 allocs/op
Benchmark_Zero_Basic-8        	  235779	      5025 ns/op	    2000 B/op	      59 allocs/op
```

### Playground:

https://play.golang.org/p/lDkvau0Ot1P

### Example:
```go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/mralexzee/acllibgo"
)

type User struct {
	Name   string
	Active bool
	Token  string `acl:"admin"`
	Groups []*Group
}

type Group struct {
	Name       string
	SystemMode string `acl:"admin"`
}

func main() {
	// Scrubbing the object with groups that include an admin
	usr := User{
		Name:   "John",
		Active: true,
		Token:  "secret-token",
		Groups: []*Group{
			&Group{
				"Users",
				"home-system"},
		},
	}

	scrubStruct(&usr, "admin", "user")

	// Scrubbing the object with groups that does not include admin
	usr = User{
		Name:   "John",
		Active: true,
		Token:  "secret-token",
		Groups: []*Group{
			&Group{
				"Users",
				"home-system"},
		},
	}
	scrubStruct(&usr, "user", "reporter")

	// Scrubbing array of objects
	usrList := []*User{&User{
		Name:   "John",
		Active: true,
		Token:  "secret-token",
		Groups: []*Group{
			&Group{
				"Users",
				"home-system"},
		},
	}, &User{
		Name:   "Jane Doe",
		Active: true,
		Token:  "secret-token",
		Groups: []*Group{
			&Group{
				"Users",
				"home-system"},
		},
	},
	}
	scrubStruct(usrList, "user", "reporter")

}

func scrubStruct(i interface{}, groups ...string) {
	fmt.Printf("Scrubbing object based on groups: %v\n", groups)
	fmt.Println("Before:   => " + toJson(i))
	acllibgo.Scrub(i, groups)
	fmt.Println("Scrubbed: => " + toJson(i) + "\n")
}

func toJson(i interface{}) string {
	j, _ := json.Marshal(i)
	return string(j)
}

```

```bash

Scrubbing object based on groups: [admin user]
Before:   => {"Name":"John","Active":true,"Token":"secret-token","Groups":[{"Name":"Users","SystemMode":"home-system"}]}
Scrubbed: => {"Name":"John","Active":true,"Token":"secret-token","Groups":[{"Name":"Users","SystemMode":"home-system"}]}

Scrubbing object based on groups: [user reporter]
Before:   => {"Name":"John","Active":true,"Token":"secret-token","Groups":[{"Name":"Users","SystemMode":"home-system"}]}
Scrubbed: => {"Name":"John","Active":true,"Token":"","Groups":[{"Name":"Users","SystemMode":""}]}

Scrubbing object based on groups: [user reporter]
Before:   => [{"Name":"John","Active":true,"Token":"secret-token","Groups":[{"Name":"Users","SystemMode":"home-system"}]},{"Name":"Jane Doe","Active":true,"Token":"secret-token","Groups":[{"Name":"Users","SystemMode":"home-system"}]}]
Scrubbed: => [{"Name":"John","Active":true,"Token":"","Groups":[{"Name":"Users","SystemMode":""}]},{"Name":"Jane Doe","Active":true,"Token":"","Groups":[{"Name":"Users","SystemMode":""}]}]

```
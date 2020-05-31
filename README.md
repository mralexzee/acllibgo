# acllibgo

Provides a method to scrub property values of a struct based on acl groups. Library uses reflection and caching whenever possible. Intended for low frequency calls to ensure security compliance. Avoid using in high call frequency such as logging or tight loops.
- Good use cases: 
    - Scrubbing before returning object model by API ( ~ 200/second )
    - Scrubbing before persisting model to database ( ~ 200/second )
- Bad use cases: 
    - Scrubbing before writing to debug logs ( ~ 1000+/second )

### Performance

Mid 2014 15" Macbook Pro i7 2.5GHz/16GB macOS 10.15
``` 
2020-05-31
Benchmark_ScrubNilAcl-8      	32795670	        31.8 ns/op	      16 B/op	       1 allocs/op
Benchmark_ScrubEmptyAcl-8    	  279805	      4166 ns/op	    1568 B/op	      37 allocs/op
Benchmark_ScrubSingleAcl-8   	  282964	      4214 ns/op	    1568 B/op	      37 allocs/op
Benchmark_ScrubMultiAcl-8    	  262330	      4485 ns/op	    1568 B/op	      37 allocs/op
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
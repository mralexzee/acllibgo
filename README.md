# acllibgo
Library to help with security

Provides a method to scrub property values of the struct based on acl groups.

###Playground:

https://play.golang.org/p/lDkvau0Ot1P

###Example:
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
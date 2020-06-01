// Copyright 2020 Alexander Zherdev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acllibgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStruct_Nil(t *testing.T) {

	err := Scrub(nil, []string{})
	assert.Error(t, err)
	var testItem *Person
	err = Scrub(testItem, []string{"test"})
	assert.Error(t, err)
	err = Scrub(&testItem, []string{"test"})
	assert.Error(t, err)
}

func TestStruct_Basic(t *testing.T) {

	testItem := newPerson()

	// Pass by value - not ok
	err := Scrub(testItem, []string{})
	assert.Error(t, err)

	// Pass by reference - ok
	err = Scrub(&testItem, []string{})
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, int32(0), testItem.Height)
	assert.NotNil(t, testItem.Father)
	assert.NotNil(t, testItem.Children)
	assert.Nil(t, testItem.FullName)
	assert.Nil(t, testItem.Mother)
	assert.Nil(t, testItem.FullName)

	assert.True(t, len(testItem.Children) > 0)
	assert.Equal(t, int32(0), testItem.Children[0].Height)
	assert.Equal(t, int32(0), testItem.Children[1].Height)
	assert.Nil(t, testItem.Children[0].FullName)
	assert.Nil(t, testItem.Children[0].FullName)
}

func TestStruct_Slice(t *testing.T) {

	one := newPerson()
	two := newPerson()
	three := newPerson()

	testItems := []*Person{&one, &two, &three}

	// Pass by value - not ok
	err := Scrub(testItems, []string{})
	assert.NoError(t, err)
}

func TestStruct_Map(t *testing.T) {

	one := newPerson()
	two := newPerson()
	three := newPerson()

	testItems := map[string]*Person{"one": &one, "two": &two, "three": &three}

	// Pass by value - not ok
	err := Scrub(testItems, []string{})
	assert.NoError(t, err)
}

func newPerson() Person {
	return Person{
		Created:   time.Now().UTC(),
		Birthdate: time.Now().Truncate(time.Hour * 24).Add(-24 * 365 * 7 * time.Hour),
		Age:       21,
		Height:    68,
		Groups:    map[string]bool{"chessclub": true, "pianoleague": false},
		FullName:  []string{"John", "Smith", "Doe"},
		Nickname:  "John",
		Mother: &Person{
			Age:      78,
			Height:   64,
			FullName: []string{"Penny", "Angela", "Smith"},
			Nickname: "Penny",
		},
		Father: &Person{
			Age:      82,
			Height:   71,
			FullName: []string{"Anthony", "Smith", "Sr"},
			Nickname: "Tony",
		},
		Children: []*Person{
			&Person{
				Age:      7,
				Height:   45,
				FullName: []string{"Johnny", "Knox", "Jr"},
				Nickname: "Johnny Boy",
			},
			&Person{
				Age:      11,
				Height:   49,
				FullName: []string{"Cindy", "Lou"},
				Nickname: "Sin",
			}},
		PetCat: Cat{
			Name: "Fluffy",
			Type: "furry",
		},
		Friends: map[string]*Person{
			"best": &Person{
				Age:      34,
				Height:   68,
				FullName: []string{"Philarmon", "Carter"},
				Nickname: "Phil",
			},
			"sweetheart": &Person{
				Age:      34,
				Height:   68,
				FullName: []string{"Becky", "Hair"},
				Nickname: "Becks",
			},
		},
	}
}

type Person struct {
	Age       int                `json:"age,omitempty"`
	Height    int32              `json:"height,omitempty" acl:"tester"`
	Groups    map[string]bool    `json:"groups,omitempty"  acl:"tester"`
	FullName  []string           `json:"fullName,omitempty" acl:"tester"`
	Nickname  string             `json:"nickname,omitempty" acl:"*"`
	Mother    *Person            `json:"mother,omitempty" acl:"tester"`
	Father    *Person            `json:"father,omitempty"`
	Children  []*Person          `json:"children,omitempty"`
	PetCat    Cat                `json:"petCat,omitempty" acl:"tester"`
	Friends   map[string]*Person `json:"friends,omitempty"`
	Created   time.Time          `json:"created,omitempty"`
	Birthdate time.Time          `json:"birthdate,omitempty" acl:"tester"`
}

type Cat struct {
	Name string `json:"name" acl:"root,account"`
	Type string `json:"type" acl:"root"`
}

func TestInvalidType_Int(t *testing.T) {

	var testVar int = 10
	err := Scrub(testVar, []string{"owner"})
	assert.Error(t, err)
	err = Scrub(&testVar, []string{"owner"})
	assert.Error(t, err)

}

func TestInvalidType_Float(t *testing.T) {

	var testVar float32 = 10
	err := Scrub(testVar, []string{"owner"})
	assert.Error(t, err)
	err = Scrub(&testVar, []string{"owner"})
	assert.Error(t, err)

}

func TestInvalidType_Bool(t *testing.T) {

	var testVar bool = true
	err := Scrub(testVar, []string{"owner"})
	assert.Error(t, err)
	err = Scrub(&testVar, []string{"owner"})
	assert.Error(t, err)

}

func TestInvalidType_InterfaceInderection(t *testing.T) {

	var origVar bool = true
	var testVar interface{} = origVar

	err := Scrub(testVar, []string{"owner"})
	assert.Error(t, err)
	err = Scrub(&testVar, []string{"owner"})
	assert.Error(t, err)
}

func TestInvalidType_Slice(t *testing.T) {

	var testVar []int = []int{1, 2, 3, 4}

	err := Scrub(testVar, []string{"owner"})
	assert.Error(t, err)
	err = Scrub(&testVar, []string{"owner"})
	assert.Error(t, err)
}

func TestInvalidType_Map(t *testing.T) {

	var testVar map[string]int = map[string]int{"one": 1, "two": 2, "three": 3}

	err := Scrub(testVar, []string{"owner"})
	assert.Error(t, err)
	err = Scrub(&testVar, []string{"owner"})
	assert.Error(t, err)
}

func Benchmark_ScrubNilAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, nil)
	}
}

func Benchmark_ScrubEmptyAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, []string{})
	}
}

func Benchmark_ScrubSingleAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, []string{"root"})
	}
}

func Benchmark_ScrubMultiAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, []string{"access", "login"})
	}
}

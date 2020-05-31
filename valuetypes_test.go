// Copyright 2020 Alexander Zherdev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acllibgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

// Copyright 2020 Alexander Zherdev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acllibgo

import "testing"

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

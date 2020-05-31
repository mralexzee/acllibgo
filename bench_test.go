package acllibgo

import "testing"

func BenchmarkScrubNilAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, nil)
	}
}

func BenchmarkScrubEmptyAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, []string{})
	}
}

func BenchmarkScrubSingleAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, []string{"root"})
	}
}

func Benchmark_Scrub_MultiAcl(b *testing.B) {
	testItem := newPerson()
	for n := 0; n < b.N; n++ {
		Scrub(&testItem, []string{"access", "login"})
	}
}

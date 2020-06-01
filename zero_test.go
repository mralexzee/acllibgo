package acllibgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Zero_Basic(t *testing.T) {

	testItem := newPerson()

	// Pass by value - not ok
	err := Zero(testItem, []StructField{{Name: "Nickname", Fields: nil}})
	assert.Error(t, err)

	// Pass by reference - ok
	err = Zero(&testItem, []StructField{
		{Name: "Age"},
		{Name: "Nickname"},
		{Name: "Children", Fields: []StructField{{Name: "Age"}}},
		{Name: "Mother"},
		{Name: "Friends", Fields: []StructField{{Name: "*"}}},
	})
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.True(t, testItem.Height > 0)
	assert.NotNil(t, testItem.Father)

	assert.True(t, testItem.Age == 0)
	assert.Nil(t, testItem.Mother)
	assert.Nil(t, testItem.Friends)

	assert.NotNil(t, testItem.Children)
	assert.True(t, len(testItem.Children) > 0)
	assert.True(t, testItem.Children[0].Age == 0)
	assert.True(t, testItem.Children[1].Age == 0)
	assert.True(t, testItem.Children[0].Height > 0)
	assert.True(t, testItem.Children[1].Height > 0)
	assert.True(t, len(testItem.Children[0].FullName) > 0)
	assert.True(t, len(testItem.Children[1].FullName) > 0)
}

func Benchmark_Zero_Basic(b *testing.B) {
	testItem := newPerson()
	fields := []StructField{
		{Name: "Age"},
		{Name: "Nickname"},
		{Name: "Children", Fields: []StructField{{Name: "Age"}}},
		{Name: "Mother", Fields: []StructField{{Name: "*"}}},
	}
	for n := 0; n < b.N; n++ {
		Zero(&testItem, fields)
	}
}

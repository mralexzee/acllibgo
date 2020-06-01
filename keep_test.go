package acllibgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Keep_Basic(t *testing.T) {

	testItem := newPerson()

	// Pass by value - not ok
	err := Keep(testItem, []StructField{{Name: "Nickname", Fields: nil}})
	assert.Error(t, err)

	// Pass by reference - ok
	err = Keep(&testItem, []StructField{
		{Name: "Age"},
		{Name: "Nickname"},
		{Name: "Children", Fields: []StructField{{Name: "Age"}}},
		{Name: "Mother", Fields: []StructField{{Name: "*"}}},
	})
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, int32(0), testItem.Height)
	assert.Nil(t, testItem.Father)
	assert.Nil(t, testItem.FullName)

	assert.NotNil(t, testItem.Mother)
	assert.True(t, testItem.Mother.Age > 0)
	assert.True(t, len(testItem.Mother.Nickname) > 0)

	assert.NotNil(t, testItem.Children)
	assert.True(t, len(testItem.Children) > 0)
	assert.Equal(t, int32(0), testItem.Children[0].Height)
	assert.Equal(t, int32(0), testItem.Children[1].Height)
	assert.Nil(t, testItem.Children[0].FullName)
	assert.Nil(t, testItem.Children[1].FullName)
}

func Benchmark_Keep_Basic(b *testing.B) {
	testItem := newPerson()
	fields := []StructField{
		{Name: "Age"},
		{Name: "Nickname"},
		{Name: "Children", Fields: []StructField{{Name: "Age"}}},
		{Name: "Mother", Fields: []StructField{{Name: "*"}}},
	}
	for n := 0; n < b.N; n++ {
		Keep(&testItem, fields)
	}
}

package acllibgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse_NoPanic(t *testing.T) {
	samples := []string{
		"",
		"id",
		"id,password",
		"id,password,account(username,type)",
		"id,password,account(username,type),group(key,name)",
		"id,password,account(username,type,parent(id,name)),group(key,name)",
		"   ",
		",,(),,",
		")",
		",",
		"(",
		"aaaa(",
		"aaaa(aaaa",
		"(aaaaa",
		"),,,,,",
		",,,)",
		"(,,)",
	}

	for _, text := range samples {
		assert.NotPanics(t, func() { Parse(text) })
	}
}

func Test_Parse_ValidateComplex(t *testing.T) {
	propString := "id,password,account(username,type,parent(id,name)),group(key,name)"

	r, e := Parse(propString)

	assert.Nil(t, e)
	assert.NotNil(t, r)
	assert.True(t, len(r) == 4)
	assert.True(t, r[0].Name == "id")
	assert.True(t, r[1].Name == "password")

	assert.True(t, r[2].Name == "account")
	assert.True(t, len(r[2].Fields) > 0)
	assert.True(t, r[2].Fields[0].Name == "username")
	assert.True(t, r[2].Fields[1].Name == "type")
	assert.True(t, r[2].Fields[2].Name == "parent")
	assert.True(t, len(r[2].Fields[2].Fields) > 0)
	assert.True(t, r[2].Fields[2].Fields[0].Name == "id")
	assert.True(t, r[2].Fields[2].Fields[1].Name == "name")

	assert.True(t, r[3].Name == "group")
	assert.True(t, len(r[3].Fields) > 0)
	assert.True(t, r[3].Fields[0].Name == "key")
	assert.True(t, r[3].Fields[1].Name == "name")
}

func Test_Parse_EmptyInput(t *testing.T) {
	propString := ""

	r, e := Parse(propString)

	assert.Nil(t, e)
	assert.NotNil(t, r)
	assert.True(t, len(r) == 0)
}

func Test_Parse_BadInput_NoCloseP(t *testing.T) {
	propString := ",("

	r, e := Parse(propString)

	assert.NotNil(t, e)
	assert.NotNil(t, r)
	assert.True(t, len(r) == 0)
}

func Test_Parse_BadInput_NoOpenP(t *testing.T) {
	propString := ",)"

	r, e := Parse(propString)

	assert.NotNil(t, e)
	assert.NotNil(t, r)
	assert.True(t, len(r) == 0)
}

func Test_Parse_BadInput_BadParenthesis(t *testing.T) {
	propString := "id,account((id,name)"

	r, e := Parse(propString)

	assert.NotNil(t, e)
	assert.NotNil(t, r)
	assert.True(t, len(r) == 0)
}

func Benchmark_ParseComplex(b *testing.B) {
	str := "id,password,account(username,type,parent(id,name)),group(key,name)"
	for n := 0; n < b.N; n++ {
		Parse(str)
	}
}

func Benchmark_ParseSimple(b *testing.B) {
	str := "id,password,account"
	for n := 0; n < b.N; n++ {
		Parse(str)
	}
}

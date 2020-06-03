package acllibgo

import (
	"errors"
	"strings"
)

// Parse converts string representing property names to array of StructField objects
// Example Text: id,password,account(username,type,parent(id,name)),group(key,name)
func Parse(text string) ([]StructField, error) {
	runeString := []rune(strings.TrimSpace(text))
	runLen := len(runeString)

	oc := strings.Count(text, "(")
	cc := strings.Count(text, ")")

	if oc != cc {
		return []StructField{}, errors.New("parse: parenthesis count mismatch")
	}

	if runLen == 0 {
		return []StructField{}, nil
	}

	rv, _, e := parseText(runeString, 0, runLen)

	return rv, e
}

func parseText(text []rune, offset int, strLen int) ([]StructField, int, error) {
	if offset >= strLen {
		return []StructField{}, offset, errors.New("parse: offset out of bounds")
	}

	rv := make([]StructField, 0)
	var queue []rune
	for x := offset; x < strLen; x++ {
		r := text[x]
		switch r {
		case ',':
			cleanString := string(queue)
			if len(cleanString) > 0 {
				rv = append(rv, StructField{Name: cleanString})
			}
			queue = make([]rune, 0)
		case '(':
			x++
			cleanString := string(queue)
			if len(cleanString) > 0 {
				children, newOffset, _ := parseText(text, x, strLen)
				rv = append(rv, StructField{Name: string(queue), Fields: children})
				queue = make([]rune, 0)
				x = newOffset
			}
		case ')':
			cleanString := string(queue)
			if len(cleanString) > 0 {
				rv = append(rv, StructField{Name: cleanString})
			}
			return rv, x, nil
		case ' ':
			// do nothing - we remove spaces
			continue
		default:
			queue = append(queue, r)
		}
	}

	cleanString := string(queue)
	if len(cleanString) > 0 {
		rv = append(rv, StructField{Name: cleanString})
	}

	return rv, strLen, nil
}

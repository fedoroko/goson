package goson

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// iPatternGenerator generates a string according to a pattern: <0\0\0>
// where the first number is a number of letters
// the second is a number of digits
// and the third is a number of special characters.
// It allowed to discard unused number from the end,
// so <0\0> and <0> are valid too.
// Also, it supports prefix and suffix.
// For example pattern "<5\1>@test.com"
// will generate a string of length 6,
// that contains exact 5 letter and exact 1 digit
// and after that it will append suffix "@test.com".
type iPatternGenerator interface {
	Generate() string
}

type strPattern struct {
	prefix, suffix string
	letters        int
	digits         int
	specials       int
	length         int
}

type generator struct {
	pattern strPattern
}

// Generate concatenates prefix, generated value and suffix
func (generator generator) Generate() string {
	return fmt.Sprintf(
		"%s%s%s",
		generator.pattern.prefix,
		generator.generate(),
		generator.pattern.suffix,
	)
}

// Character pool: 52 letters, 10 digits and 19 specials.
const bytePoolLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const bytePoolDigits = "0123456789"
const bytePoolSpecials = "~!@#$%^&*()_+-={}[]"

// generate takes a random character from a specific byte pool and insert it into a slice of bytes.
// At the end shuffling a slice of bytes.
func (generator generator) generate() string {
	n := generator.pattern.length
	bytes := make([]byte, n)

	var i int
	for j := 0; j < generator.pattern.letters; j++ {
		char := getRandomCharFromSource(bytePoolLetters)
		bytes[i] = char
		i++
	}

	for j := 0; j < generator.pattern.digits; j++ {
		char := getRandomCharFromSource(bytePoolDigits)
		bytes[i] = char
		i++
	}

	for j := 0; j < generator.pattern.specials; j++ {
		char := getRandomCharFromSource(bytePoolSpecials)
		bytes[i] = char
		i++
	}

	rand.Shuffle(n, func(i, j int) {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	})

	return string(bytes)
}

func getRandomCharFromSource(source string) byte {
	n := len(source)
	char := source[rand.Intn(n)]

	return char
}

func newPatternGenerator(raw string) (iPatternGenerator, error) {
	pattern, err := buildPattern(raw)
	if err != nil {
		return nil, throwInvalidPatternError(raw, err)
	}

	return &generator{
		pattern: pattern,
	}, nil
}

func buildPattern(raw string) (strPattern, error) {
	var pattern strPattern

	start, end := locatePattern(raw)
	if start == -1 || end == -1 || start > end {
		return pattern, errors.New("can't locate pattern")
	}

	patternBody := raw[start+1 : end]
	pattern, err := parsePattern(patternBody)
	if err != nil {
		return pattern, err
	}

	pattern.prefix = unescapeString(raw[:start])
	pattern.suffix = unescapeString(raw[end+1:])

	return pattern, nil
}

func unescapeString(raw string) string {
	res := make([]byte, 0)
	for i := 0; i < len(raw); i++ {
		if i < len(raw) && raw[i] == '\\' && raw[i+1] != '\\' {
			continue
		}

		res = append(res, raw[i])
	}

	return string(res)
}

func locatePattern(raw string) (int, int) {
	start, end := -1, -1
	for i, char := range raw {
		switch {
		case i > 0 && raw[i-1] == '\\':
			continue
		case char == '<' && start == -1:
			start = i
		case char == '<' && start != -1:
			start = -1
		case char == '>' && end == -1:
			end = i
		case char == '>' && end != -1:
			end = -1
		}
	}

	return start, end
}

const (
	patternLegendLetter = iota
	patternLegendDigit
	patternLegendSpecial
)

const patternSeparator = "/"
const patternMaxLen = 64

func parsePattern(body string) (strPattern, error) {
	var pattern strPattern
	exploded := strings.Split(body, patternSeparator)
	if len(exploded) > 3 {
		return pattern, errors.New("can't split pattern correctly")
	}

	for i, str := range exploded {
		num, err := strconv.Atoi(str)
		if err != nil {
			return pattern, err
		}
		switch i {
		case patternLegendLetter:
			pattern.letters = num
			pattern.length += num
		case patternLegendDigit:
			pattern.digits = num
			pattern.length += num
		case patternLegendSpecial:
			pattern.specials = num
			pattern.length += num
		}
	}

	if pattern.length > patternMaxLen {
		return pattern, errors.New(
			fmt.Sprintf(
				"max pattern len is %d. actual: %d",
				patternMaxLen,
				pattern.length,
			),
		)
	}

	return pattern, nil
}

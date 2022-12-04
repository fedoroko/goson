package goson

import (
	"encoding/json"
	"reflect"
	"strings"
)

// iParser parses an input bytes and transforms it into a ParsedFields
type iParser interface {
	Parse([]byte) (map[string]ParsedField, error)
}

type parser struct {
	keywordSet map[string]NewFieldFunc
	rawFields  map[string]interface{} // original fields as key:value map. need it for a referenceFields
}

func (parser *parser) Parse(body []byte) (map[string]ParsedField, error) {
	return parser.parse(body)
}

func (parser *parser) parse(body []byte) (map[string]ParsedField, error) {
	parser.rawFields = make(map[string]interface{})
	if err := json.Unmarshal(body, &parser.rawFields); err != nil {
		return nil, err
	}

	fields := make(map[string]ParsedField)
	for key, value := range parser.rawFields {
		field, err := parser.parseField(key, value)
		if err != nil {
			return nil, err
		}
		fields[key] = field
	}

	return fields, nil
}

func (parser *parser) parseField(key string, rawField interface{}) (ParsedField, error) {
	var field ParsedField
	var err error

	rawType, rawValue := reflect.TypeOf(rawField), reflect.ValueOf(rawField)

	switch rawType.Kind() {
	case reflect.String:
		field, err = parser.parseStringField(key, rawValue.String())
	default:
		field = newStaticField(key, rawField)
	}

	return field, err
}

const goPrefix = "_go:"

func (parser *parser) parseStringField(key, value string) (ParsedField, error) {
	if parser.isStaticString(value) {
		return newStaticField(key, value), nil
	}

	value = value[len(goPrefix):]
	if parser.isValidKeywordString(value) {
		fn := parser.keywordSet[value]
		return fn(key, value)
	}

	if parser.isValidReferenceString(value) {
		return newReferenceField(key, value[1:]), nil
	}

	return newPatternField(key, value)
}

func (parser *parser) isStaticString(value string) bool {
	return !strings.Contains(value, goPrefix)
}

func (parser *parser) isValidKeywordString(value string) bool {
	for fieldKey := range parser.keywordSet {
		if value == fieldKey {
			return true
		}
	}

	return false
}

func (parser *parser) isValidReferenceString(value string) bool {
	if _, ok := parser.rawFields[value[1:]]; ok {
		return true
	}

	return false
}

func newParser() iParser {
	defaultKeywordSet := getDefaultKeywordSet()
	return &parser{
		keywordSet: defaultKeywordSet,
	}
}

func getDefaultKeywordSet() map[string]NewFieldFunc {
	return map[string]NewFieldFunc{
		"timestamp": newTimestampField,
		"uuid":      newUUIDField,
	}
}

func newParserWithCustomKeywords(fields map[string]NewFieldFunc) iParser {
	customKeywordSet := getDefaultKeywordSet()
	for key, fn := range fields {
		customKeywordSet[key] = fn
	}

	return &parser{
		keywordSet: customKeywordSet,
	}
}

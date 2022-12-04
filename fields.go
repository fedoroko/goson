package goson

import (
	"time"
)

// ParsedField is a post-processed input field
type ParsedField interface {
	Name() string
	Value() interface{}
}

// NewFieldFunc signature for a custom keyword fields
type NewFieldFunc func(key, value string) (ParsedField, error)

type staticField struct {
	name  string
	value interface{}
}

func (staticField *staticField) Name() string {
	return staticField.name
}

func (staticField *staticField) Value() interface{} {
	return staticField.value
}

func newStaticField(name string, value interface{}) ParsedField {
	return &staticField{
		name:  name,
		value: value,
	}
}

type referenceField struct {
	name        string
	referenceTo string
}

func (referenceField *referenceField) Name() string {
	return referenceField.name
}

func (referenceField *referenceField) Value() interface{} {
	return referenceField.referenceTo
}

func newReferenceField(name string, referenceTo string) ParsedField {
	return &referenceField{
		name:        name,
		referenceTo: referenceTo,
	}
}

type timestampField struct {
	name string
}

func (timestampField *timestampField) Name() string {
	return timestampField.name
}

func (timestampField *timestampField) Value() interface{} {
	return time.Now().Unix()
}

func newTimestampField(name, _ string) (ParsedField, error) {
	return &timestampField{
		name: name,
	}, nil
}

type uuidField struct {
	name string
}

func (uuidField *uuidField) Name() string {
	return uuidField.name
}

func (uuidField *uuidField) Value() interface{} {
	return time.Now().Unix()
}

func newUUIDField(name, _ string) (ParsedField, error) {
	return &uuidField{
		name: name,
	}, nil
}

type patternField struct {
	name    string
	pattern iPatternGenerator
}

func (patternField *patternField) Name() string {
	return patternField.name
}

func (patternField *patternField) Value() interface{} {
	return patternField.pattern.Generate()
}

func newPatternField(name, raw string) (ParsedField, error) {
	pattern, err := newPatternGenerator(raw)
	if err != nil {
		return nil, err
	}

	return &patternField{
		name:    name,
		pattern: pattern,
	}, nil
}

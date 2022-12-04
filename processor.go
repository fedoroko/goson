package goson

import (
	"encoding/json"
	"math/rand"
	"time"
)

// Processor constructs a parser and force it to parse an input bytes
type Processor interface {
	Generate() []byte
}

type processor struct {
	fields map[string]ParsedField
}

// Generate generates and returns a record according to an input fields
func (processor *processor) Generate() []byte {
	fields := make(map[string]interface{})
	for _, field := range processor.fields {
		fields[field.Name()] = field.Value()
	}

	for _, field := range processor.fields {
		switch field.(type) {
		case *referenceField:
			fields[field.Name()] = fields[field.Value().(string)]
		}
	}

	b, err := json.Marshal(fields)
	if err != nil {
		panic(err)
	}

	return b
}

// New builds a Processor with a default set of keywords
func New(body []byte) (Processor, error) {
	defaultParser := newParser()
	fields, err := defaultParser.Parse(body)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().Unix())

	return &processor{
		fields: fields,
	}, nil
}

// NewWithCustomKeywords allows to build a Processor that can parse and handle a custom set of keywords
func NewWithCustomKeywords(body []byte, keywordSet map[string]NewFieldFunc) (Processor, error) {
	defaultParser := newParserWithCustomKeywords(keywordSet)
	fields, err := defaultParser.Parse(body)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().Unix())

	return &processor{
		fields: fields,
	}, nil
}

package goson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{
			name: "positive",
			input: []byte(`
				{
					"email": "_go:<5/1>@test.com",
					"password": "_go:<5/2/1>",
					"confirm_password": "_go:_password",
					"escaped_field": "_go:\\<0\\> <5>",
					"timestamp": "_go:timestamp"
				}
			`),
		},
		{
			name: "invalid pattern #1",
			input: []byte(`
				{
					"password": "_go:<5/2/1",
				}
			`),
			wantErr: true,
		},
		{
			name: "invalid pattern #2",
			input: []byte(`
				{
					"password": "_go:<5/2/1/>",
				}
			`),
			wantErr: true,
		},
		{
			name: "non numeric value #1",
			input: []byte(`
				{
					"password": "_go:<5/2/a>",
				}
			`),
			wantErr: true,
		},
		{
			name: "non numeric value #2",
			input: []byte(`
				{
					"password": "_go:<>",
				}
			`),
			wantErr: true,
		},
		{
			name: "pattern is too long",
			input: []byte(`
				{
					"password": "_go:<65>",
				}
			`),
			wantErr: true,
		},
		{
			name: "invalid json",
			input: []byte(`
				{
					"password": "_go:<5/2/1>",
				
			`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type testDummyField struct {
	name string
}

func (testDummyField *testDummyField) Name() string {
	return testDummyField.name
}

func (testDummyField *testDummyField) Value() interface{} {
	return "dummy"
}

func newTestDummyField(key, _ string) (ParsedField, error) {
	return &testDummyField{
		name: key,
	}, nil
}

func TestNewWithCustomKeywords(t *testing.T) {
	keywordSet := getDefaultKeywordSet()
	keywordSet["dummy"] = newTestDummyField
	body := []byte(`{"dummy":"_go:dummy"}`)
	expect := []byte(`{"dummy":"dummy"}`)

	t.Run("positive", func(t *testing.T) {
		testProcessor, err := NewWithCustomKeywords(body, keywordSet)
		require.NoError(t, err)
		got := testProcessor.Generate()
		assert.Equal(t, expect, got)
	})
}

func Test_processor_Generate(t *testing.T) {
	seedTestDate()
	defaultBody := []byte(`
		{
			"email": "_go:<5/1>@test.com",
			"password": "_go:<5/2/1>",
			"confirm_password": "_go:_password",
			"escaped_field": "_go:\\<0\\> <5>",
			"timestamp": "_go:timestamp"
		}
	`)
	defaultProcessor, err := New(defaultBody)
	t.Run("positive default processor", func(t *testing.T) {
		require.NoError(t, err)
		var prev []byte
		for i := 0; i < 10; i++ {
			got := defaultProcessor.Generate()
			assert.NotEqual(t, prev, got)
		}
	})

	customBody := []byte(`{"dummy": "_go:dummy"}`)
	expect := []byte(`{"dummy":"dummy"}`)
	keywordSet := getDefaultKeywordSet()
	keywordSet["dummy"] = newTestDummyField

	customProcessor, err := NewWithCustomKeywords(customBody, keywordSet)
	t.Run("positive custom processor", func(t *testing.T) {
		require.NoError(t, err)
		got := customProcessor.Generate()
		assert.Equal(t, expect, got)
	})
}

package goson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestParserStruct() parser {
	return parser{
		keywordSet: getDefaultKeywordSet(),
		rawFields: map[string]interface{}{
			"email":            "_go:<4/1>@test.com",
			"password":         "_go:<5/1/1>",
			"confirm_password": "_go:_password",
			"time":             "timestamp",
			"escaped_field":    "_go:\\<0\\> <5>",
		},
	}
}

func Test_parser_isStaticString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "positive",
			input: "some string",
			want:  true,
		},
		{
			name:  "negative",
			input: "_go:some string",
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testParser.isStaticString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parser_isValidKeywordString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "positive #1",
			input: "timestamp",
			want:  true,
		},
		{
			name:  "positive #1",
			input: "uuid",
			want:  true,
		},
		{
			name:  "negative #1",
			input: "city",
			want:  false,
		},
		{
			name:  "negative #2",
			input: "name",
			want:  false,
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testParser.isValidKeywordString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parser_isValidReferenceString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "positive",
			input: "_password",
			want:  true,
		},
		{
			name:  "negative #1",
			input: "_state",
		},
		{
			name:  "negative #2",
			input: "password",
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testParser.isValidReferenceString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parser_parseStringField(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    ParsedField
		wantErr bool
	}{
		{
			name: "positive #1",
			args: args{
				"email", "_go:<4/1>@test.com",
			},
			want: &patternField{
				name: "email",
				pattern: &generator{
					pattern: strPattern{
						suffix:  "@test.com",
						letters: 4,
						digits:  1,
						length:  5,
					},
				},
			},
		},
		{
			name: "positive #2",
			args: args{
				"time", "_go:timestamp",
			},
			want: &timestampField{
				name: "time",
			},
		},
		{
			name: "positive #3",
			args: args{
				"confirm_password", "_go:_password",
			},
			want: &referenceField{
				name:        "confirm_password",
				referenceTo: "password",
			},
		},
		{
			name: "invalid reference",
			args: args{
				"confirm_password", "_go:_name",
			},
			wantErr: true,
		},
		{
			name: "invalid keyword",
			args: args{
				"confirm_password", "_go:state",
			},
			wantErr: true,
		},
		{
			name: "invalid pattern #1",
			args: args{
				"password", "_go:<5/1/>",
			},
			wantErr: true,
		},
		{
			name: "invalid pattern #2",
			args: args{
				"password", "_go:<5/1/c>",
			},
			wantErr: true,
		},
		{
			name: "invalid pattern #3",
			args: args{
				"password", "_go:<65>",
			},
			wantErr: true,
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testParser.parseStringField(tt.args.key, tt.args.value)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_parser_parseField(t *testing.T) {
	type args struct {
		key      string
		rawField interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    ParsedField
		wantErr bool
	}{
		{
			name: "positive #1",
			args: args{
				"email", "_go:<4/1>@test.com",
			},
			want: &patternField{
				name: "email",
				pattern: &generator{
					pattern: strPattern{
						suffix:  "@test.com",
						letters: 4,
						digits:  1,
						length:  5,
					},
				},
			},
		},
		{
			name: "positive #2",
			args: args{
				"time", "_go:timestamp",
			},
			want: &timestampField{
				name: "time",
			},
		},
		{
			name: "positive #3",
			args: args{
				"confirm_password", "_go:_password",
			},
			want: &referenceField{
				name:        "confirm_password",
				referenceTo: "password",
			},
		},
		{
			name: "invalid reference",
			args: args{
				"confirm_password", "_go:_name",
			},
			wantErr: true,
		},
		{
			name: "invalid keyword",
			args: args{
				"confirm_password", "_go:state",
			},
			wantErr: true,
		},
		{
			name: "invalid pattern #1",
			args: args{
				"password", "_go:<5/1/>",
			},
			wantErr: true,
		},
		{
			name: "invalid pattern #2",
			args: args{
				"password", "_go:<5/1/c>",
			},
			wantErr: true,
		},
		{
			name: "invalid pattern #3",
			args: args{
				"password", "_go:<65>",
			},
			wantErr: true,
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testParser.parseField(tt.args.key, tt.args.rawField)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_parser_parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]ParsedField
		wantErr bool
	}{
		{
			name: "positive #1",
			input: []byte(`
				{
					"email": "_go:<4/1>@test.com",
					"password": "_go:<5/2/1>",
					"confirm_password": "_go:_password",
					"escaped_field": "_go:\\<0\\> <5>"
				}
			`),
			want: map[string]ParsedField{
				"confirm_password": &referenceField{
					name:        "confirm_password",
					referenceTo: "password",
				},
				"email": &patternField{
					name: "email",
					pattern: &generator{
						strPattern{
							suffix:  "@test.com",
							letters: 4,
							digits:  1,
							length:  5,
						},
					},
				},
				"escaped_field": &patternField{
					name: "escaped_field",
					pattern: &generator{
						strPattern{
							prefix:  "<0> ",
							letters: 5,
							length:  5,
						},
					},
				},
				"password": &patternField{
					name: "password",
					pattern: &generator{
						strPattern{
							letters:  5,
							digits:   2,
							specials: 1,
							length:   8,
						},
					},
				},
			},
		},
		{
			name: "invalid pattern",
			input: []byte(`
				{
					"email": "_go:<4/1@test.com",
					"password": "_go:<5/2/1>",
					"confirm_password": "_go:_password"
				}
			`),
			wantErr: true,
		},
		{
			name: "unescaped suffix",
			input: []byte(`
				{
					"email": "_go:<4/1><@test.com>",
				}
			`),
			wantErr: true,
		},
		{
			name: "pattern is too long",
			input: []byte(`
				{
					"email": "_go:<67>@test.com",
				}
			`),
			wantErr: true,
		},
		{
			name: "unknown keyword",
			input: []byte(`
				{
					"password": "_go:password",
				}
			`),
			wantErr: true,
		},
		{
			name: "invalid reference",
			input: []byte(`
				{
					"confirm_password": "_go:_password",
				}
			`),
			wantErr: true,
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testParser.parse(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_parser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]ParsedField
		wantErr bool
	}{
		{
			name: "positive #1",
			input: []byte(`
				{
					"email": "_go:<4/1>@test.com",
					"password": "_go:<5/2/1>",
					"confirm_password": "_go:_password",
					"escaped_field": "_go:\\<0\\> <5>"
				}
			`),
			want: map[string]ParsedField{
				"confirm_password": &referenceField{
					name:        "confirm_password",
					referenceTo: "password",
				},
				"email": &patternField{
					name: "email",
					pattern: &generator{
						strPattern{
							suffix:  "@test.com",
							letters: 4,
							digits:  1,
							length:  5,
						},
					},
				},
				"escaped_field": &patternField{
					name: "escaped_field",
					pattern: &generator{
						strPattern{
							prefix:  "<0> ",
							letters: 5,
							length:  5,
						},
					},
				},
				"password": &patternField{
					name: "password",
					pattern: &generator{
						strPattern{
							letters:  5,
							digits:   2,
							specials: 1,
							length:   8,
						},
					},
				},
			},
		},
		{
			name: "invalid pattern",
			input: []byte(`
				{
					"email": "_go:<4/1@test.com",
					"password": "_go:<5/2/1>",
					"confirm_password": "_go:_password"
				}
			`),
			wantErr: true,
		},
		{
			name: "unescaped suffix",
			input: []byte(`
				{
					"email": "_go:<4/1><@test.com>",
				}
			`),
			wantErr: true,
		},
		{
			name: "pattern is too long",
			input: []byte(`
				{
					"email": "_go:<67>@test.com",
				}
			`),
			wantErr: true,
		},
		{
			name: "unknown keyword",
			input: []byte(`
				{
					"password": "_go:password",
				}
			`),
			wantErr: true,
		},
		{
			name: "invalid reference",
			input: []byte(`
				{
					"confirm_password": "_go:_password",
				}
			`),
			wantErr: true,
		},
	}

	testParser := getTestParserStruct()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testParser.Parse(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

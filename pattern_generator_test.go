package goson

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_locatePattern(t *testing.T) {
	tests := []struct {
		name  string
		input string
		start int
		end   int
	}{
		{
			name:  "positive pattern only",
			input: "<4/2/1>",
			start: 0,
			end:   6,
		},
		{
			name:  "positive with prefix and suffix",
			input: "prefix <4/2/1> suffix",
			start: 7,
			end:   13,
		},
		{
			name:  "positive char escape",
			input: "\\<prefix\\><4/2/1>",
			start: 10,
			end:   16,
		},
		{
			name:  "invalid #1",
			input: "4/2/1>",
			start: -1,
			end:   5,
		},
		{
			name:  "invalid #2",
			input: "<4/2/1",
			start: 0,
			end:   -1,
		},
		{
			name:  "invalid #3",
			input: "4/2/1",
			start: -1,
			end:   -1,
		},
		{
			name:  "invalid #4",
			input: "<4/<2/1>",
			start: -1,
			end:   7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := locatePattern(tt.input)
			assert.Equal(t, tt.start, gotStart)
			assert.Equal(t, tt.end, gotEnd)
		})
	}
}

func Test_unescapeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "positive #1",
			input: "plain text",
			want:  "plain text",
		},
		{
			name:  "positive #2",
			input: "\\<field\\>",
			want:  "<field>",
		},
		{
			name:  "positive #3",
			input: "\\<\\<\\>\\>",
			want:  "<<>>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := unescapeString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parsePattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    strPattern
		wantErr bool
	}{
		{
			name:  "positive full pattern",
			input: "4/2/1",
			want: strPattern{
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			wantErr: false,
		},
		{
			name:  "positive without digits",
			input: "4/2",
			want: strPattern{
				letters: 4,
				digits:  2,
				length:  6,
			},
			wantErr: false,
		},
		{
			name:  "positive only letters",
			input: "4",
			want: strPattern{
				letters: 4,
				length:  4,
			},
			wantErr: false,
		},
		{
			name:    "invalid pattern len",
			input:   "4/2/1/",
			wantErr: true,
		},
		{
			name:    "empty value",
			input:   "4//1",
			wantErr: true,
		},
		{
			name:    "non numeric value #1",
			input:   "4/2/c",
			wantErr: true,
		},
		{
			name:    "non numeric value #2",
			input:   "a",
			wantErr: true,
		},
		{
			name:    "pattern is to long",
			input:   "32/32/1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePattern(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_buildPattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    strPattern
		wantErr bool
	}{
		{
			name:  "positive body only",
			input: "<4/2/1>",
			want: strPattern{
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
		},
		{
			name:  "positive with prefix",
			input: "some prefix <4/2/1>",
			want: strPattern{
				prefix:   "some prefix ",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
		},
		{
			name:  "positive with suffix",
			input: "<4/2/1> some suffix",
			want: strPattern{
				suffix:   " some suffix",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
		},
		{
			name:  "positive with prefix and suffix",
			input: "some prefix <4/2/1> some suffix",
			want: strPattern{
				prefix:   "some prefix ",
				suffix:   " some suffix",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
		},
		{
			name:  "positive with escapes",
			input: "I \\<3 U <4/2/1> 1 \\> 0",
			want: strPattern{
				prefix:   "I <3 U ",
				suffix:   " 1 > 0",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
		},
		{
			name:    "invalid pattern #1",
			input:   "<4/2/1",
			wantErr: true,
		},
		{
			name:    "invalid pattern #2",
			input:   "<4/2/>",
			wantErr: true,
		},
		{
			name:    "invalid escape",
			input:   "<some prefix> <4/2/1>",
			wantErr: true,
		},
		{
			name:    "invalid pattern is too long",
			input:   "<32/32/1>",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildPattern(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_getRandomCharFromSource(t *testing.T) {
	seedTestDate()

	source := "0123456789"
	target := "0261344689" // pre-generated string with a test date
	for i := 0; i < len(target); i++ {
		t.Run(fmt.Sprintf("positive %d", i+1), func(t *testing.T) {
			assert.Equal(t, target[i], getRandomCharFromSource(source))
		})
	}
}

func seedTestDate() {
	date, _ := time.Parse("2006-01-02", "1975-02-24")
	rand.Seed(date.Unix())
}

func Test_generator_generate(t *testing.T) {
	tests := []struct {
		name    string
		pattern strPattern
		want    string
	}{
		{
			name: "positive #1",
			pattern: strPattern{
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "wQET3~4",
		},
		{
			name: "positive #2",
			pattern: strPattern{
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "f)j0vx2",
		},
		{
			name: "positive #3",
			pattern: strPattern{
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "+UGyD83",
		},
	}

	seedTestDate()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testGenerator := &generator{
				pattern: tt.pattern,
			}
			got := testGenerator.generate()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_generator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		pattern strPattern
		want    string
	}{
		{
			name: "body only",
			pattern: strPattern{
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "wQET3~4",
		},
		{
			name: "with prefix",
			pattern: strPattern{
				prefix:   "some prefix ",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "some prefix f)j0vx2",
		},
		{
			name: "with suffix",
			pattern: strPattern{
				suffix:   " some suffix",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "+UGyD83 some suffix",
		},
		{
			name: "with prefix and suffix",
			pattern: strPattern{
				prefix:   "some prefix ",
				suffix:   " some suffix",
				letters:  4,
				digits:   2,
				specials: 1,
				length:   7,
			},
			want: "some prefix b$Tl7G1 some suffix",
		},
	}

	seedTestDate()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testGenerator := &generator{
				pattern: tt.pattern,
			}

			got := testGenerator.Generate()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_newPatternGenerator(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "positive body only #1",
			input: "<4/2/1>",
		},
		{
			name:  "positive body only #2",
			input: "<4/2>",
		},
		{
			name:  "positive body only #3",
			input: "<4>",
		},
		{
			name:  "positive with prefix and suffix",
			input: "some prefix <4/2/1> some suffix",
		},
		{
			name:  "positive with escaped suffix",
			input: "\\<prefix\\> <4/2/1>",
		},
		{
			name:    "invalid pattern #1",
			input:   "<4/2/1",
			wantErr: true,
		},
		{
			name:    "invalid pattern #2",
			input:   "<4/2/>",
			wantErr: true,
		},
		{
			name:    "invalid unescaped prefix",
			input:   "<prefix><4/2/1>",
			wantErr: true,
		},
		{
			name:    "invalid non numeric value #1",
			input:   "<4/2/c>",
			wantErr: true,
		},
		{
			name:    "invalid non numeric value #2",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid pattern is too long",
			input:   "<34/32/1>",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newPatternGenerator(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

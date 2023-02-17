package framework

import (
	"reflect"
	"testing"
)

func TestClean_Success(t *testing.T) {

	testCases := []struct {
		name string
		in   struct {
			description        SchemaDescription
			removeTrailingStop bool
		}
		expected SchemaDescription
	}{
		{
			name: "passthrough-no-changes",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "test123 this is a description - not formatted",
					MarkdownDescription: "test123 this is a description - md",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted",
				MarkdownDescription: "test123 this is a description - md",
			},
		},
		{
			name: "blank-description",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "",
					MarkdownDescription: "test123 this is a description - md",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - md",
				MarkdownDescription: "test123 this is a description - md",
			},
		},
		{
			name: "blank-md-description",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "test123 this is a description - not formatted",
					MarkdownDescription: "",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted",
				MarkdownDescription: "test123 this is a description - not formatted",
			},
		},
		{
			name: "nil-description",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					MarkdownDescription: "test123 this is a description - md",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - md",
				MarkdownDescription: "test123 this is a description - md",
			},
		},
		{
			name: "nil-md-description",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description: "test123 this is a description - not formatted",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted",
				MarkdownDescription: "test123 this is a description - not formatted",
			},
		},
		{
			name: "trailing-spaces",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "       test123 this is a description - not formatted         ",
					MarkdownDescription: "       test123 this is a description - md         ",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted",
				MarkdownDescription: "test123 this is a description - md",
			},
		},
		{
			name: "trailing-stop-don't-remove",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "test123 this is a description - not formatted.",
					MarkdownDescription: "test123 this is a description - md.",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted.",
				MarkdownDescription: "test123 this is a description - md.",
			},
		},
		{
			name: "trailing-stop-don't-remove-with-spaces",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "test123 this is a description - not formatted.   ",
					MarkdownDescription: "test123 this is a description - md.  ",
				},
				removeTrailingStop: false,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted.",
				MarkdownDescription: "test123 this is a description - md.",
			},
		},
		{
			name: "trailing-stop-remove",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "test123 this is a description - not formatted.",
					MarkdownDescription: "test123 this is a description - md.",
				},
				removeTrailingStop: true,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted",
				MarkdownDescription: "test123 this is a description - md",
			},
		},
		{
			name: "trailing-stop-remove-with-spaces",
			in: struct {
				description        SchemaDescription
				removeTrailingStop bool
			}{
				description: SchemaDescription{
					Description:         "test123 this is a description - not formatted  .  ",
					MarkdownDescription: "test123 this is a description - md  .  ",
				},
				removeTrailingStop: true,
			},
			expected: SchemaDescription{
				Description:         "test123 this is a description - not formatted",
				MarkdownDescription: "test123 this is a description - md",
			},
		},
	}

	for _, test := range testCases {
		got := test.in.description
		got.Clean(test.in.removeTrailingStop)
		if !reflect.DeepEqual(test.expected, got) {
			t.Fatalf("\nTest: \t\t%s\nExpected: \t%s\ngot:\t\t%s", test.name, test.expected, got)
		}
	}

}

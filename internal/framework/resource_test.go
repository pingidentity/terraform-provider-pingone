// Copyright © 2026 Ping Identity Corporation

package framework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringToTF_Success(t *testing.T) {

	testCases := []struct {
		name     string
		in       string
		expected basetypes.StringValue
	}{
		{
			name:     "success-with-value",
			in:       "testvalue",
			expected: types.StringValue("testvalue"),
		},
		{
			name:     "success-with-blank-string",
			in:       "",
			expected: types.StringNull(),
		},
	}

	for _, test := range testCases {
		if got := StringToTF(test.in); !test.expected.Equal(got) {
			t.Fatalf("\nTest: \t\t%s\nExpected: \t%s\ngot:\t\t%s", test.name, test.expected, got)
		}
	}

}

func TestStringSliceToTF_Success(t *testing.T) {

	testCases := []struct {
		name     string
		in       []string
		expected basetypes.ListValue
	}{
		{
			name: "success-with-values",
			in:   []string{"value1", "value2", "value3"},
			expected: func() basetypes.ListValue {

				v := []string{"value1", "value2", "value3"}

				list := make([]attr.Value, 0)
				for _, item := range v {
					list = append(list, StringToTF(item))
				}

				values, _ := types.ListValue(types.StringType, list)
				return values
			}(),
		},
		{
			name: "success-with-no-values",
			in:   []string{},
			expected: func() basetypes.ListValue {
				values, _ := types.ListValue(types.StringType, make([]attr.Value, 0))
				return values
			}(),
		},
		{
			name:     "success-with-nil",
			in:       nil,
			expected: types.ListNull(types.StringType),
		},
	}

	for _, test := range testCases {
		if got, _ := StringSliceToTF(test.in); !test.expected.Equal(got) {
			t.Fatalf("\nTest: \t\t%s\nExpected: \t%s\ngot:\t\t%s", test.name, test.expected, got)
		}
	}

}

// func TestTFListToStringSlice_Success(t *testing.T) {

// 	testCases := []struct {
// 		name     string
// 		in       basetypes.ListValue
// 		expected []*string
// 	}{
// 		{
// 			name: "success-with-values",
// 			in: func() basetypes.ListValue {

// 				v := []string{"value1", "value2", "value3"}

// 				list := make([]attr.Value, 0)
// 				for _, item := range v {
// 					list = append(list, StringToTF(item))
// 				}

// 				values, _ := types.ListValue(types.StringType, list)
// 				return values
// 			}(),
// 			expected: func() []*string {

// 				v := []string{"value1", "value2", "value3"}

// 				list := make([]*string, 0)
// 				for _, item := range v {
// 					list = append(list, &item)
// 				}

// 				return list
// 			}(),
// 		},
// 		{
// 			name: "success-with-no-values",
// 			in: func() basetypes.ListValue {
// 				values, _ := types.ListValue(types.StringType, make([]attr.Value, 0))
// 				return values
// 			}(),
// 			expected: []*string{},
// 		},
// 		{
// 			name:     "success-with-nil",
// 			in:       types.ListNull(types.StringType),
// 			expected: nil,
// 		},
// 	}

// 	for _, test := range testCases {
// 		ctx := context.TODO()

// 		if got := TFListToStringSlice(ctx, test.in); !reflect.DeepEqual(test.expected, got) {
// 			t.Fatalf("\nTest: \t\t%s\nExpected: \t%+v\ngot:\t\t%+v", test.name, test.expected, got)
// 		}
// 	}

// }

package configo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testcases := []struct {
		description    string
		filepaths      []string
		expectedKeys   []string
		expectedValues []string
		err            error
	}{
		{
			description: "sample",
			filepaths:   []string{"./fixtures/sample.yaml"},
			expectedKeys: []string{
				"string",
				"bool",
				"int",
				"int8",
				"int16",
				"int32",
				"int64",
				"uint",
				"uint8",
				"uint16",
				"uint32",
				"uint64",
				"float32",
				"float64",
				"map.string",
				"map.bool",
				"map.int",
				"map.int8",
				"map.int16",
				"map.int32",
				"map.int64",
				"map.uint",
				"map.uint8",
				"map.uint16",
				"map.uint32",
				"map.uint64",
				"map.float32",
				"map.float64",
				"slice.0",
				"slice.1",
				"slice.2",
				"slice.3",
				"slice.4",
				"slice.5",
				"slice.6",
				"slice.7",
				"slice.8",
				"slice.9",
				"slice.10",
				"slice.11",
				"slice.12",
				"slice.13",
				"slice.map.0.string",
				"slice.map.0.bool",
				"slice.map.0.int",
				"slice.map.0.int8",
				"slice.map.0.int16",
				"slice.map.0.int32",
				"slice.map.0.int64",
				"slice.map.0.uint",
				"slice.map.0.uint8",
				"slice.map.0.uint16",
				"slice.map.0.uint32",
				"slice.map.0.uint64",
				"slice.map.0.float32",
				"slice.map.0.float64",
				"slice.map.1.string",
				"slice.map.1.bool",
				"slice.map.1.int",
				"slice.map.1.int8",
				"slice.map.1.int16",
				"slice.map.1.int32",
				"slice.map.1.int64",
				"slice.map.1.uint",
				"slice.map.1.uint8",
				"slice.map.1.uint16",
				"slice.map.1.uint32",
				"slice.map.1.uint64",
				"slice.map.1.float32",
				"slice.map.1.float64",
			},
			expectedValues: []string{
				"test",
				"true",
				"9223372036854775807",
				"127",
				"32767",
				"2147483647",
				"9223372036854775807",
				"18446744073709551615",
				"255",
				"65535",
				"4294967295",
				"18446744073709551615",
				"65535.00390625",
				"65535.00390625",
				"test",
				"true",
				"9223372036854775807",
				"127",
				"32767",
				"2147483647",
				"9223372036854775807",
				"18446744073709551615",
				"255",
				"65535",
				"4294967295",
				"18446744073709551615",
				"65535.00390625",
				"65535.00390625",
				"test",
				"true",
				"9223372036854775807",
				"127",
				"32767",
				"2147483647",
				"9223372036854775807",
				"18446744073709551615",
				"255",
				"65535",
				"4294967295",
				"18446744073709551615",
				"65535.00390625",
				"65535.00390625",
				"testA",
				"true",
				"9223372036854775807",
				"127",
				"32767",
				"2147483647",
				"9223372036854775807",
				"18446744073709551615",
				"255",
				"65535",
				"4294967295",
				"18446744073709551615",
				"65535.00390625",
				"65535.00390625",
				"testB",
				"false",
				"-9223372036854775808",
				"-128",
				"-32768",
				"-2147483648",
				"-9223372036854775808",
				"-18446744073709551616",
				"-256",
				"-65536",
				"-4294967296",
				"-18446744073709551616",
				"-65535.00390625",
				"-65535.00390625",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.description, func(t *testing.T) {
			result, err := ParseFile(testcase.filepaths)

			if testcase.err == nil {
				require.NoError(t, err)
				assert.Equal(t, testcase.expectedKeys, result.Keys())
				assert.Equal(t, testcase.expectedValues, result.Values())
			} else {
				require.NotNil(t, err)
				require.EqualError(t, testcase.err, err.Error())
			}
		})
	}
}

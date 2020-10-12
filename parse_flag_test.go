package configo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlag(t *testing.T) {
	testcases := []struct {
		description    string
		args           []string
		expectedKeys   []string
		expectedValues []string
		err            error
	}{
		{
			description:    "simple",
			args:           []string{"--foo", "bar", "--baz", "cat"},
			expectedKeys:   []string{"foo", "baz"},
			expectedValues: []string{"bar", "cat"},
		},
		{
			description:    "boolean in beginning",
			args:           []string{"--boolean", "--foo", "bar", "--baz", "cat"},
			expectedKeys:   []string{"boolean", "foo", "baz"},
			expectedValues: []string{"true", "bar", "cat"},
		},
		{
			description:    "boolean in middle",
			args:           []string{"--foo", "bar", "--boolean", "--baz", "cat"},
			expectedKeys:   []string{"foo", "boolean", "baz"},
			expectedValues: []string{"bar", "true", "cat"},
		},
		{
			description:    "boolean in end",
			args:           []string{"--foo", "bar", "--baz", "cat", "--boolean"},
			expectedKeys:   []string{"foo", "baz", "boolean"},
			expectedValues: []string{"bar", "cat", "true"},
		},
		{
			description:    "complicated keys",
			args:           []string{"--foo-slice-0-bar", "baz", "--foo-map-key", "value"},
			expectedKeys:   []string{"foo.slice.0.bar", "foo.map.key"},
			expectedValues: []string{"baz", "value"},
		},
		{
			description:    "repeated keys",
			args:           []string{"--foo", "baz", "--foo", "bar"},
			expectedKeys:   []string{"foo.0", "foo.1"},
			expectedValues: []string{"baz", "bar"},
		},
		{
			description:    "terminator",
			args:           []string{"--foo", "bar", "--", "baz"},
			expectedKeys:   []string{"foo"},
			expectedValues: []string{"bar"},
		},
		{
			description:    "empty",
			args:           []string{"foo", "bar", "baz"},
			expectedValues: []string{},
		},
		{
			description: "bad syntax",
			args:        []string{"---foo", "baz"},
			err:         errors.New("bad flag syntax: ---foo"),
		},
		{
			description: "bad syntax",
			args:        []string{"--=foo", "baz"},
			err:         errors.New("bad flag syntax: --=foo"),
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.description, func(t *testing.T) {
			result, err := ParseFlag(testcase.args)

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

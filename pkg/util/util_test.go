package util_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"go-rest-api-boilerplate/pkg/util"
)

func TestStringToInt64(t *testing.T) {
	cases := []struct {
		str      string
		expected int64
		error    error
	}{
		{str: "10", expected: 10, error: nil},
		{str: "10x", expected: 0, error: strconv.ErrSyntax},
	}

	for _, s := range cases {
		n, err := util.StringToInt64(s.str)
		assert.Equal(t, s.expected, n)
		assert.ErrorIs(t, err, s.error)
	}
}

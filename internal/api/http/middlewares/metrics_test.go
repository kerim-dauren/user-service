package middlewares

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHandlerName(t *testing.T) {
	tests := []struct {
		fullName string
		expected string
	}{
		{
			fullName: "github.com/kerim-dauren/user-service/api/http/v1/routes.(*UserHandler).createUser",
			expected: "UserHandler",
		},
		{
			fullName: "unexpectedFormatHandler",
			expected: "unexpectedFormatHandler",
		},
		{
			fullName: "some/other/package.(*Struct).Method",
			expected: "Struct",
		},
	}

	for _, tt := range tests {
		actual := parseHandlerName(tt.fullName)
		assert.Equal(t, tt.expected, actual)
	}
}

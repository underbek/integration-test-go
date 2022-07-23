package step_9_api_fixtures

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func JSONEq(t *testing.T, expected, actual any) bool {
	return assert.JSONEq(t, jsonMarshal(t, expected), jsonMarshal(t, actual))
}

func jsonMarshal(t *testing.T, data any) string {
	switch v := data.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case io.Reader:
		data, err := io.ReadAll(v)
		require.NoError(t, err)
		return string(data)
	default:
		res, err := json.Marshal(v)
		require.NoError(t, err)
		return string(res)
	}
}

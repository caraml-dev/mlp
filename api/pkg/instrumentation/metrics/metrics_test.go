package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusString(t *testing.T) {
	assert.Equal(t, "success", GetStatusString(true))
	assert.Equal(t, "failure", GetStatusString(false))
}

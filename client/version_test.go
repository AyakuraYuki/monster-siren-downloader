package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	buildVersion = "test-001"
	assert.EqualValues(t, "test-001", Version())
}

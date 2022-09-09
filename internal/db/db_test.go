package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	pg := DB(&postgre{})
	assert.NotNil(t, pg)

	unknown := DB(nil)
	assert.Nil(t, unknown)
}

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgreeDb(t *testing.T) {
	pg := NewPostgreeDb("", "", "", "", "")
	assert.NotNil(t, pg)
}

func TestPostgre_DSN(t *testing.T) {
	dsn := NewPostgreeDb("localhost", "3352", "db_name", "postgres", "postgres").DSN()
	t.Log(dsn)
	assert.Equal(t, "host=localhost port=3352 user=postgres dbname=db_name password=postgres sslmode=disable", dsn)
}

func TestPostgre_Connect(t *testing.T) {
	//mock
}

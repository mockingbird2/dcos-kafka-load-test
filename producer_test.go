package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeTickerInterval(t *testing.T) {
	assert.Equal(t, computeTickerInterval(1, 1), int64(1e9))
}

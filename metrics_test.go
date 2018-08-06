package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddBatch(t *testing.T) {
	metrics := CreateMetrics()
	metrics.AddBatch()
	assert.Equal(t, metrics.SentBatches(), uint64(1))
}

func TestAddError(t *testing.T) {
	metrics := CreateMetrics()
	metrics.AddError()
	assert.Equal(t, metrics.Errors(), uint64(1))
}

func TestStartReporting(t *testing.T) {
	metrics := CreateMetrics()
	metrics.StartReporting()
	metrics.AddBatch()
	metrics.StopReporting()
	assert.True(t, metrics.SentBatches() == 1)
}

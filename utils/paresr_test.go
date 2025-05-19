package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseCustomReportDates(t *testing.T) {
	input := "2025-05-01 - 2025-05-15"
	outputFrom := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)
	outputTo := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)
	from, to, err := ParseCustomReportDates(input)

	assert.NoError(t, err, "Unexpected error while parsing custom report")
	assert.Equal(t, outputFrom, from)
	assert.Equal(t, outputTo, to)
}

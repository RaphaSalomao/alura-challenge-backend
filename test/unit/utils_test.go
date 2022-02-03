package test_test

import (
	"testing"
	"time"

	"github.com/RaphaSalomao/alura-challenge-backend/utils"
	"github.com/stretchr/testify/assert"
)

func TestMonthInterval(t *testing.T) {
	t1, t2, err := utils.MonthInterval("2022-01")
	assert.NoError(t, err)
	assert.Equal(t, t1, time.Date(2022, time.January, 1, 0, 0, 0, 0, time.Local))
	assert.Equal(t, t2, time.Date(2022, time.February, 1, 0, 0, 0, -1, time.Local))
}

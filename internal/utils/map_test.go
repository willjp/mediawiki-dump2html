package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	integers := []int{1, 2, 3, 4, 5}
	expects := []string{"1", "2", "3", "4", "5"}
	result := Map(integers, func(val int) string {
		return strconv.Itoa(val)
	})
	assert.Equal(t, expects, result)
}

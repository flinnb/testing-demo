package extras

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtras(t *testing.T) {
	t.Run("zero allocation swap", func(t *testing.T) {
		x := 23
		y := 17
		x, y = y, x
		assert.Equal(t, 17, x)
		assert.Equal(t, 23, y)
	})
	t.Run("interface games", func(t *testing.T) {
		var foo interface{}
		foo = 7
		bar := foo.(int)
		var expectedType int
		assert.Equal(t, 7, bar)
		assert.IsType(t, expectedType, bar)
	})
}

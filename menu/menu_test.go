package menu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_randChar(t *testing.T) {
	m := NewMenu(nil)
	iterations := 20
	length := []int{
		1,
		5,
		7,
		8,
		16,
		20,
		33,
	}

	for i := 0; i < iterations; i++ {
		for _, v := range length {
			rndString := m.randChar(v)
			assert.Len(t, rndString, v)
		}
	}
}

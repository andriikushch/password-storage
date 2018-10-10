package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_randChar(t *testing.T) {
	iterations := 20
	lenght := []int{
		1,
		5,
		7,
		8,
		16,
		20,
		33,
	}

	for i := 0; i < iterations; i++ {
		for _, v := range lenght {
			rndString := randChar(v)
			assert.Len(t, rndString, v)
		}
	}
}

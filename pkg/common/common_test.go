package common

import (
	"testing"

	. "github.com/franela/goblin"
)

func TestAdd(t *testing.T) {
	gob := Goblin(t)
	gob.Describe("Common Add Function", func() {
		gob.It("should add two numbers ", func() {
			gob.Assert(Add(1, 2)).Equal(3)
			gob.Assert(Add(1, 0)).Equal(1)
			gob.Assert(Add(2, -2)).Equal(0)
		})
	})
}

package escape

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {

	ts := []struct {
		before string

		after string
	}{
		{
			`HELLO "WORLD"`,
			`HELLO \"WORLD\"`,
		},
		{
			`HELLO \tWORLD`,
			`HELLO \\tWORLD`,
		},
	}

	for _, tc := range ts {

		got := String(tc.before)

		want := tc.after
		assert.Equal(t, want, got)
	}

}

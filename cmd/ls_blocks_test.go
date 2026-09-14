package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintBlocks(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		width  int
		blocks []block
		want   string
	}{
		"wraps the description under the width": {
			width: 40,
			blocks: []block{{
				name:        "go-review",
				labels:      "go, review",
				description: "Reviews Go code for correctness and idiom and reports every finding.",
			}},
			want: "  go-review   go, review\n" +
				"      Reviews Go code for correctness\n" +
				"      and idiom and reports every\n" +
				"      finding.\n",
		},
		"one blank line between blocks and none after the last": {
			width: 80,
			blocks: []block{
				{name: "go-review", labels: "go, review", description: "Reviews Go code."},
				{name: "sql-review", labels: "sql", description: "Reviews SQL queries."},
			},
			want: "  go-review   go, review\n" +
				"      Reviews Go code.\n" +
				"\n" +
				"  sql-review   sql\n" +
				"      Reviews SQL queries.\n",
		},
		"empty labels leave the name alone on the header": {
			width:  80,
			blocks: []block{{name: "wip", description: "Not ready yet."}},
			want:   "  wip\n      Not ready yet.\n",
		},
		"a word longer than the width stays whole on its own line": {
			width:  20,
			blocks: []block{{name: "x", description: "a verylongwordthatdoesnotfit b"}},
			want: "  x\n" +
				"      a\n" +
				"      verylongwordthatdoesnotfit\n" +
				"      b\n",
		},
		"wrapping counts runes, not bytes": {
			width:  16,
			blocks: []block{{name: "x", description: "éé éé éé éé"}},
			want: "  x\n" +
				"      éé éé éé\n" +
				"      éé\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer

			err := printBlocks(&out, tt.width, tt.blocks)

			require.NoError(t, err)
			assert.Equal(t, tt.want, out.String())
		})
	}
}

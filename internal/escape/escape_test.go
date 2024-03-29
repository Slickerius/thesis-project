// Copyright 2021 The Mellium Contributors.
// Use of this source code is governed by the BSD 2-clause
// license that can be found in the LICENSE file.

package escape_test

import (
	"strconv"
	"testing"

	"github.com/rivo/tview"
	"golang.org/x/text/transform"

	"mellium.im/communique/internal/escape"
)

var escapeTests = [...]struct {
	in, out string
}{
	0: {},
	1: {in: `["abc"][""][][red]`, out: `["abc"[][""[][][red[]`},
	2: {in: `[""[[[]`, out: `[""[[[[]`},
	3: {in: `["a[bc"]`, out: `["a[bc"[]`},
	4: {in: `["a]bc"]`, out: `["a[]bc"]`},
}

func TestEscape(t *testing.T) {
	for i, tc := range escapeTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			et := escape.Transformer()
			out, _, err := transform.String(et, tc.in)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if out != tc.out {
				t.Errorf("want=`%s`, got=`%s`", tc.out, out)
			}
		})
	}
}

const benchEscape = `["abc"][""][][red][""[[[]["a[bc"]["a]bc"]`

func BenchmarkEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = tview.Escape(benchEscape)
	}
}

func BenchmarkTransform(b *testing.B) {
	t := escape.Transformer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = transform.String(t, benchEscape)
	}
}

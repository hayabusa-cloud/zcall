// Copyright 2025 Hayabusa Cloud Co., Ltd. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package zcall

import "testing"

// TestUitoa tests the internal uitoa function directly.
func TestUitoa(t *testing.T) {
	tests := []struct {
		input uint
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{10, "10"},
		{42, "42"},
		{100, "100"},
		{999, "999"},
		{1234567890, "1234567890"},
		{18446744073709551615, "18446744073709551615"}, // max uint64
	}
	for _, tt := range tests {
		got := uitoa(tt.input)
		if got != tt.want {
			t.Errorf("uitoa(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

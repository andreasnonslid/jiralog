package main

import "testing"

func TestSplitComma(t *testing.T) {
	in := "a,b,c"
	out := splitComma(in)
	if len(out) != 3 || out[0] != "a" || out[1] != "b" || out[2] != "c" {
		t.Fatalf("unexpected output: %v", out)
	}
}

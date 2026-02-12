package main

import "testing"

func TestPrepare(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
	}{
		{
			name:  "empty statement",
			input: "",
		},
		{
			name:  "invalid attribute type",
			input: "select one",
		},
		{
			name:  "invalid statement type",
			input: "meme 1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stmt Statement
			_, err := stmt.Prepare(tt.input)
			if err == nil {
				t.Errorf("test: %s failed. input %s should have been validated", tt.name, tt.input)
			}
		})
	}
}

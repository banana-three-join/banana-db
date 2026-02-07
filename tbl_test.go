package main

import "testing"

func TestGetRowsByPage(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input int
		tbl   Table
	}{
		{"fail on out of bounds page request [negative value]", -1, Table{}},
		{"fail on out of bounds page request [positive value]", 101, Table{}},
		{"fail on request of empty page", 30, Table{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.tbl.GetRowsByPage(tt.input); err == nil {
				t.Errorf("test: %s failed", tt.name)
			}
		})
	}

}

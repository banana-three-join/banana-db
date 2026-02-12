package main

import "testing"

func TestDoCommand(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
	}{
		{".help input", ".help"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := doMetaCommand(tt.input); err != nil {
				t.Fatalf("test: %s failed. the %s isn't recognized", tt.name, tt.input)
			}
		})
	}

	t.Run("fail on unknown cmd", func(t *testing.T) {
		if err := doMetaCommand("nonExistantCmd"); err == nil {
			t.Fatalf("test: %s failed. meta command shouldn't allow unknown commands", t.Name())
		}
	})
}

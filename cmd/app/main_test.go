package main

import (
	"strings"
	"testing"

	"atad-project/internal/version"
)

func TestVersionReportsContractFields(t *testing.T) {
	var b strings.Builder
	if err := version.WriteJSON(&b); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	for _, field := range []string{`"app"`, `"version"`, `"commit"`, `"built_at"`} {
		if !strings.Contains(b.String(), field) {
			t.Errorf("version output is missing %s: %s", field, b.String())
		}
	}
}

func TestYourLogic(t *testing.T) {
	t.Skip("write your own tests here")
}

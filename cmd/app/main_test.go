package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"atad-project/internal/version"
)

func TestContractRoutes(t *testing.T) {
	mux := newMux()
	do := func(method, path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
		return rec
	}

	rec := do("GET", "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("/health status = %d, want 200", rec.Code)
	}
	var h struct {
		Status string `json:"status"`
		Uptime *int   `json:"uptime_seconds"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &h); err != nil || h.Status != "ok" || h.Uptime == nil {
		t.Errorf("/health body = %s", rec.Body)
	}

	rec = do("GET", "/version")
	if rec.Code != http.StatusOK {
		t.Fatalf("/version status = %d, want 200", rec.Code)
	}
	var v map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("/version body: %v", err)
	}
	for _, k := range []string{"app", "version", "commit", "built_at"} {
		if _, ok := v[k]; !ok {
			t.Errorf("/version is missing %q", k)
		}
	}

	rec = do("GET", "/")
	if rec.Code != http.StatusOK || !strings.HasPrefix(rec.Header().Get("Content-Type"), "text/html") {
		t.Errorf("/ = %d %q, want 200 text/html", rec.Code, rec.Header().Get("Content-Type"))
	}
	for _, want := range []string{"Personal Finance Tracker", "Nistor Carina", "DMWD", "Project 2"} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("/ does not mention %q", want)
		}
	}

	rec = do("POST", "/reset")
	if rec.Code != http.StatusNoContent || rec.Body.Len() != 0 {
		t.Errorf("POST /reset = %d with body %q, want 204 and no body", rec.Code, rec.Body)
	}
}

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

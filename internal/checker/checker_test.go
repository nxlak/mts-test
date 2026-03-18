package checker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nxlak/mts-test/internal/models"
)

func newTestServer(t *testing.T, responses map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
		modulePath := path[:len(path)-len("/@latest")]

		version, ok := responses[modulePath]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(latestResponse{Version: version})
	}))
}

func TestChecker_FindUpdates(t *testing.T) {
	srv := newTestServer(t, map[string]string{
		"github.com/stretchr/testify": "v1.9.0", // newer than v1.8.0
		"github.com/gomodule/redigo":  "v1.9.3", // same as pinned
	})
	defer srv.Close()

	c := NewChecker(
		WithProxyBase(srv.URL),
		WithWorkers(2),
		WithHTTPClient(srv.Client()),
	)

	deps := []models.Dependency{
		{Path: "github.com/stretchr/testify", CurrentVersion: "v1.8.0"},
		{Path: "github.com/gomodule/redigo", CurrentVersion: "v1.9.3"},
		{Path: "github.com/golang-jwt/jwt", CurrentVersion: "v4.0.0"}, // 404 → skipped
	}

	updates, err := c.FindUpdates(context.Background(), deps)
	if err != nil {
		t.Fatalf("FindUpdates: %v", err)
	}

	if len(updates) != 1 {
		t.Fatalf("len(updates) = %d; want 1", len(updates))
	}
	if updates[0].Path != "github.com/stretchr/testify" {
		t.Errorf("updates[0].Path = %q; want github.com/stretchr/testify", updates[0].Path)
	}
	if updates[0].LatestVersion != "v1.9.0" {
		t.Errorf("updates[0].LatestVersion = %q; want v1.9.0", updates[0].LatestVersion)
	}
}

func TestChecker_FindUpdates_EmptyList(t *testing.T) {
	srv := newTestServer(t, nil)
	defer srv.Close()

	c := NewChecker(WithProxyBase(srv.URL), WithHTTPClient(srv.Client()))
	updates, err := c.FindUpdates(context.Background(), nil)
	if err != nil {
		t.Fatalf("FindUpdates: %v", err)
	}
	if len(updates) != 0 {
		t.Errorf("expected empty updates, got %d", len(updates))
	}
}

func Test_isNewer(t *testing.T) {
	tests := []struct {
		latest, current string
		want            bool
	}{
		{"v1.1.0", "v1.0.0", true},
		{"v1.0.0", "v1.0.0", false},
		{"v0.9.0", "v1.0.0", false},
		{"not-semver", "v1.0.0", false},
		{"v1.0.0", "not-semver", false},
	}

	for _, tt := range tests {
		got := isNewer(tt.latest, tt.current)
		if got != tt.want {
			t.Errorf("isNewer(%q, %q) = %v; want %v", tt.latest, tt.current, got, tt.want)
		}
	}
}

func Test_encodePath(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"github.com/stretchr/testify", "github.com/stretchr/testify"},
		{"github.com/sirupsen/logrus", "github.com/sirupsen/logrus"},
		{"github.com/golang-jwt/jwt", "github.com/golang-jwt/jwt"},
	}

	for _, tt := range tests {
		got := encodePath(tt.input)
		if got != tt.want {
			t.Errorf("encodePath(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

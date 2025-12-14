package kv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "input.env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestParseFile(t *testing.T) {
	tests := map[string]struct {
		content       string
		want          map[string]string
		expectErrText string
	}{
		"valid file with comments and whitespace": {
			content: `
				# comment
				DB_HOST=localhost
				DB_USER = root

				API_KEY=   secret
			`,
			want: map[string]string{
				"DB_HOST": "localhost",
				"DB_USER": "root",
				"API_KEY": "secret",
			},
		},
		"missing separator": {
			content:       "INVALID_LINE",
			expectErrText: "missing '=' separator",
		},
		"empty key": {
			content:       "=value",
			expectErrText: "empty key",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := writeTempFile(t, tc.content)
			got, err := ParseFile(path)

			if tc.expectErrText != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.expectErrText)
				}
				if !contains(err.Error(), tc.expectErrText) {
					t.Fatalf("expected error to contain %q, got %v", tc.expectErrText, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseFile returned error: %v", err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("expected %d entries, got %d", len(tc.want), len(got))
			}

			for key, wantVal := range tc.want {
				if gotVal, ok := got[key]; !ok || gotVal != wantVal {
					t.Fatalf("key %s mismatch: want %q got %q (present=%v)", key, wantVal, got[key], ok)
				}
			}
		})
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

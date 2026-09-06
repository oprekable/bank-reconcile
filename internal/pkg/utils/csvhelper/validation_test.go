package csvhelper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDeletionPath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantErr     bool
		errContains string
	}{
		{
			name:        "path traversal with ..",
			path:        "/home/user/../../etc/passwd",
			wantErr:     true,
			errContains: "path traversal",
		},
		{
			name:        "root directory",
			path:        "/etc/passwd",
			wantErr:     true,
			errContains: "system directory",
		},
		{
			name:        "home directory too shallow",
			path:        "/home/file.txt",
			wantErr:     true,
			errContains: "system directory",
		},
		{
			name:        "tmp directory too shallow",
			path:        "/tmp/file.txt",
			wantErr:     true,
			errContains: "system directory",
		},
		{
			name:        "valid safe path",
			path:        "/tmp/data/sample/system/123.csv",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "valid nested path",
			path:        "/home/user/projects/bank-reconcile/data/file.csv",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "var log too shallow",
			path:        "/var/log/app.log",
			wantErr:     true,
			errContains: "system directory",
		},
		{
			name:        "usr bin too shallow",
			path:        "/usr/bin/evil",
			wantErr:     true,
			errContains: "system directory",
		},
		{
			name:        "relative path with traversal",
			path:        "./../../../etc/shadow",
			wantErr:     true,
			errContains: "path traversal",
		},
		{
			name:        "valid deep path",
			path:        "/var/lib/docker/volumes/data/_data/file.csv",
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "depth check - 2 levels too shallow",
			path:        "/tmp/data/file.txt",
			wantErr:     true,
			errContains: "system directory", // /tmp/data is under /tmp
		},
		{
			name:        "depth check - 3 levels OK",
			path:        "/tmp/data/sample/file.txt",
			wantErr:     false,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeletionPath(tt.path)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDeletionPath_Symlink(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create a safe directory
	safeDir := filepath.Join(tmpDir, "safe", "data", "files")
	require.NoError(t, os.MkdirAll(safeDir, 0755))

	// Create a test file in safe directory
	safeFile := filepath.Join(safeDir, "test.csv")
	require.NoError(t, os.WriteFile(safeFile, []byte("data"), 0644))

	t.Run("valid symlink within safe directory", func(t *testing.T) {
		// Create a symlink within the safe directory
		symlink := filepath.Join(safeDir, "link.csv")
		require.NoError(t, os.Symlink(safeFile, symlink))
		defer os.Remove(symlink)

		err := validateDeletionPath(symlink)
		assert.NoError(t, err)
	})

	t.Run("symlink to /etc should be blocked", func(t *testing.T) {
		// Create a symlink pointing to /etc/passwd
		dangerousLink := filepath.Join(safeDir, "dangerous")
		require.NoError(t, os.Symlink("/etc/passwd", dangerousLink))
		defer os.Remove(dangerousLink)

		err := validateDeletionPath(dangerousLink)
		require.Error(t, err)
		// On macOS, /etc resolves to /private/etc which is still blocked
		assert.Contains(t, err.Error(), "symlink resolves")
	})

	t.Run("symlink escape to shallow path", func(t *testing.T) {
		// Create a symlink pointing to /tmp (shallow)
		shallowLink := filepath.Join(safeDir, "shallow")
		require.NoError(t, os.Symlink("/tmp", shallowLink))
		defer os.Remove(shallowLink)

		err := validateDeletionPath(shallowLink)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "symlink resolves")
	})

	t.Run("symlink to direct child of dangerous root", func(t *testing.T) {
		// Create a symlink pointing to /var/log (direct child of /var)
		varLogLink := filepath.Join(safeDir, "varlog")
		require.NoError(t, os.Symlink("/var/log", varLogLink))
		defer os.Remove(varLogLink)

		err := validateDeletionPath(varLogLink)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "symlink resolves")
	})

	t.Run("symlink to deep safe path under dangerous root", func(t *testing.T) {
		// /var/lib/docker/volumes/data/_data is deep enough under /var
		deepLink := filepath.Join(safeDir, "deeplink")
		target := "/var/lib/docker/volumes/data/_data"
		// Only test if target exists
		if _, statErr := os.Stat(target); statErr == nil {
			require.NoError(t, os.Symlink(target, deepLink))
			defer os.Remove(deepLink)

			err := validateDeletionPath(deepLink)
			// Should pass - deep enough under /var
			assert.NoError(t, err)
		}
	})
}

func TestValidateDeletionPath_NonExistent(t *testing.T) {
	// Non-existent paths should still be validated for structure
	t.Run("non-existent but valid structure", func(t *testing.T) {
		err := validateDeletionPath("/tmp/data/nonexistent/file.csv")
		assert.NoError(t, err)
	})

	t.Run("non-existent with traversal", func(t *testing.T) {
		err := validateDeletionPath("/home/user/../../nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "path traversal")
	})
}

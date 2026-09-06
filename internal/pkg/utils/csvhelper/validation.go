package csvhelper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrDangerousPath indicates a path that could cause data loss or security issues
var ErrDangerousPath = errors.New("dangerous path detected")

// dangerousRoots lists directories that should never be deleted
var dangerousRoots = []string{
	// Unix system directories
	"/", "/bin", "/boot", "/dev", "/etc", "/home", "/lib", "/lib64",
	"/proc", "/root", "/sbin", "/sys", "/tmp", "/usr", "/var",
	// Windows system drives (case-insensitive comparison later)
	"C:\\", "D:\\", "E:\\", "C:/", "D:/", "E:/",
}

// validateDeletionPath checks if a path is safe to delete contents from.
// It performs multiple validation layers:
// 1. Detect path traversal attempts (..)
// 2. Block dangerous system directories (including parents)
// 3. Require minimum path depth
// 4. Resolve symlinks to detect escape attempts
func validateDeletionPath(filePath string) error {
	// 1. Detect path traversal attempts BEFORE cleaning
	if strings.Contains(filePath, "..") {
		return fmt.Errorf("%w: path traversal detected in %q", ErrDangerousPath, filePath)
	}

	// 2. Clean and resolve to absolute path
	cleanPath := filepath.Clean(filePath)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// 3. Get the base directory that would be deleted
	baseToDelete := filepath.Dir(absPath)

	// 4. Block dangerous system directories and their parents
	// Check the base directory and all its ancestors
	for _, root := range dangerousRoots {
		// Check if baseToDelete is exactly a dangerous root or under it
		// For Unix: /var/log -> check if it starts with /var
		// For Windows: C:\Windows\System32 -> check if it starts with C:\
		if strings.EqualFold(baseToDelete, root) || strings.HasPrefix(strings.ToLower(baseToDelete), strings.ToLower(root)) {
			// Special case: allow paths that are deeper than the dangerous root
			// e.g., /var/log/docker/volumes is OK, but /var/log is not
			if strings.EqualFold(baseToDelete, root) {
				return fmt.Errorf("%w: refusing to delete contents of system directory %q",
					ErrDangerousPath, baseToDelete)
			}

			// Check if it's a direct child of dangerous root (still dangerous)
			// e.g., /var/log is direct child of /var
			relPath, err := filepath.Rel(root, baseToDelete)
			if err == nil && !strings.Contains(relPath, string(filepath.Separator)) {
				return fmt.Errorf("%w: refusing to delete contents of system directory %q (under %s)",
					ErrDangerousPath, baseToDelete, root)
			}
		}
	}

	// 5. Require minimum depth (at least 3 components: /a/b/c or C:\a\b\c)
	// This prevents: /, /tmp, /home, C:\, etc.
	parts := strings.Split(baseToDelete, string(filepath.Separator))
	nonEmpty := 0
	for _, p := range parts {
		if p != "" {
			nonEmpty++
		}
	}

	// Require at least 3 non-empty components for safety
	// e.g., "tmp/data/files" or "Users/name/projects"
	if nonEmpty < 3 {
		return fmt.Errorf("%w: path too shallow (minimum 3 levels deep): %q",
			ErrDangerousPath, baseToDelete)
	}

	// 6. Resolve symlinks to detect escape attempts
	resolvedPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If path doesn't exist yet, that's OK - use the original path
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to resolve symlinks: %w", err)
		}
		resolvedPath = absPath
	}

	// 7. Re-validate the resolved path if symlinks were involved
	if resolvedPath != absPath {
		resolvedBase := filepath.Dir(resolvedPath)

		// Check if resolved path escapes to dangerous location
		for _, root := range dangerousRoots {
			if strings.EqualFold(resolvedBase, root) || strings.HasPrefix(strings.ToLower(resolvedBase), strings.ToLower(root)) {
				if strings.EqualFold(resolvedBase, root) {
					return fmt.Errorf("%w: symlink resolves to dangerous directory %q",
						ErrDangerousPath, resolvedBase)
				}

				// Check if it's a direct child of dangerous root
				relPath, err := filepath.Rel(root, resolvedBase)
				if err == nil && !strings.Contains(relPath, string(filepath.Separator)) {
					return fmt.Errorf("%w: symlink resolves to dangerous directory %q (under %s)",
						ErrDangerousPath, resolvedBase, root)
				}
			}
		}

		// Check if resolved path has minimum depth
		resolvedParts := strings.Split(resolvedBase, string(filepath.Separator))
		resolvedNonEmpty := 0
		for _, p := range resolvedParts {
			if p != "" {
				resolvedNonEmpty++
			}
		}

		if resolvedNonEmpty < 3 {
			return fmt.Errorf("%w: symlink resolves to shallow path %q",
				ErrDangerousPath, resolvedBase)
		}
	}

	return nil
}

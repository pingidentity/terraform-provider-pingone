package testhcl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// reads HCL file contents used in tests
func ReadTestHcl(path string) (string, error) {
	_, currentPath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current path when reading test HCL")
	}
	fileContents, err := os.ReadFile(filepath.Join(filepath.Dir(currentPath), filepath.Clean(path)))
	if err != nil {
		return "", err
	}

	return string(fileContents), nil
}

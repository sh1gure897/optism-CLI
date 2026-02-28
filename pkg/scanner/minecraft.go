package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// LocateMinecraft finds the default .minecraft directory based on OS.
func LocateMinecraft() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("APPDATA"), ".minecraft")
	case "darwin":
		path = filepath.Join(home, "Library", "Application Support", "minecraft")
	case "linux":
		path = filepath.Join(home, ".minecraft")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("path not found: %s", path)
	}
	return path, nil
}

// LocatePrismInstances detects the Prism Launcher instance storage.
func LocatePrismInstances() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var path string
	if runtime.GOOS == "windows" {
		path = filepath.Join(os.Getenv("APPDATA"), "PrismLauncher", "instances")
	} else if runtime.GOOS == "darwin" {
		path = filepath.Join(home, "Library", "Application Support", "PrismLauncher", "instances")
	}

	if path == "" || isDirMissing(path) {
		return "", fmt.Errorf("prism directory missing")
	}
	return path, nil
}

func isDirMissing(path string) bool {
	info, err := os.Stat(path)
	return os.IsNotExist(err) || !info.IsDir()
}

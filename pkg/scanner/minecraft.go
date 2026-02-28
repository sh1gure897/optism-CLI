package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

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

func LocatePrismInstances() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var path string
	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(os.Getenv("APPDATA"), "PrismLauncher", "instances")
	case "darwin":
		path = filepath.Join(home, "Library", "Application Support", "PrismLauncher", "instances")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("prism directory missing")
	}
	return path, nil
}

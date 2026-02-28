package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppConfig はツールの永続的な設定を保持します
type AppConfig struct {
	Language string `json:"language"`
}

// getConfigPath はOS標準の設定ディレクトリパスを解決します
func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(dir, "optism-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "settings.json"), nil
}

// Load は設定ファイルを読み込みます。存在しない場合は空の設定を返します。
func Load() (*AppConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &AppConfig{Language: ""}, nil // 初回起動時
		}
		return nil, err
	}

	var cfg AppConfig
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

// Save は現在の設定をファイルに書き込みます。
func Save(cfg *AppConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

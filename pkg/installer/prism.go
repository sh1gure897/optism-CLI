package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CreatePrismInstance は新しい起動構成をゼロから構築します
func CreatePrismInstance(prismBaseDir, name, mcVersion, loader, icon string) (string, error) {
	instanceDir := filepath.Join(prismBaseDir, name)
	if err := os.MkdirAll(instanceDir, 0755); err != nil {
		return "", err
	}

	// 1. instance.cfg の作成 (アイコンや表示名の設定)
	cfgContent := fmt.Sprintf("InstanceType=OneSix\niconKey=%s\nname=%s\n", icon, name)
	if err := os.WriteFile(filepath.Join(instanceDir, "instance.cfg"), []byte(cfgContent), 0644); err != nil {
		return "", err
	}

	// 2. mmc-pack.json の作成 (バージョンとローダーの定義)
	pack := map[string]interface{}{
		"formatVersion": 1,
		"components": []map[string]interface{}{
			{
				"uid":       "net.minecraft",
				"version":   mcVersion,
				"important": true,
			},
		},
	}

	// Fabricローダーを構成に追加
	if loader == "fabric" {
		pack["components"] = append(pack["components"].([]map[string]interface{}),
			map[string]interface{}{
				"uid":     "net.fabricmc.intermediary",
				"version": mcVersion,
			},
			map[string]interface{}{
				"uid":     "net.fabricmc.fabric-loader",
				"version": "0.15.7", // 安定板のバージョン
			},
		)
	}

	packData, _ := json.MarshalIndent(pack, "", "  ")
	if err := os.WriteFile(filepath.Join(instanceDir, "mmc-pack.json"), packData, 0644); err != nil {
		return "", err
	}

	// .minecraft フォルダもあらかじめ作っておく
	os.MkdirAll(filepath.Join(instanceDir, ".minecraft"), 0755)

	return instanceDir, nil
}

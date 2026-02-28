package optimizer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InjectConfig は指定されたパスの options.txt を解析し、最適化値を書き込みます
func InjectConfig(mcPath string, plan *OptimizationPlan) error {
	configPath := filepath.Join(mcPath, "options.txt")

	// 既存ファイルの読み込み
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("options.txt が見つかりません（一度ゲームを起動する必要があります）: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// 特定のキーを置換
		if strings.HasPrefix(line, "renderDistance:") {
			line = fmt.Sprintf("renderDistance:%d", plan.RenderDistance)
		} else if strings.HasPrefix(line, "graphicsMode:") {
			mode := 0 // 0: Fast
			if plan.GraphicsMode == "fancy" {
				mode = 1
			}
			line = fmt.Sprintf("graphicsMode:%d", mode)
		} else if strings.HasPrefix(line, "maxFps:") {
			line = fmt.Sprintf("maxFps:%d", plan.MaxFPS)
		} else if strings.HasPrefix(line, "enableVsync:") {
			line = "enableVsync:false" // PvP環境ではVsyncは強制オフ
		}

		lines = append(lines, line)
	}

	// 変更内容をファイルに書き戻す
	output := strings.Join(lines, "\n")
	return os.WriteFile(configPath, []byte(output), 0644)
}

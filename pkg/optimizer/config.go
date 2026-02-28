package optimizer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func InjectConfig(mcPath string, plan *OptimizationPlan) error {
	// フォルダが存在しない場合は作成（新規インスタンス対応）
	if err := os.MkdirAll(mcPath, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(mcPath, "options.txt")
	mode := 0 // Fast
	if plan.GraphicsMode == "fancy" {
		mode = 1
	}

	// ファイルを開く
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// まだ起動していない完全新規のインスタンスなら、最低限の最適化設定だけを書き込んだファイルを生成
			baseConfig := fmt.Sprintf("renderDistance:%d\ngraphicsMode:%d\nmaxFps:%d\nenableVsync:false\n", plan.RenderDistance, mode, plan.MaxFPS)
			return os.WriteFile(configPath, []byte(baseConfig), 0644)
		}
		return fmt.Errorf("failed to access options.txt: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "renderDistance:") {
			line = fmt.Sprintf("renderDistance:%d", plan.RenderDistance)
		} else if strings.HasPrefix(line, "graphicsMode:") {
			line = fmt.Sprintf("graphicsMode:%d", mode)
		} else if strings.HasPrefix(line, "maxFps:") {
			line = fmt.Sprintf("maxFps:%d", plan.MaxFPS)
		} else if strings.HasPrefix(line, "enableVsync:") {
			line = "enableVsync:false"
		}
		lines = append(lines, line)
	}

	output := strings.Join(lines, "\n")
	return os.WriteFile(configPath, []byte(output), 0644)
}

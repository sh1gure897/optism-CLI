//go:build windows

package optimizer

import (
	"fmt"
	"os/exec"
)

// OptimizeNetwork はWindowsのTCPスタックをPvP向けにチューニングします
func OptimizeNetwork() error {
	commands := [][]string{
		// 1. 受信ウィンドウの自動チューニング無効化（パケット到達のブレを抑える）
		{"netsh", "int", "tcp", "set", "global", "autotuninglevel=disabled"},
		// 2. 最新のWindows仕様に合わせたCTCP（輻輳制御）の有効化（ロス時のリカバリ高速化）
		{"netsh", "int", "tcp", "set", "supplemental", "template=internet", "congestionprovider=ctcp"},
	}

	for _, cmdArgs := range commands {
		out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("コマンド [%s] 失敗\n理由: %s", cmdArgs[1], string(out))
		}
	}
	return nil
}

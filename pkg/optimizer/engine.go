package optimizer

import (
	"github.com/sh1gure897/optism-cli/pkg/scanner"
)

// PresetType は最適化の方向性を定義します
type PresetType string

const (
	Competitive PresetType = "CPvP"    // FPS最優先・遅延最小
	Balanced    PresetType = "Balance" // 描画と速度の両立
	Quality     PresetType = "Quality" // 描画距離重視
)

type OptimizationPlan struct {
	JavaXmx        string
	JavaXms        string
	RenderDistance int
	GraphicsMode   string
	MaxFPS         int
	PresetName     PresetType
}

// GeneratePlan はプリセットに基づいて設定を生成します
func GeneratePlan(info *scanner.SystemInfo, pType PresetType) *OptimizationPlan {
	plan := &OptimizationPlan{
		MaxFPS:     260,
		PresetName: pType,
	}

	// RAM割り当て (ここはハードウェア依存)
	if info.TotalRAM_MB <= 8192 {
		plan.JavaXmx = "4G"
		plan.JavaXms = "2G"
	} else {
		plan.JavaXmx = "6G"
		plan.JavaXms = "3G"
	}

	// プリセット別ロジック
	switch pType {
	case Competitive:
		plan.RenderDistance = 6 // CPvPでは視認性と軽さのバランスがここ
		plan.GraphicsMode = "fast"
		plan.MaxFPS = 0 // 無制限
	case Balanced:
		plan.RenderDistance = 10
		plan.GraphicsMode = "fancy"
	case Quality:
		plan.RenderDistance = 16
		plan.GraphicsMode = "fancy"
	}

	return plan
}

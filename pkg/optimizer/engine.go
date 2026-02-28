package optimizer

import (
	"github.com/sh1gure897/optism-cli/pkg/scanner"
)

// OptimizationPlan holds the calculated values for Minecraft settings.
type OptimizationPlan struct {
	JavaXmx        string
	JavaXms        string
	RenderDistance int
	GraphicsMode   string
	MaxFPS         int
}

// GeneratePlan translates raw hardware specs into optimal game settings.
// This logic is tuned for modern PvP where stability > quality.
func GeneratePlan(info *scanner.SystemInfo) *OptimizationPlan {
	plan := &OptimizationPlan{
		MaxFPS: 260, // Uncapped but stable for 144Hz/240Hz monitors
	}

	// RAM Allocation Strategy
	// We keep a 2GB-4GB buffer for the OS to prevent swapping.
	switch {
	case info.TotalRAM_MB <= 4096:
		plan.JavaXmx = "2G"
		plan.JavaXms = "1G"
	case info.TotalRAM_MB <= 8192:
		plan.JavaXmx = "4G"
		plan.JavaXms = "2G"
	default:
		plan.JavaXmx = "6G"
		plan.JavaXms = "3G"
	}

	// Graphic & CPU Strategy
	// Lower cores get aggressive optimization to reduce draw-call overhead.
	if info.CPUCores <= 4 {
		plan.RenderDistance = 8
		plan.GraphicsMode = "fast"
	} else {
		plan.RenderDistance = 12
		plan.GraphicsMode = "fancy"
	}

	return plan
}

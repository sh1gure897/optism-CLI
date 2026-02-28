/*
 * Optism-CLI - Minecraft Performance & Network Optimizer
 * Developed by !時雨/s1gure.dev (sh1gure897)
 * * This tool automates hardware-based configuration injection for
 * competitive Minecraft environments, specifically tuned for Crystal PvP.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sh1gure897/optism-cli/pkg/optimizer"
	"github.com/sh1gure897/optism-cli/pkg/scanner"
)

const Version = "0.1.0-alpha"

func main() {
	printBanner()

	// 1. システムプロファイルの取得
	info, err := scanner.ScanHardware()
	if err != nil {
		log.Fatalf("[CRITICAL] Hardware scan failed: %v", err)
	}

	// 2. 最適化戦略の立案 (ロジックエンジン)
	plan := optimizer.GeneratePlan(info)

	fmt.Printf("[*] Host System: %s (%d Cores, %d MB RAM)\n", info.CPUName, info.CPUCores, info.TotalRAM_MB)
	fmt.Printf("[*] Strategy:    %s RAM Allocation / %d Chunks / Graphics: %s\n",
		plan.JavaXmx, plan.RenderDistance, plan.GraphicsMode)

	// 3. ネットワークレイテンシの最適化 (Windows専用)
	fmt.Print("\n[>] Applying OS-level network tuning...")
	if err := optimizer.OptimizeNetwork(); err != nil {
		fmt.Printf("\n    [!] Skip/Error: %v\n", err)
	} else {
		fmt.Println(" DONE (Zero-Lag Mode Engaged)")
	}

	// 4. バニラ環境への設定注入
	if vPath, err := scanner.LocateMinecraft(); err == nil {
		fmt.Printf("[>] Injecting settings to Vanilla: %s...", vPath)
		if err := optimizer.InjectConfig(vPath, plan); err != nil {
			fmt.Printf(" FAIL: %v\n", err)
		} else {
			fmt.Println(" DONE")
		}
	}

	// 5. Prism Launcher インスタンスの一括処理
	fmt.Println("[>] Scanning Prism Launcher instances...")
	processPrismInstances(plan)

	fmt.Println("\n---------------------------------------")
	fmt.Println("All systems optimized. Dominate the field.")
}

func printBanner() {
	fmt.Printf("⚡ Optism-CLI v%s | Developed by s1gure.dev\n", Version)
	fmt.Println("---------------------------------------")
}

func processPrismInstances(plan *optimizer.OptimizationPlan) {
	prismDir, err := scanner.LocatePrismInstances()
	if err != nil {
		fmt.Printf("    [!] Prism Launcher not detected: %v\n", err)
		return
	}

	entries, _ := os.ReadDir(prismDir)
	var count int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// インスタンス直下、または .minecraft フォルダを探す柔軟な設計
		target := filepath.Join(prismDir, entry.Name())
		if _, err := os.Stat(filepath.Join(target, "options.txt")); os.IsNotExist(err) {
			target = filepath.Join(target, ".minecraft")
		}

		if err := optimizer.InjectConfig(target, plan); err == nil {
			fmt.Printf("    - Optimized: %s\n", entry.Name())
			count++
		}
	}
	fmt.Printf("[+] Processed %d Prism instances.\n", count)
}

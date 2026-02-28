package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sh1gure897/optism-cli/internal/config"
	"github.com/sh1gure897/optism-cli/internal/profiles"
	"github.com/sh1gure897/optism-cli/pkg/installer"
	"github.com/sh1gure897/optism-cli/pkg/optimizer"
	"github.com/sh1gure897/optism-cli/pkg/scanner"
)

const (
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

func main() {
	var targetInstance string
	var customModDir string
	flag.StringVar(&targetInstance, "target", "", "Target Prism instance name")
	flag.StringVar(&customModDir, "mod-dir", "", "Custom mod directory path")
	flag.Parse()

	printBanner()

	cfg, _ := config.Load()
	if cfg == nil {
		cfg = &config.AppConfig{}
	}
	if cfg.Language == "" {
		cfg.Language = promptLanguage()
		config.Save(cfg)
	}

	info, _ := scanner.ScanHardware()
	prismDir, _ := scanner.LocatePrismInstances()

	// プリセット選択
	fmt.Printf("\n%s [1] Competitive (CPvP) [2] Balanced [3] Quality: ", ColorBold+"Select Preset:"+ColorReset)
	var pInput string
	fmt.Scanln(&pInput)
	preset := optimizer.Competitive
	switch pInput {
	case "2":
		preset = optimizer.Balanced
	case "3":
		preset = optimizer.Quality
	}
	plan := optimizer.GeneratePlan(info, preset)

	for {
		lang, _ := profiles.LoadLanguage(cfg.Language)
		fmt.Printf("\n%s\n", ColorCyan+lang.ScanStart+ColorReset)
		fmt.Println("---------------------------------------")
		fmt.Printf(lang.CPUInfo+"\n", info.CPUName, info.CPUCores)
		fmt.Printf(lang.PlanInfo+"\n", plan.JavaXmx, plan.RenderDistance, plan.GraphicsMode)

		if err := optimizer.OptimizeNetwork(); err == nil {
			fmt.Printf("[*] %s\n", lang.NetSuccess)
		}
		fmt.Println("---------------------------------------")

		if targetInstance != "" || customModDir != "" {
			runProcess(targetInstance, customModDir, prismDir, plan, lang)
			break
		}

		choice, selected := interactiveMenu(prismDir, lang)
		if choice == "QUIT" {
			break
		}
		if choice == "LANG" {
			cfg.Language = promptLanguage()
			config.Save(cfg)
			continue
		}

		if choice == "CREATE" {
			selected = createFlow(prismDir, lang)
		}

		if selected != "" {
			runProcess(selected, "", prismDir, plan, lang)
			break // 処理完了後に終了
		}
	}

	lang, _ := profiles.LoadLanguage(cfg.Language)
	fmt.Println("\n" + lang.Finish)
}

func printBanner() {
	banner := `
  ____  _____ _______ _____  _____ __  __ 
 / __ \|  __ \__   __|_   _|/ ____|  \/  |
| |  | | |__) | | |    | | | (___ | \  / |
| |  | |  ___/  | |    | |  \___ \| |\/| |
| |__| | |      | |   _| |_ ____) | |  | |
 \____/|_|      |_|  |_____|_____/|_|  |_|
`
	fmt.Println(ColorCyan + banner + ColorReset)
}

func promptLanguage() string {
	fmt.Print("Select Language (en/ja) [default:en]: ")
	var l string
	fmt.Scanln(&l)
	if l == "ja" {
		return "ja"
	}
	return "en"
}

func runProcess(target, customDir, prismDir string, plan *optimizer.OptimizationPlan, lang *profiles.LanguageBundle) {
	if customDir != "" {
		installer.InstallPerformanceMods(customDir)
		return
	}
	targetPath := filepath.Join(prismDir, target)
	mcPath := filepath.Join(targetPath, ".minecraft")
	if _, err := os.Stat(mcPath); os.IsNotExist(err) {
		mcPath = targetPath
	}

	fmt.Printf("\n"+lang.TargetLock+"\n", target)
	optimizer.InjectConfig(mcPath, plan)

	fmt.Print(lang.ModPrompt)
	var ans string
	fmt.Scanln(&ans)
	if strings.ToLower(ans) == "y" {
		installer.InstallPerformanceMods(filepath.Join(mcPath, "mods"))
	}
}

func interactiveMenu(prismDir string, lang *profiles.LanguageBundle) (string, string) {
	entries, _ := os.ReadDir(prismDir)
	var instances []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			instances = append(instances, e.Name())
		}
	}

	fmt.Printf("\n%s\n", lang.MenuTitle)
	for i, name := range instances {
		fmt.Printf("  "+ColorYellow+"[%d]"+ColorReset+" "+lang.MenuOptimize+"\n", i+1, name)
	}
	fmt.Printf("  [N] %s\n  [L] %s\n  [Q] %s\n", lang.MenuCreate, lang.MenuLang, lang.MenuQuit)
	fmt.Print("\n" + lang.MenuChoice)

	var input string
	fmt.Scanln(&input)
	input = strings.ToUpper(input)

	if input == "Q" {
		return "QUIT", ""
	}
	if input == "L" {
		return "LANG", ""
	}
	if input == "N" {
		return "CREATE", ""
	}

	idx, err := strconv.Atoi(input)
	if err == nil && idx > 0 && idx <= len(instances) {
		return "OPTIMIZE", instances[idx-1]
	}
	return "QUIT", ""
}

func createFlow(prismDir string, lang *profiles.LanguageBundle) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n%s\n", lang.CreateTitle)
	fmt.Print(lang.PromptName)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print(lang.PromptVer)
	ver, _ := reader.ReadString('\n')
	ver = strings.TrimSpace(ver)
	if ver == "" {
		ver = "1.20.1"
	}

	fmt.Print(lang.PromptLoader)
	loader, _ := reader.ReadString('\n')
	loader = strings.TrimSpace(loader)
	if loader == "" {
		loader = "fabric"
	}

	fmt.Print(lang.PromptIcon)
	icon, _ := reader.ReadString('\n')
	icon = strings.TrimSpace(icon)
	if icon == "" {
		icon = "default"
	}

	fmt.Println(lang.ScaffoldLog)
	installer.CreatePrismInstance(prismDir, name, ver, loader, icon)
	return name
}

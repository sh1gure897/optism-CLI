package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
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

func main() {
	var targetInstance string
	var customModDir string
	flag.StringVar(&targetInstance, "target", "", "Target Prism instance name")
	flag.StringVar(&customModDir, "mod-dir", "", "Custom mod directory path")
	flag.Parse()

	cfg, _ := config.Load()
	if cfg == nil {
		cfg = &config.AppConfig{}
	}

	if cfg.Language == "" {
		cfg.Language = promptLanguage()
		config.Save(cfg)
	}

	info, _ := scanner.ScanHardware()
	plan := optimizer.GeneratePlan(info)
	prismDir, _ := scanner.LocatePrismInstances()

	for {
		lang, err := profiles.LoadLanguage(cfg.Language)
		if err != nil {
			log.Fatalf("i18n Load Error: %v", err)
		}

		fmt.Printf("\n%s\n", lang.ScanStart)
		fmt.Println("---------------------------------------")
		fmt.Printf(lang.CPUInfo+"\n", info.CPUName, info.CPUCores)
		fmt.Printf(lang.PlanInfo+"\n", plan.JavaXmx, plan.RenderDistance, plan.GraphicsMode)

		if err := optimizer.OptimizeNetwork(); err == nil {
			fmt.Printf("[*] %s\n", lang.NetSuccess)
		}
		fmt.Println("---------------------------------------")

		if targetInstance != "" || customModDir != "" {
			runDirect(targetInstance, customModDir, prismDir, plan, lang)
			break
		}

		choice, selectedInstance := interactiveMenu(prismDir, lang)

		if choice == "LANG" {
			cfg.Language = promptLanguage()
			config.Save(cfg)
			continue
		} else if choice == "CREATE" {
			targetInstance = createFlow(prismDir, lang)
			if targetInstance != "" {
				runDirect(targetInstance, "", prismDir, plan, lang)
			}
			break
		} else if choice == "OPTIMIZE" {
			runDirect(selectedInstance, "", prismDir, plan, lang)
			break
		} else {
			break
		}
	}

	lang, _ := profiles.LoadLanguage(cfg.Language)
	fmt.Println("---------------------------------------")
	fmt.Println(lang.Finish)
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

func runDirect(target, customDir, prismDir string, plan *optimizer.OptimizationPlan, lang *profiles.LanguageBundle) {
	if customDir != "" {
		// Custom Dir モード
		installer.InstallPerformanceMods(customDir) // 本来はここも翻訳を渡すべきですが、一旦パス表示で対応
		return
	}

	fmt.Printf("\n"+lang.TargetLock+"\n", target)
	targetPath := filepath.Join(prismDir, target)
	mcPath := filepath.Join(targetPath, ".minecraft")
	if _, err := os.Stat(mcPath); os.IsNotExist(err) {
		mcPath = targetPath
	}

	fmt.Println(lang.InjectStart)
	optimizer.InjectConfig(mcPath, plan)

	fmt.Print(lang.ModPrompt)
	var ans string
	fmt.Scanln(&ans)
	if strings.ToLower(ans) == "y" {
		// ※installerパッケージへlangを渡せるようにすると完璧です
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
		fmt.Printf("  [%d] "+lang.MenuOptimize+"\n", i+1, name)
	}
	fmt.Printf("  [N] %s\n", lang.MenuCreate)
	fmt.Printf("  [L] %s\n", lang.MenuLang)
	fmt.Printf("  [Q] %s\n", lang.MenuQuit)
	fmt.Print("\n" + lang.MenuChoice)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToUpper(strings.TrimSpace(input))

	if input == "L" {
		return "LANG", ""
	}
	if input == "N" {
		return "CREATE", ""
	}
	if input == "Q" {
		return "QUIT", ""
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

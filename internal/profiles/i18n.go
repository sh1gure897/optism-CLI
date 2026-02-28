package profiles

import (
	"encoding/json"
	"os"
)

type LanguageBundle struct {
	// スキャン・基本情報
	ScanStart  string `json:"scan_start"`
	CPUInfo    string `json:"cpu_info"`
	PlanInfo   string `json:"plan_info"`
	NetSuccess string `json:"net_success"`

	// メインメニュー
	MenuTitle    string `json:"menu_title"`
	MenuOptimize string `json:"menu_optimize"`
	MenuCreate   string `json:"menu_create"`
	MenuLang     string `json:"menu_lang"`
	MenuQuit     string `json:"menu_quit"`
	MenuChoice   string `json:"menu_choice"`

	// インスタンス作成
	CreateTitle  string `json:"create_title"`
	PromptName   string `json:"prompt_name"`
	PromptVer    string `json:"prompt_ver"`
	PromptLoader string `json:"prompt_loader"`
	PromptIcon   string `json:"prompt_icon"`
	ScaffoldLog  string `json:"scaffold_log"`

	// 最適化・Mod導入
	TargetLock     string `json:"target_lock"`
	InjectStart    string `json:"inject_start"`
	ModPrompt      string `json:"mod_prompt"`
	ModDownloading string `json:"mod_downloading"`
	ModExists      string `json:"mod_exists"`
	ModSaved       string `json:"mod_saved"`

	Finish string `json:"finish"`
}

func LoadLanguage(code string) (*LanguageBundle, error) {
	if code != "ja" && code != "en" {
		code = "en"
	}
	path := "assets/i18n/" + code + ".json"
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bundle LanguageBundle
	err = json.Unmarshal(data, &bundle)
	return &bundle, err
}

package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type ModVersion struct {
	Files []struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
	} `json:"files"`
}

// InstallPerformanceMods は指定されたディレクトリに軽量化Modを直接導入します
func InstallPerformanceMods(modsDir string) error {
	// ユーザーが分かりやすいように絶対パス（フルパス）に変換
	absModsDir, err := filepath.Abs(modsDir)
	if err != nil {
		absModsDir = modsDir
	}

	if err := os.MkdirAll(absModsDir, 0755); err != nil {
		return fmt.Errorf("modsディレクトリの作成に失敗: %v", err)
	}

	fmt.Printf("    [*] Target Directory: %s\n", absModsDir)

	mods := []string{"sodium", "lithium", "entityculling"}
	mcVersion := "1.20.1"

	for _, mod := range mods {
		if err := downloadLatestMod(mod, mcVersion, absModsDir); err != nil {
			fmt.Printf("        [!] %s: %v\n", mod, err)
		}
	}
	return nil
}

func downloadLatestMod(projectID, mcVersion, destDir string) error {
	query := url.Values{}
	query.Add("loaders", `["fabric"]`)
	query.Add("game_versions", fmt.Sprintf(`["%s"]`, mcVersion))

	apiURL := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version?%s", projectID, query.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "sh1gure897/optism-CLI/0.1.0 (s1gure.dev)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API blocked (Status: %d)", resp.StatusCode)
	}

	var versions []ModVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return fmt.Errorf("failed to parse JSON")
	}

	if len(versions) == 0 || len(versions[0].Files) == 0 {
		return fmt.Errorf("no compatible files found")
	}

	fileInfo := versions[0].Files[0]
	destPath := filepath.Join(destDir, fileInfo.Filename)

	if _, err := os.Stat(destPath); err == nil {
		fmt.Printf("      - %s: Already exists.\n        -> %s\n", projectID, destPath)
		return nil
	}

	fmt.Printf("      - %s: Downloading %s...\n", projectID, fileInfo.Filename)

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	fileResp, err := http.Get(fileInfo.URL)
	if err != nil {
		return err
	}
	defer fileResp.Body.Close()

	_, err = io.Copy(out, fileResp.Body)
	if err == nil {
		// 保存先のフルパスを緑色のチェックマークっぽく出力
		fmt.Printf("        [+] Saved to: %s\n", destPath)
	}
	return err
}

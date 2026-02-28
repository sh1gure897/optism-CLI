# ⚡ Optism-CLI

<p align="center">
  <img src="assets/banner.png" alt="Optism-CLI Banner" width="100%">
</p>

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)

Minecraft の実行環境をハードウェアレベルで解析し、パフォーマンスを極限まで引き出すための次世代最適化エンジン。
Developed by [!時雨/s1gure.dev](https://s1gure.dev).

---

## ✨ Key Features

- **Hardware-Aware Optimization**: CPU/RAMをスキャンし、Javaメモリ割当や描画設定を動的に算出.
- **Zero-Lag Network Tuning**: WindowsのTCPスタックを調整し、パケット遅延を最小化.
- **Prism Launcher Sync**: インスタンスを自動検出し、一括または個別に最適化を適用.
- **Mod Management**: Sodium, Lithium などの必須軽量化ModをModrinth APIから自動導入.

## 🏹 Optimization Presets

| プリセット | ターゲット | 主な設定 | 推奨シーン |
| :--- | :--- | :--- | :--- |
| **Competitive** | FPS / 低遅延 | Render: 6 / FPS: Uncapped | Crystal PvP / ランク戦 |
| **Balanced** | 安定性 | Render: 10 / FPS: 260 | 普段使い |
| **Quality** | 描画品質 | Render: 16 / FPS: 144 | 建築 / 影Mod |

## 🛠 Usage

### Quick Start
```bash
# 対話型モードで起動
go run main.go
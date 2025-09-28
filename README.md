# 📊 codex_usage_report

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)

Go tool to analyze **OpenAI Codex** session files (`.jsonl`) and generate reports on token usage, consumption evolution, and cooldown (rate limit recharge) estimates.

---

## 🚀 Features

- Reads multiple session files from `~/.codex/sessions/`  
- Robust **JSON parser** with regex fallback  
- Global aggregated token report  
- Clean timeline (`--timeline`)  
- Full timeline without filtering (`--full-timeline`)  
- Cooldown estimation (rate limit recharge)  
- Supports Linux, macOS, and Windows  
- Option to disable emojis (`--no-emoji`)  

---

## 🛠️ Build

Requirements: **Go 1.21+**

Clone the project and run:

    make build

The binaries will be placed in `dist/`.

To generate builds for **all supported platforms**:

    make release

---

## ▶️ Usage

    ./dist/codex_usage_report [flags]

### Available flags:
- `--timeline` → shows chronological evolution (without repetitions)  
- `--full-timeline` → shows the raw timeline (with repetitions)  
- `--no-emoji` → disables emojis in the output  
- `--debug` → prints detailed parsing logs  

Example:

    ./dist/codex_usage_report --timeline

---

## 💡 How Codex Limits Work

Codex applies two different rate limit counters:

- **Primary** → resets every **5 hours** (short-term limit).  
- **Secondary** → resets every **7 days** (long-term weekly limit).  

This tool helps you estimate **when each limit will recharge**, so you can plan usage better.

---

## 📈 Example Output

```text
📈 Global usage timeline:
  001 | 2025-09-24T21:11:51.883Z → Primary: 1% | Secondary: 0%
  002 | 2025-09-24T21:14:23.962Z → Primary: 3% | Secondary: 1%
  ...
  185 | 2025-09-26T18:38:53.533Z → Primary: 98% | Secondary: 99%
  186 | 2025-09-26T18:38:57.141Z → Primary: 99% | Secondary: 100%
  187 | 2025-09-26T18:39:02.626Z → Primary: 100% | Secondary: 100%

✅ Last values → Primary: 100% | Secondary: 100% (ts=2025-09-26T18:39:02.626Z)  
⏳ Estimated recharge in: 2 days, 23 hours and 13 minutes  
📊 Max total tokens used: 35010120  
📊 Sum of last task tokens: 72496393  

---

## 💻 Installation

### Linux/macOS

1. Build the binary:
   
       make build

2. Copy it into your PATH:
   
       sudo cp dist/codex_usage_report /usr/local/bin/

3. Run from anywhere:
   
       codex_usage_report --timeline

### Windows

1. Build:
   
       make build

   or grab the binary from `dist/codex_usage_report_windows_amd64.exe`.

2. Copy it to a folder in your **PATH** (e.g. `C:\Windows\System32` or set up `C:\bin` in PATH).  

3. Run from anywhere in PowerShell or Command Prompt:
   
       codex_usage_report.exe --timeline

---

## ⚡ Quick install via script

### Linux/macOS

    ./install.sh

The script copies the binary to `/usr/local/bin`.

### Windows

    install.bat

The script copies the `.exe` to your Windows PATH.

---

## 📂 Project structure

    codex_usage_report/
    ├── cmd/
    │   └── codex_usage_report/   # entrypoint (main.go)
    ├── internal/                 # internal project code
    │   ├── parser/               # JSONL parser
    │   ├── timeline/             # timeline + cooldown logic
    │   └── report/               # global summary + printing
    ├── pkg/
    │   └── utils/                # helper functions (timefmt.go)
    ├── dist/                     # build outputs
    ├── Makefile                  # build/release
    ├── go.mod
    └── README.md


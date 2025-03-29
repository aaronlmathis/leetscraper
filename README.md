# 🧠 LeetScraper
[![GitHub release](https://img.shields.io/github/v/release/aaronlmathis/leetscraper?label=Download)](https://github.com/aaronlmathis/leetscraper/releases/latest)

**LeetScraper** is a command-line tool that fetches the [LeetCode](https://leetcode.com) Daily Challenge (or any specific problem) and saves it as well-formatted, language-specific source files — complete with problem description, difficulty, and starter code.

> 🔓 Open Source | 🐹 Written in Go | 📘 GPL-3.0 Licensed

---

## ✨ Features

- ✅ Fetches **daily challenge** or specific problem by slug  
- ✅ Outputs clean, readable **source files** with problem description  
- ✅ Supports **multiple languages** (Go, Python, Java, Rust, etc.)  
- ✅ Fully configurable with `--flags` or `~/.leetscraper.json`  
- ✅ Clean CLI and **GitHub Actions-based release automation**

---

## 📦 Installation

### 🐳 Download prebuilt binary


Precompiled binaries are available inside archives for major platforms:

- **Linux (x86_64)**: [leetscraper-linux.tar.gz](https://github.com/aaronlmathis/leetscraper/releases/latest/download/leetscraper-linux.tar.gz)
- **macOS (arm64)**: [leetscraper-darwin.tar.gz](https://github.com/aaronlmathis/leetscraper/releases/latest/download/leetscraper-darwin.tar.gz)
- **Windows (x86_64)**: [leetscraper-windows.zip](https://github.com/aaronlmathis/leetscraper/releases/latest/download/leetscraper-windows.zip)

Or visit the [Releases page](https://github.com/aaronlmathis/leetscraper/releases) to download and verify checksums.

```sh
# Example (Linux):
curl -LO https://github.com/aaronlmathis/leetscraper/releases/latest/download/leetscraper-linux.tar.gz
tar -xzf leetscraper-linux.tar.gz
chmod +x leetscraper
./leetscraper --help
```

### 🛠 Or build from source:

```bash
git clone https://github.com/aaronlmathis/leetscraper.git
cd leetscraper
make build
./dist/leetscraper --help
```

---

## 🚀 Usage

### ➤ Daily Challenge (default)
Will default to pulling golang snippet, saving to working directory, using filename format: {id}-{difficulty}-{slug}.{ext} if .leetscraper.json is not present in home directory.
```bash
leetscraper
```

### ➤ Specific Problem by Slug

```bash
leetscraper --slug two-sum
```

### ➤ Custom Output Directory

```bash
leetscraper --out ~/leetcode/daily
```

### ➤ Multiple Languages

```bash
leetscraper --langs golang,python3,rust
```
### ➤ File Naming 
```bash
leetscraper --format {id}-{difficulty}-{slug}.{ext}
```
---

## ⚙️ Configuration
Set your preferences (optional) via .leetscraper.json saved to home directory.

See `leetscraper.json.sample` for example.

`~/.leetscraper.json`:

```json
{
  "outputDir": "/home/you/leetcode",
  "filenameFormat": "{id}-{difficulty}-{slug}.{ext}",
  "languages": ["golang", "python3"]
}
```

Command-line flags will **override** config file values.

---

## 📁 Output Example

```text
123-easy-two-sum.go
123-easy-two-sum.py
123-easy-two-sum.rs
```

Each file includes:

- Problem title, difficulty, and link  
- Full problem description (Markdown)  
- Starter function in your language(s)

---

## 🔧 Development

```bash
make build        # Build binary into dist/
make test         # Run integration tests
make release      # Cross-compile for release
make package      # Create .tar.gz and .zip archives
```

---

## 🧪 Test Example

```bash
./test/test_leetscraper.sh
```

Or run as part of `make test`.

---

## 📜 License

**GPL-3.0-or-later**  
© 2025 [Aaron Mathis](mailto:aaron.mathis@gmail.com)

---

## ⭐️ Star This Project

If this tool saves you time, consider [⭐ starring the repo](https://github.com/aaronlmathis/leetscraper) — it helps more developers discover it!

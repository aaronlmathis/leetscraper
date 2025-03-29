# 🧠 LeetScraper

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

```bash
# Coming soon via GitHub Releases
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

---

## ⚙️ Configuration

Optional `~/.leetscraper.json`:

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

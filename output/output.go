/*
SPDX-License-Identifier: GPL-3.0-or-later

Copyright (C) 2025 Aaron Mathis aaron.mathis@gmail.com

This file is part of LeetScraper.

LeetScraper is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

LeetScraper is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with LeetScraper. If not, see https://www.gnu.org/licenses/.
*/

package output

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aaronlmathis/leetscraper/config"
	"github.com/aaronlmathis/leetscraper/scraper"
)

type OutputFile struct {
	Path     string
	Lang     string
	LangSlug string
	Code     string
}

var langExtensions = map[string]string{
	"c":          "c",
	"cpp":        "cpp",
	"csharp":     "cs",
	"java":       "java",
	"javascript": "js",
	"typescript": "ts",
	"python":     "py",
	"python3":    "py",
	"ruby":       "rb",
	"swift":      "swift",
	"kotlin":     "kt",
	"dart":       "dart",
	"golang":     "go",
	"rust":       "rs",
	"scala":      "scala",
	"php":        "php",
	"racket":     "rkt",
	"erlang":     "erl",
	"elixir":     "ex",
	"mysql":      "sql",
	"mssql":      "sql",
	"oraclesql":  "sql",
}

type CommentStyle struct {
	LinePrefix string
	BlockStart string
	BlockEnd   string
}

var commentStyles = map[string]CommentStyle{
	"c":          {BlockStart: "/*", BlockEnd: "*/"},
	"cpp":        {BlockStart: "/*", BlockEnd: "*/"},
	"csharp":     {BlockStart: "/*", BlockEnd: "*/"},
	"java":       {BlockStart: "/*", BlockEnd: "*/"},
	"javascript": {BlockStart: "/*", BlockEnd: "*/"},
	"typescript": {BlockStart: "/*", BlockEnd: "*/"},
	"python":     {BlockStart: `"""`, BlockEnd: `"""`},
	"python3":    {BlockStart: `"""`, BlockEnd: `"""`},
	"golang":     {BlockStart: `/*`, BlockEnd: `*/`},
	"ruby":       {LinePrefix: "#"},
	"swift":      {LinePrefix: "//"},
	"kotlin":     {LinePrefix: "//"},
	"dart":       {LinePrefix: "//"},

	"rust":      {LinePrefix: "//"},
	"scala":     {LinePrefix: "//"},
	"php":       {BlockStart: "/*", BlockEnd: "*/"},
	"racket":    {LinePrefix: ";"},
	"erlang":    {LinePrefix: "%"},
	"elixir":    {LinePrefix: "#"},
	"mysql":     {LinePrefix: "--"},
	"mssql":     {LinePrefix: "--"},
	"oraclesql": {LinePrefix: "--"},
}

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)

func sanitizeTitle(title string) string {
	// Lowercase and replace spaces with dashes
	s := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

	// Remove all characters except letters, numbers, and dashes
	s = unsafeChars.ReplaceAllString(s, "")

	return s
}

func dirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func FormatComment(langSlug, content string) string {
	style, ok := commentStyles[langSlug]
	if !ok {
		style = CommentStyle{LinePrefix: "//"}
	}

	if style.LinePrefix != "" {
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			lines[i] = style.LinePrefix + " " + line
		}
		return strings.Join(lines, "\n")
	}

	return fmt.Sprintf("%s\n%s\n%s", style.BlockStart, content, style.BlockEnd)
}

func formatFilename(format string, problem scraper.LeetCodeProblem, ext, langSlug string) string {

	safeTitle := sanitizeTitle(problem.Title)

	replacer := strings.NewReplacer(
		"{id}", problem.Id,
		"{difficulty}", strings.ToLower(problem.Difficulty),
		"{slug}", problem.TitleSlug,
		"{title}", safeTitle,
		"{ext}", ext,
		"{lang}", langSlug,
	)
	return replacer.Replace(format)
}

func WriteProblemFiles(problem scraper.LeetCodeProblem, cfg *config.Config) error {

	leetcodeDir := filepath.Clean(cfg.OutputDir)

	// Expand "~" if user accidentally used it in config file
	if strings.HasPrefix(leetcodeDir, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not resolve ~ in path: %v", err)
		}
		leetcodeDir = filepath.Join(home, strings.TrimPrefix(leetcodeDir, "~"))

	}
	exists, err := dirExists(leetcodeDir)
	if err != nil {
		return fmt.Errorf("error loading leetcode directory: %v", err)
	}

	if !exists {
		err := os.MkdirAll(leetcodeDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating leetcode directory: %v", err)
		}
	}

	for _, snippet := range problem.CodeSnippets {
		ext, ok := langExtensions[snippet.LangSlug]
		if !ok {
			log.Printf("Error generating file for %v. Unknown language", snippet.LangSlug)
			continue
		}
		filename := formatFilename(cfg.FilenameFormat, problem, ext, snippet.LangSlug)
		fullPath := filepath.Join(leetcodeDir, filename)
		fullUrl := fmt.Sprintf("http://leetcode.com/problems/%s/", problem.TitleSlug)

		header := fmt.Sprintf(
			"%s\n\n%s",
			fmt.Sprintf("#%s - %s\nDifficulty: %s\n%s", problem.Id, problem.Title, problem.Difficulty, fullUrl),
			problem.Description,
		)

		commentHeader := FormatComment(snippet.LangSlug, header)

		content := fmt.Sprintf("%s\n\n%s", commentHeader, snippet.Code)

		f, err := os.Create(fullPath)
		if err != nil {
			log.Printf("Failed to create file: %s", fullPath)
			continue
		}

		_, err = f.WriteString(content)
		if err != nil {
			log.Printf("Failed to write file: %s", fullPath)
			return err
		}

		if err := f.Close(); err != nil {
			log.Printf("Error closing file %s: %v", fullPath, err)
			return err
		}

	}

	return nil
}

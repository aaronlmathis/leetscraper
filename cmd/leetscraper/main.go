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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronlmathis/leetscraper/config"
	"github.com/aaronlmathis/leetscraper/internal/version"
	"github.com/aaronlmathis/leetscraper/output"
	"github.com/aaronlmathis/leetscraper/scraper"
)

func trimSplit(input, sep string) []string {
	parts := strings.Split(input, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func main() {
	// Load config from file if it exists
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".leetscraper.json")

	var cfg *config.Config
	var err error

	if _, err := os.Stat(configPath); err == nil {
		cfg, err = config.LoadFromFile(configPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	} else {
		cfg = config.Default()
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Define CLI flags using config values as defaults
	slug := flag.String("slug", "", "LeetCode problem titleSlug (e.g., two-sum)")
	outDir := flag.String("out", cfg.OutputDir, "Output directory (e.g., /home/leetcoder/daily)")
	langs := flag.String("langs", strings.Join(cfg.Languages, ","), "Comma-separated list of preferred languages (e.g., python3,golang,ruby,cpp)")
	format := flag.String("format", cfg.FilenameFormat, "Filename format (e.g., {id}-{slug}-{title}.{ext})")
	versionFlag := flag.Bool("version", false, "Print LeetScraper version and exit")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("LeetScraper version: %s\n", version.Version)
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Unknown argument(s): %v\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	// Apply CLI overrides to config
	cfg.OutputDir = *outDir
	cfg.FilenameFormat = *format
	cfg.Languages = trimSplit(*langs, ",")

	// Fetch problem (by slug or daily)
	var problem scraper.LeetCodeProblem

	if *slug != "" {
		problem, err = scraper.GetProblemDetails(*slug, cfg.Languages)
		if err != nil {
			log.Fatalf("Failed to get problem details: %v", err)
		}
	} else {
		daily, err := scraper.GetDailyCodingChallengeQuestion()
		if err != nil {
			log.Fatalf("Failed to get daily challenge: %v", err)
		}
		problem, err = scraper.GetProblemDetails(daily.TitleSlug, cfg.Languages)
		if err != nil {
			log.Fatalf("Failed to get daily problem details: %v", err)
		}
	}

	log.Printf("📘 Problem: #%s - %s (%s)", problem.Id, problem.Title, problem.Difficulty)

	if err := output.WriteProblemFiles(problem, cfg); err != nil {
		log.Fatalf("Failed to write output files: %v", err)
	}
}

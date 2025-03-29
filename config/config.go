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

package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	OutputDir      string   `json:"outputDir"`
	FilenameFormat string   `json:"filenameFormat"`
	Languages      []string `json:"languages"`
}

// Default returns fallback values for config
func Default() *Config {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "." // fallback to relative if os.Getwd fails
	}
	return &Config{
		OutputDir:      cwd,
		FilenameFormat: "{id}-{difficulty}-{slug}.{ext}",
		Languages:      []string{"golang"},
	}
}

// LoadFromFile reads config from a JSON file (returns partial config if missing)
func LoadFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	cfg := Default()
	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}

// MergeFlags overrides values in the given config using CLI flags
func MergeFlags(cfg *Config) {
	outputDir := flag.String("out", cfg.OutputDir, "Output directory")
	format := flag.String("format", cfg.FilenameFormat, "Filename format")
	langs := flag.String("langs", strings.Join(cfg.Languages, ","), "Preferred languages (comma-separated)")

	flag.Parse()

	cfg.OutputDir = *outputDir
	cfg.FilenameFormat = *format
	if *langs != "" {
		cfg.Languages = strings.Split(*langs, ",")
	}
}

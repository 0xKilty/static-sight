package main

import (
	"gopkg.in/yaml.v2"
	"time"
)

type frontMatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
	Latex bool     `yaml:"latex"`
	Code  bool     `yaml:"code"`
}

func parseFrontMatter(content string) (frontMatter, int) {
	var frontMatter frontMatter
	frontMatterEnd := 0
	if content[0:3] == "---" {
		frontMatterEnd = 3
		for i := 3; i < len(content); i++ {
			if content[i:i + 3] == "---" {
				frontMatterEnd = i + 3
				break
			}
		}
	}
	if frontMatterEnd > 0 {
		if err := yaml.Unmarshal([]byte(content[3:frontMatterEnd-3]), &frontMatter); err != nil {
			return frontMatter, -1
		}
	}
	return frontMatter, frontMatterEnd
}

func frontMatterToArticleEntry(frontMatter frontMatter, filePath string) articleEntry {
	date, err := time.Parse("1/2/2006", frontMatter.Date)
	check(err)
	return articleEntry{Title: frontMatter.Title, Date: date, Tags: frontMatter.Tags, FilePath: filePath}
}
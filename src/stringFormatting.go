package main

import (
	"strings"
	"github.com/gomarkdown/markdown"
)

func slugify(str string) string {
	return strings.ReplaceAll(strings.ToLower(str), " ", "-")
}

func MDtoHTML(str string) string {
	return strings.Replace(str, ".md", ".html", 1)
}

func removeDotMD(str string) string {
	return strings.Replace(str, ".md", "", 1)
}

func splitBySlash(str string) []string {
	return strings.Split(str, "/")
}

func getBuildDirPath(path string) string {
	buildDirSlice := splitBySlash(path)
	buildDirSlice[0] = "build"
	buildDirSlice[1] = strings.ToLower(buildDirSlice[1])
	buildDirString := "./" + strings.Join(buildDirSlice, "/")
	return buildDirString
}

func getLocalPath(wd string) string {
	slice := splitBySlash(wd)
	path := strings.Join(slice[2:], "/")
	return path
}

func getMarkdown(bytes []byte) string {
	return string(markdown.ToHTML(bytes, nil, nil))
}
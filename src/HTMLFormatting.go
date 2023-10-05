package main

import (
	"fmt"
	"strconv"
	"time"
)

type articleEntry struct {
	FilePath string
	Date     time.Time
	Title    string
	Tags     []string
}

func formatRegularEntryHTML(entry articleEntry) string {
	tags := getTagsSting(entry.Tags)
	return fmt.Sprintf("<li><a href=\"%s\">%s</a> - <em>%s</em>&nbsp;&nbsp;<div class=\"tags\">%s</div></li>", slugify(removeDotMD(entry.FilePath)), entry.Title, entry.Date.Format("1/2/2006"), tags)
}

func formatTagsHTML(tags []string) string {
	var tagsString string
	for _, tag := range tags {
		tagsString += fmt.Sprintf("<a href=\"%s\">%s</a>&nbsp;", "/tags/" + slugify(removeDotMD(tag)), tag)
	}
	return tagsString
}

func getTagsSting(tags []string) string {
	var tagsString string
	for _, tag := range tags {
		tagsString += fmt.Sprintf("<a href=\"/tags/%s\">%s</a>&nbsp;&nbsp;", slugify(removeDotMD(tag)), tag)
	}
	return tagsString
}

func getNotesDirPageContent(dirPath string) (string, string) {
	var readmeContent string
	if checkExists(dirPath + "/README.md") {
		readmeContent = getMarkdown(readFile(dirPath + "/README.md"))
	}
	return getFileStructureHTML(dirPath, 1), readmeContent
}

func getFileStructureHTML(dirPath string, depth int) string {
	var content string
	noteDir := readDir("./content/" + dirPath)
	arrowPadding := strconv.Itoa(((depth - 1) * 5) + 15)
	normalPadding := strconv.Itoa(((depth) * 5) + 15)
	for _, item := range noteDir {
		url := slugify("/" + dirPath + "/" + removeDotMD(item.Name()))
		if item.IsDir() {
			content += fmt.Sprintf("<div><div class=\"dirName\"><i class=\"arrow right\" onclick=\"arrowClick(this)\"></i></div>&nbsp;<strong><a href=\"%s\">", url)
			content += item.Name()
			content += fmt.Sprintf("</a></strong><div style=\"padding-left: %spx;\" hidden>", normalPadding)
			innerContent := getFileStructureHTML(dirPath + "/" + item.Name(), depth + 1)
			content += innerContent
			content += "</div></div>"
		} else if item.Name() != "README.md" {
			content += fmt.Sprintf("<a href=\"%s\" style=\"padding-left: %spx;\">%s</a><br>", url, arrowPadding, removeDotMD(item.Name()))
		}
	}
	return content
}

func getRegularDirPageContent(localPath string) (string, string) {
	var content string
	entries := readDir("./content/" + localPath)
	for _, entry := range entries {
		if entry.Name() != "README" {
			file := readFile("./content/" + localPath + "/" + entry.Name())
			frontMatter, _ := parseFrontMatter(string(file))
			content += formatRegularEntryHTML(frontMatterToArticleEntry(frontMatter, "/"+localPath+"/"+slugify(removeDotMD(entry.Name()))))
		}
	}
	return content, ""
}
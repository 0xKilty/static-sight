package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	clearBuild()
	copyAssets()
	postOrderTraversal("./content")
	updateIndexAndTags()

	os.Chdir("./build")

	localhostPort := "3000"
	hostBuild(localhostPort)
}

type articlePage struct {
	Path    string
	Title   string
	Date    string
	Tags    string
	Content string
	Latex   bool
	Code    bool
}

type notePage struct {
	Title string
	FilePath string
	Content string
}

type directoryPage struct {
	Title   string
	Path    string
	Content string
	ReadMe  string
}

func getPathLinks(path string, caser cases.Caser) string {
	var paths string
	pathLinks := "/"
	for _, file := range strings.Split(path, "/") {
		cleanFile := removeDotMD(file)
		paths += "/" + cleanFile
		pathLinks += fmt.Sprintf("<a href=\"%s\">%s</a>/", slugify(paths), caser.String(cleanFile))
	}
	return pathLinks
}

func postOrderTraversal(root string) error {
	caser := cases.Title(language.English)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		check(err)
		buildDirPath := getBuildDirPath(path)
		localPath := getLocalPath(buildDirPath)
		pathSlice := splitBySlash(buildDirPath)
		lastPathSliceElement := pathSlice[len(pathSlice)-1]
		topDirectory := splitBySlash(localPath)[0]
		var template string

		if info.IsDir() {
			if path == root || lastPathSliceElement == "README" { return nil }

			if !checkExists(buildDirPath) { createDir(slugify(buildDirPath)) }

			var content, readmeContent string
			if topDirectory == "blog" || topDirectory == "projects" {
				content, readmeContent = getRegularDirPageContent(localPath)
				template = "dir_temp.html"
			} else {
				content, readmeContent = getNotesDirPageContent(localPath)
				template = "note_dir_temp.html"
			}
			dirIndexHTML := openOrCreateFile(slugify(buildDirPath) + "/index.html")
			pathLinks := getPathLinks(localPath, caser)
			dirPage := directoryPage{Title: caser.String(lastPathSliceElement), Path: pathLinks, Content: content, ReadMe: readmeContent}
			writeFileWithTemplate(dirIndexHTML, template, dirPage)

			return nil
		}

		file := readFile(path)
		fileOut := openOrCreateFile(slugify(removeDotMD(buildDirPath)))
		var contentStruct interface{}
		if topDirectory == "blog" || topDirectory == "projects" {
			frontMatter, frontMatterEnd := parseFrontMatter(string(file))
			contentStruct = articlePage {
				Title: frontMatter.Title, 
				Path: getPathLinks(localPath, caser),
				Tags: formatTagsHTML(frontMatter.Tags), 
				Date: frontMatter.Date, 
				Content: getMarkdown(file[frontMatterEnd:]), 
				Latex: frontMatter.Latex, 
				Code: frontMatter.Code,
			}
			template = "art_temp.html"
		} else {
			contentStruct = notePage {
				Title: removeDotMD(lastPathSliceElement), 
				Content: getMarkdown(file), 
				FilePath: getPathLinks(localPath, caser),
			}
			template = "note_temp.html"
		}
		writeFileWithTemplate(fileOut, template, contentStruct)
		return nil
	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	clearBuild()
	copyAssets()
	postOrderTraversal("./content")
	updateIndexAndTags()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("(1) Rebuild (default)\n(2) Deploy\n> ")
    name, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Hello,", name)

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

func postOrderTraversal(root string) error {
	caser := cases.Title(language.English)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		check(err)
		buildDirPath := getBuildDirPath(path)
		localPath := getLocalPath(buildDirPath)
		pathSlice := splitBySlash(buildDirPath)
		lastPathSliceElement := pathSlice[len(pathSlice)-1]
		topDirectory := splitBySlash(localPath)[0]
		pathLinks := getPathLinks(localPath, caser)
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
			dirPage := directoryPage{Title: caser.String(lastPathSliceElement), Path: pathLinks, Content: content, ReadMe: readmeContent}
			writeFileWithTemplate(dirIndexHTML, template, dirPage)

			return nil
		}

		file := readFile(path)
		fileOut := openOrCreateFile(slugify(removeDotMD(buildDirPath)))
		var contentStruct interface{}
		
		if topDirectory == "blog" || topDirectory == "projects" {
			contentStruct = getRegularPageContent(file, pathLinks)
			template = "art_temp.html"
		} else {
			contentStruct = getNotePageContent(file, removeDotMD(lastPathSliceElement), pathLinks)
			template = "note_temp.html"
		}
		writeFileWithTemplate(fileOut, template, contentStruct)
		return nil
	})
}

func getRegularPageContent(file []byte, pathLinks string) (articlePage) {
	frontMatter, frontMatterEnd := parseFrontMatter(string(file))
	contentStruct := articlePage {
		Title: frontMatter.Title, 
		Path: pathLinks,
		Tags: formatTagsHTML(frontMatter.Tags), 
		Date: frontMatter.Date, 
		Content: getMarkdown(file[frontMatterEnd:]), 
		Latex: frontMatter.Latex, 
		Code: frontMatter.Code,
	}
	return contentStruct
}

func getNotePageContent(file []byte, title string, pathLinks string) (notePage) {
	contentStruct := notePage {
		Title: title,
		Content: getMarkdown(file), 
		FilePath: pathLinks,
	}
	return contentStruct
}

func getAllDirContent() (string, string) {
	return 
}
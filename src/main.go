package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	build()
	os.Chdir("./build")
	localhostPort := "3000"
	go hostBuild(localhostPort)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("(1) Rebuild (default)\n(2) Deploy\n(3) Quit\n> ")
		input, err := reader.ReadString('\n')
		fmt.Print("\033[H\033[2J")
		check(err)
		if input == "1\n" || input == "" {
			build()
			fmt.Print("\033[H\033[2J")
			fmt.Print("Rebuilt")
		} else if input == "2\n" {
			fmt.Println("Deploying")
		} else if input == "3" {
			fmt.Println("Quitting")
			break
		} else {
			fmt.Println("Please enter either 1, 2, or 3")
		}
	}
	os.Chdir("./build")
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
		inBlogOrProjects := topDirectory == "blog" || topDirectory == "projects"
		pathLinks := getPathLinks(localPath, caser)
		var template string

		if info.IsDir() {
			if path == root { return nil }
			if !checkExists(buildDirPath) { createDir(slugify(buildDirPath)) }

			var content, readmeContent string
			if inBlogOrProjects {
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

		if lastPathSliceElement == "README" { return nil }
		file := readFile(path)
		fileOut := openOrCreateFile(slugify(removeDotMD(buildDirPath)))
		var contentStruct interface{}

		if inBlogOrProjects {
			contentStruct, template = getRegularPageContent(file, pathLinks)
		} else {
			contentStruct, template = getNotePageContent(file, removeDotMD(lastPathSliceElement), pathLinks)
		}
		writeFileWithTemplate(fileOut, template, contentStruct)
		return nil
	})
}
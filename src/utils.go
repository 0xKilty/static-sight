package main

import (
	"fmt"
	"net/http"
	"os"
)

func copyAssets() {
	files := readDir("./templates")
	for _, file := range files {
		if file.IsDir() {
			createDir("./build/" + file.Name())
			dirFiles := readDir("./templates/" + file.Name())
			for _, dirFile := range dirFiles {
				copyFileTemplatesToBuild(file.Name() + "/" + dirFile.Name())
			}
		} else if file.Name() == "resume.pdf" {
			copyFileTemplatesToBuild(file.Name())
		}
	}
}

func copyFileTemplatesToBuild(fileName string) {
	assetContents := readFile("./templates/" + fileName)
	fileOut := openOrCreateFile("./build/" + fileName)
	writeFile(fileOut, assetContents)
}

func clearDir(path string) {
	err := os.RemoveAll("./build")
	check(err)
	createDir("./build")
}

func hostBuild(port string) {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	fmt.Println("Serving site on localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func build() {
	clearDir("./build")
	copyAssets()
	postOrderTraversal("./content")
	updateIndexAndTags()
}

package main

import (
	"os"
	"fmt"
	"net/http"
)

func copyAssets() {
	files := readDir("./templates")
	for _, file := range files {
		if file.IsDir() {
			createDir("./build/" + file.Name())
			dirFiles := readDir("./templates/" + file.Name())
			for _, dirFile := range dirFiles {
				assetContents := readFile("./templates/" + file.Name() + "/" + dirFile.Name())
				fileOut := openOrCreateFile("./build/" + file.Name() + "/" + dirFile.Name())
				writeFile(fileOut, assetContents)
			}
		}
	}
}

func clearBuild() {
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
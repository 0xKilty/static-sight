package main

import (
	"os"
	"io"
	"io/fs"
	"text/template"
)

func createDir(name string) {
	err := os.Mkdir(name, 0777)
	check(err)
}

func readDir(path string) []fs.DirEntry {
	files, err := os.ReadDir(path)
	check(err)
	return files
}

func openOrCreateFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	return file
}

func readFile(path string) []byte {
	file, err := os.Open(path)
	check(err)
	contents, err := io.ReadAll(file)
	check(err)
	err = file.Close()
	check(err)
	return contents
}

func writeFile(file *os.File, content []byte) {
	_, err := file.Write(content)
	check(err)
}

func writeFileWithTemplate(file *os.File, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles("./templates/" + templatePath)
	check(err)
	tmpl.Execute(file, data)
	check(err)
}

func checkExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
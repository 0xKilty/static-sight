package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
	"github.com/gomarkdown/markdown"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

func main() {
	postOrderTraversal("./content")
	updateIndex()

	os.Chdir("./build")

	localhostPort := "3000"
	hostBuild(localhostPort)
}

func hostBuild(port string) {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	fmt.Println("Serving site on localhost:" + port);
	http.ListenAndServe(":" + port, nil)
}

type articleEntry struct {
	FilePath string
	Date     time.Time
	Title    string
	Tags     []string
}

type tag struct {
	Name    string
	Entries []articleEntry
}

type indexPage struct {
	RecentPosts    string
	RecentProjects string
}

type articlePage struct {
	Title   string
	Date    string
	Content string
	Latex   bool
	Code    bool
}

type directoryPage struct {
	Title   string
	Path    string
	Content string
	ReadMe  string
}

type frontMatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
	Latex bool     `yaml:"latex"`
	Code  bool     `yaml:"code"`
}

func removeDotMD(str string) string {
	return strings.Replace(str, ".md", "", 1)
}

func addArticleEntryInOrder(slice []articleEntry, path string, title string, tags []string, date time.Time) []articleEntry {
	index := sort.Search(len(slice), func(i int) bool {
		return !slice[i].Date.Before(date)
	})

	slice = append(slice, articleEntry{})
	copy(slice[index+1:], slice[index:])
	slice[index].Date = date
	slice[index].FilePath = path
	slice[index].Title = title
	slice[index].Tags = tags

	return slice
}

func add_limit_convert_html(recent_blogs []articleEntry, limit int) string {
	blogs_content := ""
	for i, e := range recent_blogs {
		if i == limit {
			break
		}
		blogs_content += "<li><a href=\"" + removeDotMD(e.FilePath) + "\">" + e.Title + "</a> - <em>" + e.Date.Format("1/2/2006") + "</em></li>"
	}
	return blogs_content
}

func updateIndex() {
	recent_blogs := get_ordered_dates("blog")
	recent_projects := get_ordered_dates("projects")
	limit := 5
	blogs_content := add_limit_convert_html(recent_blogs, limit)
	projects_content := add_limit_convert_html(recent_projects, limit)

	file_out, err := os.OpenFile("./build/index.html", os.O_RDWR|os.O_CREATE, 0755)
	check(err)

	home := indexPage{RecentPosts: blogs_content, RecentProjects: projects_content}

	tmpl, err := template.ParseFiles("./templates/index_temp.html")
	check(err)

	tmpl.Execute(file_out, home)
	check(err)
}

func get_ordered_dates(folder string) []articleEntry {
	blog_posts, err := os.ReadDir("./content/" + folder)
	check(err)
	var date_slice []articleEntry
	for _, e := range blog_posts {
		if e.Name() != "README" {
			file, err := os.Open("./content/" + folder + "/" + e.Name())
			check(err)
			contents, err := io.ReadAll(file)
			check(err)
			fm, _ := parseFrontMatter(string(contents))
			check(err)
			date, err := time.Parse("1/2/2006", fm.Date)
			check(err)
			date_slice = addArticleEntryInOrder(date_slice, "/"+folder+"/"+e.Name(), fm.Title, fm.Tags, date)
		}
	}
	return date_slice
}

func getBuildDirPath(path string) string {
	build_path := strings.Split(path, "/")
	build_path[0] = "build"
	build_path[1] = strings.ToLower(build_path[1])
	build_string := "./" + strings.Join(build_path, "/")
	return build_string
}

func checkExistsDir(dir string) bool {
	_, err := os.Stat(dir); 
	return err == nil
}

func createDir(name string) {
	err := os.Mkdir(name, 0777)
	check(err)
}

func getLocalPath(wd string) string {
	slice := splitBySlash(wd)
	path := strings.Join(slice[2:], "/")
	return path
}

func splitBySlash(str string) []string {
	return strings.Split(str, "/")
}

func readDir(path string) []fs.DirEntry {
	files, err := os.ReadDir("./content/" + path)
	check(err)
	return files;
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

func openOrCreateFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	return file
}

func formatEntryHTML(entry fs.DirEntry, localPath string) string {
	contentFilePath := removeDotMD(entry.Name())
	content := "<p><a href=\"/" + localPath + "/" + contentFilePath + "\" class=\"black\">"

	// Make it bold if its a dir
	if entry.IsDir() { 
		contentFilePath = "<strong>" + contentFilePath + "</strong>" 
	}
	content += contentFilePath + "</a></p>\n"
	return content
}

func getDirPageContent(contentPath string, localPath string) (string, string) {
	var content, readmeContent string;
	entries := readDir(localPath)
	for _, entry := range entries {
		if entry.Name() == "README" {
			fileContent := readFile(contentPath + "/" + entry.Name())
			readmeContent = getMarkdown(fileContent)
		} else if entry.Name() != "index.html" {
			content += formatEntryHTML(entry, localPath)
		}
	}
	return content, readmeContent
}

func getPathLinks(path string, caser cases.Caser) string {
	var paths string;
	pathLinks := "/"
	for _, file := range strings.Split(path, "/") {
		paths += "/" + file
		fmt.Println(paths, file)
		pathLinks += "<a href=\"" + paths + "\" class=\"black\">" + caser.String(file) + "</a>/"
	}
	return pathLinks
}

func writeFileWithTemplate(file *os.File, templatePath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	check(err)
	tmpl.Execute(file, data)
	check(err)
}

func getMarkdown(bytes []byte) string {
	return string(markdown.ToHTML(bytes, nil, nil))
}

func postOrderTraversal(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		buildDirPath := getBuildDirPath(path)
		localPath := getLocalPath(buildDirPath)
		pathSlice := splitBySlash(buildDirPath)
		lastPathSliceElement := pathSlice[len(pathSlice)-1]

		if info.IsDir() {
			if path == root { return nil }
			
			if !checkExistsDir(buildDirPath) { createDir(buildDirPath) }

			content, readmeContent := getDirPageContent("./content/" + localPath, localPath)

			dirIndexHTML := openOrCreateFile(buildDirPath + "/index.html")

			caser := cases.Title(language.English)
			pathLinks := getPathLinks(localPath, caser)
			current_page := directoryPage{Title: caser.String(lastPathSliceElement), Path: pathLinks, Content: content, ReadMe: readmeContent}

			writeFileWithTemplate(dirIndexHTML, "./templates/dir_temp.html", current_page)

			return nil
		}

		if lastPathSliceElement == "README" {
			return nil
		}

		file := readFile(path)

		frontMatter, frontMatterEnd := parseFrontMatter(string(file))

		content := getMarkdown(file[frontMatterEnd:])
		fileOut := openOrCreateFile(buildDirPath)

		articlePage := articlePage{ Title: frontMatter.Title, Date: frontMatter.Date, Content: content, Latex: frontMatter.Latex, Code: frontMatter.Code }

		writeFileWithTemplate(fileOut, "./templates/art_temp.html", articlePage)

		return nil
	})
}

func parseFrontMatter(content string) (frontMatter, int) {
	var frontMatter frontMatter
	frontMatterEnd := 0
	if content[0:3] == "---" {
		frontMatterEnd = 3
		for i := 3; i < len(content); i++ {
			if content[i:i+3] == "---" {
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

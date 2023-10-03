package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
	"strconv"
)

func main() {
	clearBuild()
	createDir("./build")
	postOrderTraversal("./content")
	updateIndexAndTags()
	copyAssets()

	os.Chdir("./build")

	localhostPort := "3000"
	hostBuild(localhostPort)
}

func hostBuild(port string) {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	fmt.Println("Serving site on localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

type articleEntry struct {
	FilePath string
	Date     time.Time
	Title    string
	Tags     []string
}

type indexPage struct {
	RecentPosts    string
	RecentProjects string
}

type articlePage struct {
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

type noteDirectoryPage struct {
	Content string
	Path    string
	ReadMe  string
}

type tagPage struct {
	Title   string
	AllTags string
	Content string
}

type frontMatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
	Latex bool     `yaml:"latex"`
	Code  bool     `yaml:"code"`
}

type DirOrFile bool

type noteEntry struct {
	FileType DirOrFile
	Path     string
}

const (
	Dir  DirOrFile = true
	File DirOrFile = false
)

func removeDotMD(str string) string {
	return strings.Replace(str, ".md", "", 1)
}

func addArticleEntryInOrder(slice []articleEntry, entry articleEntry) []articleEntry {
	index := sort.Search(len(slice), func(i int) bool {
		return !slice[i].Date.Before(entry.Date)
	})

	slice = append(slice, articleEntry{})
	copy(slice[index+1:], slice[index:])
	slice[index].Date = entry.Date
	slice[index].FilePath = entry.FilePath
	slice[index].Title = entry.Title
	slice[index].Tags = entry.Tags

	return slice
}

func addLimitConvertHtml(recentPosts []articleEntry, limit int) string {
	var blogContent string
	for index, entry := range recentPosts {
		if index == limit {
			break
		}
		blogContent += "<li><a href=\"" + removeDotMD(entry.FilePath) + "\">" + entry.Title + "</a> - <em>" + entry.Date.Format("1/2/2006") + "</em></li>"
	}
	return blogContent
}

func urlFormat(str string) string {
	return strings.ToLower(strings.ReplaceAll(str, " ", "-"))
}

func getAllTags(tagsMap map[string][]articleEntry) string {
	var allTags string
	for tagName := range tagsMap {
		allTags += "&nbsp;<a href=\"/tags/" + urlFormat(tagName) + "\">" + tagName + "</a>"
	}
	return allTags
}

func getTagPageContent(tagEntries []articleEntry) string {
	var tagPageContent string
	for _, tagEntry := range tagEntries {
		tagPageContent += formatRegularEntryHTML(tagEntry)
	}
	return tagPageContent
}

func updateIndexAndTags() {
	numberOfRecent := 5
	tagsMap := make(map[string][]articleEntry)

	recentBlogPosts, tagsMap := getFolderInfo("blog", tagsMap)
	recentProjects, tagsMap := getFolderInfo("projects", tagsMap)

	blogsContent := addLimitConvertHtml(recentBlogPosts, numberOfRecent)
	projectsContent := addLimitConvertHtml(recentProjects, numberOfRecent)

	fileOut := openOrCreateFile("./build/index.html")
	index := indexPage{RecentPosts: blogsContent, RecentProjects: projectsContent}
	writeFileWithTemplate(fileOut, "./templates/index_temp.html", index)

	if !checkExists("./build/tags") {
		createDir("./build/tags")
	}

	allTags := getAllTags(tagsMap)

	for tagName, tagEntries := range tagsMap {
		tagPageContent := getTagPageContent(tagEntries)
		tagHTMLFile := openOrCreateFile("./build/tags/" + strings.ToLower(tagName))
		tagPageData := tagPage{Title: tagName, AllTags: allTags, Content: tagPageContent}
		writeFileWithTemplate(tagHTMLFile, "./templates/tags_temp.html", tagPageData)
	}
}

func insertTag(tagsMap map[string][]articleEntry, entry articleEntry) map[string][]articleEntry {
	for _, articleTag := range entry.Tags {
		tagsMap[articleTag] = append(tagsMap[articleTag], entry)
	}
	return tagsMap
}

func getFolderInfo(folder string, tagsMap map[string][]articleEntry) ([]articleEntry, map[string][]articleEntry) {
	dir := readDir("./content/" + folder)
	var articlesByDate []articleEntry
	for _, entry := range dir {
		if entry.Name() != "README" {
			file := readFile("./content/" + folder + "/" + entry.Name())
			frontMatter, _ := parseFrontMatter(string(file))
			date, err := time.Parse("1/2/2006", frontMatter.Date)
			check(err)
			entry := articleEntry{FilePath: "/" + folder + "/" + entry.Name(), Title: frontMatter.Title, Date: date, Tags: frontMatter.Tags}
			articlesByDate = addArticleEntryInOrder(articlesByDate, entry)
			tagsMap = insertTag(tagsMap, entry)
		}
	}
	return articlesByDate, tagsMap
}

func getBuildDirPath(path string) string {
	buildDirSlice := splitBySlash(path)
	buildDirSlice[0] = "build"
	buildDirSlice[1] = strings.ToLower(buildDirSlice[1])
	buildDirString := "./" + strings.Join(buildDirSlice, "/")
	return buildDirString
}

func checkExists(dir string) bool {
	_, err := os.Stat(dir)
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
	files, err := os.ReadDir(path)
	check(err)
	return files
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

func frontMatterToArticleEntry(frontMatter frontMatter, filePath string) articleEntry {
	date, err := time.Parse("1/2/2006", frontMatter.Date)
	check(err)
	return articleEntry{Title: frontMatter.Title, Date: date, Tags: frontMatter.Tags, FilePath: filePath}
}

func getNotesPageContent(contentPath string, localPath string) (string, string) {
	var content, readmeContent string
	entries := readDir("./content/" + localPath)
	for _, entry := range entries {
		if entry.Name() == "README" {
			readmeContent = getMarkdown(readFile(contentPath + "/" + entry.Name()))
		} else if entry.Name() != "index.html" {
			content += formatNoteEntryHTML(entry, localPath)
		}
	}
	return content, readmeContent
}

func getRegularDirPageContent(contentPath string, localPath string) (string, string) {
	var content string
	entries := readDir("./content/" + localPath)
	for _, entry := range entries {
		if entry.Name() != "README" {
			file := readFile("./content/" + localPath + "/" + entry.Name())
			frontMatter, _ := parseFrontMatter(string(file))
			content += formatRegularEntryHTML(frontMatterToArticleEntry(frontMatter, "/"+localPath+"/"+removeDotMD(entry.Name())))
		}
	}
	return content, ""
}

func getPathLinks(path string, caser cases.Caser) string {
	var paths string
	pathLinks := "/"
	for _, file := range strings.Split(path, "/") {
		paths += "/" + file
		pathLinks += fmt.Sprintf("<a href=\"%s\" class=\"black\">%s</a>/", paths, caser.String(file))
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

func formatNoteEntryHTML(entry fs.DirEntry, localPath string) string {
	contentFilePath := removeDotMD(entry.Name())
	boldFormatted := contentFilePath
	if entry.IsDir() {
		boldFormatted = "<strong>" + boldFormatted + "</strong>"
	}
	content := fmt.Sprintf("<p><a href=\"/%s/%s\" class=\"black\">%s</a></p>\n", localPath, contentFilePath, boldFormatted)
	return content
}

func getTagsSting(tags []string) string {
	var tagsString string
	for _, tag := range tags {
		tagsString += fmt.Sprintf("<a href=\"/tags/%s\">%s</a>&nbsp;&nbsp;", strings.ToLower(tag), tag)
	}
	return tagsString
}

func formatRegularEntryHTML(entry articleEntry) string {
	tags := getTagsSting(entry.Tags)
	return fmt.Sprintf("<li><a href=\"%s\">%s</a> - <em>%s</em>&nbsp;&nbsp;<div class=\"tags\">%s</div></li>", entry.FilePath, entry.Title, entry.Date.Format("1/2/2006"), tags)
}

func formatTagsHTML(tags []string) string {
	var tagsString string
	for _, tag := range tags {
		tagsString += fmt.Sprintf("<a href=\"%s\">%s</a>&nbsp;", "/tags/"+strings.ToLower(tag), tag)
	}
	return tagsString
}

func getNotesDirPageContent(dirPath string) (string, string) {
	var readmeContent string
	if checkExists(dirPath + "/README") {
		readmeContent = getMarkdown(readFile(dirPath + "/README"))
	}
	return getFileStructureHTML(dirPath, 1), readmeContent
}

func getFileStructureHTML(dirPath string, depth int) string {
	var content string
	noteDir := readDir(dirPath)
	arrowPadding := strconv.Itoa(((depth - 1) * 7) + 15)
	normalPadding := strconv.Itoa(((depth) * 7) + 15)
	for _, item := range noteDir {
		if item.IsDir() {
			content += "<div><div class=\"dirName\"><i class=\"arrow right\" onclick=\"arrowClick(this)\"></i></div>&nbsp;<strong><a href=\"\">"
			content += item.Name()
			content += fmt.Sprintf("</a></strong><div style=\"padding-left: %spx;\" hidden>", normalPadding)
			innerContent := getFileStructureHTML(dirPath + "/" + item.Name(), depth + 1)
			content += innerContent
			content += "</div></div>"
		} else if item.Name() != "README" {
			url := "/" + strings.Join(splitBySlash(dirPath)[2:], "/") + "/" + removeDotMD(item.Name())
			content += fmt.Sprintf("<a href=\"%s\" style=\"padding-left: %spx;\">%s</a><br>", url, arrowPadding, removeDotMD(item.Name()))
		}
	}
	return content
}

func postOrderTraversal(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		check(err)
		buildDirPath := getBuildDirPath(path)
		localPath := getLocalPath(buildDirPath)
		pathSlice := splitBySlash(buildDirPath)
		lastPathSliceElement := pathSlice[len(pathSlice)-1]
		topDirectory := splitBySlash(localPath)[0]

		if info.IsDir() {
			if path == root {
				return nil
			}

			if !checkExists(buildDirPath) {
				createDir(buildDirPath)
			}

			caser := cases.Title(language.English)

			var content, readmeContent string
			if topDirectory == "blog" || topDirectory == "projects" {
				content, readmeContent = getRegularDirPageContent("./content/"+localPath, localPath)
			} else {
				if localPath == "notes" {
					content, readmeContent = getNotesDirPageContent("./content/notes")
					dirIndexHTML := openOrCreateFile(buildDirPath + "/index.html")
					pathLinks := getPathLinks(localPath, caser)
					currentPage := noteDirectoryPage{Content: content, Path: pathLinks, ReadMe: readmeContent}
					writeFileWithTemplate(dirIndexHTML, "./templates/note_dir_temp.html", currentPage)
					return nil
				} else {
					content, readmeContent = getNotesPageContent("./content/"+localPath, localPath)
				}
			}

			dirIndexHTML := openOrCreateFile(buildDirPath + "/index.html")
			pathLinks := getPathLinks(localPath, caser)
			currentPage := directoryPage{Title: caser.String(lastPathSliceElement), Path: pathLinks, Content: content, ReadMe: readmeContent}

			writeFileWithTemplate(dirIndexHTML, "./templates/dir_temp.html", currentPage)

			return nil
		}

		if lastPathSliceElement == "README" {
			return nil
		}

		file := readFile(path)

		if topDirectory == "blog" || topDirectory == "projects" {
			frontMatter, frontMatterEnd := parseFrontMatter(string(file))
			content := getMarkdown(file[frontMatterEnd:])
			fileOut := openOrCreateFile(removeDotMD(buildDirPath))
			tags := formatTagsHTML(frontMatter.Tags)
			articlePage := articlePage{Title: frontMatter.Title, Tags: tags, Date: frontMatter.Date, Content: content, Latex: frontMatter.Latex, Code: frontMatter.Code}
			writeFileWithTemplate(fileOut, "./templates/art_temp.html", articlePage)
		} else {
			content := getMarkdown(file)
			fileOut := openOrCreateFile(removeDotMD(buildDirPath))
			notePage := notePage{Title: "Title", Content: content, FilePath: "/something"}
			writeFileWithTemplate(fileOut, "./templates/note_temp.html", notePage)
		}

		return nil
	})
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

func writeFile(file *os.File, content []byte) {
	_, err := file.Write(content)
	check(err)
}

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
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

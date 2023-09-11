package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

func main() {
	post_order_traversal("./content")
	update_index()

	os.Chdir("./build")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	http.ListenAndServe(":3000", nil)
}

type art_entry struct {
	FilePath string
	Date     time.Time
	Title    string
	Tags     []string
}

type index_page struct {
	RecentPosts    string
	RecentProjects string
}

type art_page struct {
	Title   string
	Date    string
	Content string
	Latex   bool
	Code    bool
}

type dir_page struct {
	Title   string
	Path    string
	Content string
	ReadMe  string
}

type FrontMatter struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
	Latex bool     `yaml:"latex"`
	Code  bool     `yaml:"code"`
}

func remove_md(str string) string {
	return strings.Replace(str, ".md", "", 1)
}

func addArtEntryInOrder(slice []art_entry, path string, title string, tags []string, date time.Time) []art_entry {
	index := sort.Search(len(slice), func(i int) bool {
		return !slice[i].Date.Before(date)
	})

	slice = append(slice, art_entry{})
	copy(slice[index+1:], slice[index:])
	slice[index].Date = date
	slice[index].FilePath = path
	slice[index].Title = title
	slice[index].Tags = tags

	return slice
}

func add_limit_convert_html(recent_blogs []art_entry, limit int) string {
	blogs_content := ""
	for i, e := range recent_blogs {
		if i == limit {
			break
		}
		blogs_content += "<li><a href=\"" + remove_md(e.FilePath) + "\"><strong>" + e.Title + "</strong></a> - " + e.Date.Format("1/2/2006") + "</li>"
	}
	return blogs_content
}

func update_index() {
	recent_blogs := get_ordered_dates("blog")
	recent_projects := get_ordered_dates("projects")
	limit := 5
	blogs_content := add_limit_convert_html(recent_blogs, limit)
	projects_content := add_limit_convert_html(recent_projects, limit)

	file_out, err := os.OpenFile("./build/index.html", os.O_RDWR|os.O_CREATE, 0755)
	check(err)

	home := index_page{RecentPosts: blogs_content, RecentProjects: projects_content}

	tmpl, err := template.ParseFiles("./templates/index_temp.html")
	check(err)

	tmpl.Execute(file_out, home)
	check(err)
}

func get_ordered_dates(folder string) []art_entry {
	blog_posts, err := os.ReadDir("./content/" + folder)
	check(err)
	var date_slice []art_entry
	for _, e := range blog_posts {
		if e.Name() != "README" {
			file, err := os.Open("./content/" + folder + "/" + e.Name())
			check(err)
			contents, err := io.ReadAll(file)
			check(err)
			fm, _, err := parseFrontMatter(string(contents))
			check(err)
			date, err := time.Parse("1/2/2006", fm.Date)
			check(err)
			date_slice = addArtEntryInOrder(date_slice, "/"+folder+"/"+e.Name(), fm.Title, fm.Tags, date)
		}
	}
	return date_slice
}

func get_build_string(path string) string {
	build_path := strings.Split(path, "/")
	build_path[0] = "build"
	build_path[1] = strings.ToLower(build_path[1])
	build_string := "./" + strings.Join(build_path, "/")
	return build_string
}

func post_order_traversal(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		build_string := get_build_string(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			if _, err := os.Stat(build_string); err != nil {
				err := os.Mkdir(build_string, 0777)
				check(err)
			}
			path_slice := strings.Split(build_string, "/")
			path := strings.Join(path_slice[2:], "/")
			if path == "content" {
				return nil
			}
			content := ""

			currentWd, err := os.Getwd()
			check(err)

			entries, err := os.ReadDir("./content/" + path)
			check(err)

			readme := ""
			for _, e := range entries {
				if e.Name() == "README" {
					file, err := os.Open("./content/" + path + "/" + e.Name())
					check(err)
					contents, err := io.ReadAll(file)
					check(err)
					err = file.Close()
					check(err)
					readme = string(markdown.ToHTML(contents, nil, nil))
				} else if e.Name() != "index.html" {
					html_file := remove_md(e.Name())
					strong_if_dir := html_file
					if e.IsDir() {
						strong_if_dir = "<strong>" + html_file + "</strong>"
					}
					content += "<p><a href=\"/" + path + "/" + html_file + "\" class=\"black\">" + strong_if_dir + "</a></p>\n"
				}
			}

			os.Chdir(currentWd)

			dir_index, err := os.OpenFile(build_string+"/index.html", os.O_RDWR|os.O_CREATE, 0755)
			check(err)

			caser := cases.Title(language.English)
			paths := ""
			path_links := ""
			for _, e := range strings.Split(path, "/") {
				paths += "/" + e
				fmt.Println(paths, e)
				path_links += "<a href=\"" + paths + "\" class=\"black\">" + caser.String(e) + "</a>/"
			}
			current_page := dir_page{Title: caser.String(path_slice[len(path_slice)-1]), Path: "/" + path_links, Content: content, ReadMe: readme}

			tmpl, err := template.ParseFiles("./templates/dir_temp.html")
			check(err)

			tmpl.Execute(dir_index, current_page)
			check(err)

			if path == root {
				return nil
			}
			return nil
		}
		build_path := strings.Split(path, "/")
		if build_path[len(build_path)-1] == "README" {
			return nil
		}
		file, err := os.Open(path)
		check(err)

		contents, err := io.ReadAll(file)
		check(err)

		err = file.Close()
		check(err)

		fm, fmEnd, err := parseFrontMatter(string(contents))
		check(err)

		html := markdown.ToHTML(contents[fmEnd:], nil, nil)
		file_out, err := os.OpenFile(build_string, os.O_RDWR|os.O_CREATE, 0755)
		check(err)

		art_page := art_page{Title: fm.Title, Date: fm.Date, Content: string(html), Latex: fm.Latex, Code: fm.Code}

		tmpl, err := template.ParseFiles("./templates/art_temp.html")
		check(err)

		tmpl.Execute(file_out, art_page)
		check(err)

		return nil
	})
}

func parseFrontMatter(content string) (FrontMatter, int, error) {
	var fm FrontMatter
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
		if err := yaml.Unmarshal([]byte(content[3:frontMatterEnd-3]), &fm); err != nil {
			return fm, -1, err
		}
	}
	return fm, frontMatterEnd, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

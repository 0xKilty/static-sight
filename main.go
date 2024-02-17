package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func host_site(port string) {
	os.Chdir("./build")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	fmt.Println("Serving site on localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func open_file(path string) *os.File {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	return file
}

func read_file(path string) []byte {
	file, err := os.Open(path)
	check(err)
	contents, err := io.ReadAll(file)
	check(err)
	err = file.Close()
	check(err)
	return contents
}

func write_file(file *os.File, content []byte) {
	_, err := file.Write(content)
	check(err)
}

func clear_dir(path string) {
	err := os.RemoveAll("./build")
	check(err)
	create_dir("./build")
}

func create_dir(name string) {
	err := os.Mkdir(name, 0777)
	check(err)
}

func read_dir(path string) []fs.DirEntry {
	files, err := os.ReadDir(path)
	check(err)
	return files
}

func init_dir(to string) {
	create_dir(to)
	os.Chdir(to)
}

func copy_file(filename string, to string) {
	file_contents := read_file(filename)
	file_out := open_file(to + "/" + filename)
	write_file(file_out, file_contents)
}

func copy_assets() {
	os.Chdir("./assets")
	files := read_dir(".")
	for _, file := range files {
		copy_file(file.Name(), "../build")
	}
	os.Chdir("..")
}

func writer_with_template(writer io.Writer, template_path string, data interface{}) {
	tmpl, err := template.ParseFiles(template_path)
	check(err)
	tmpl.Execute(writer, data)
	check(err)
}

func remove_suffix(str string, to_remove string) string {
	return strings.TrimSuffix(str, ".md")
}

func md_to_url(md_file_name string) string {
	return 	dir_to_url(remove_suffix(md_file_name, ".md")) + ".html"
}

func dir_to_url(file_name string) string {
	return strings.ReplaceAll(strings.ToLower(file_name), " ", "-")
}

func get_HTML(bytes []byte) string {
	return string(markdown.ToHTML(bytes, nil, nil))
}

func string_to_date(str string) time.Time {
	date, err := time.Parse("1/2/2006", str)
	check(err)
	return date
}

func sort_posts_by_date(posts []page_template) {
	sort.Slice(posts, func(i, j int) bool {
		return string_to_date(posts[i].Date).Before(string_to_date(posts[j].Date))
	})
}

type front_matter_template struct {
	Date  string   `yaml:"date"`
	Tags  []string `yaml:"tags"`
	Latex bool     `yaml:"latex"`
	Code  bool     `yaml:"code"`
	Draft  bool     `yaml:"draft"`
}

func parseFrontMatter(content string) (front_matter_template, int) {
	var front_matter front_matter_template
	front_matter_end := 0
	if content[0:3] == "---" {
		front_matter_end = 3
		for i := 3; i < len(content); i++ {
			if content[i:i+3] == "---" {
				front_matter_end = i + 3
				break
			}
		}
	}
	if front_matter_end > 0 {
		if err := yaml.Unmarshal([]byte(content[3:front_matter_end-3]), &front_matter); err != nil {
			return front_matter, -1
		}
	}
	return front_matter, front_matter_end
}

type page_template struct {
	Title   string
	URL     string
	Date    string
	Content string
}

type index_template struct {
	Posts string
}

func generate_blog_and_index() {
	init_dir("./build/blog")
	defer os.Chdir("../..")

	var posts []page_template
	file_path := os.Args[1] + "/blog/"
	files := read_dir(file_path)
	for _, file := range files {
		input_file := read_file(file_path + file.Name())
		front_matter, front_matter_end := parseFrontMatter(string(input_file))
		if front_matter.Draft {
			continue
		}
		blog_structure := page_template{
			Title:   remove_suffix(file.Name(), ".md"),
			URL:     "/blog/" + md_to_url(file.Name()),
			Date:    front_matter.Date,
			Content: get_HTML(input_file[front_matter_end:]),
		}
		posts = append(posts, blog_structure)
		var buffer bytes.Buffer
		writer_with_template(&buffer, "../../templates/page.html", blog_structure)
		output_file := open_file(md_to_url(file.Name()))
		base := base_template{
			Content: buffer.String(),
			Code:    front_matter.Code,
			Latex:   front_matter.Latex,
		}
		writer_with_template(output_file, "../../templates/base.html", base)
	}

	sort_posts_by_date(posts)
	var posts_buffer bytes.Buffer
	for _, post := range posts {
		writer_with_template(&posts_buffer, "../../templates/page_bullet.html", post)
	}

	index_file := open_file("../index.html")
	var buffer bytes.Buffer
	writer_with_template(&buffer, "../../templates/index.html", index_template{Posts: posts_buffer.String()})
	writer_with_template(index_file, "../../templates/base.html", base_template{Content: buffer.String()})
}

type dir_template struct {
	Title string
	URL string
	Sub string
}

type base_template struct {
	Content string
	Code    bool
	Latex   bool
}

func resources_recursive(file_path string) bytes.Buffer {
	var resources bytes.Buffer
	files := read_dir(file_path)
	for _, file := range files {
		if file.IsDir() {
			create_dir(dir_to_url(file.Name()))
			sub := resources_recursive(file_path + file.Name())
			dir_bullet := dir_template{
				Title: file.Name(),
				Sub: sub.String(),
			}
			writer_with_template(&resources, "../../templates/dir_bullet.html", dir_bullet)
		} else {
			input_file := read_file(file_path + "/" + file.Name())
			page_structure := page_template{
				Title:   remove_suffix(file.Name(), ".md"),
				Content: get_HTML(input_file),
				URL: md_to_url(file_path[len(os.Args[1]):] + "/" + file.Name()),
			}

			writer_with_template(&resources, "../../templates/page_bullet.html", page_structure)

			var buffer bytes.Buffer
			writer_with_template(&buffer, "../../templates/page.html", page_structure)
			output_file := open_file(".." + md_to_url(file_path[len(os.Args[1]):] + "/" + file.Name()))
			base := base_template{
				Content: buffer.String(),
			}
			writer_with_template(output_file, "../../templates/base.html", base)
		}
	}
	return resources
}

func generate_resources() {
	init_dir("./build/resources")
	defer os.Chdir("../..")

	file_path := os.Args[1] + "/resources/"
	index_content := resources_recursive(file_path)

	var buffer bytes.Buffer
	index_file := open_file("index.html")
	writer_with_template(&buffer, "../../templates/resources_index.html", index_template{Posts: index_content.String()})
	writer_with_template(index_file, "../../templates/base.html", base_template{Content: buffer.String()})
}

func main() {
	if len(os.Args) == 2 {
		clear_dir("./build")
		copy_assets()
		generate_blog_and_index()
		generate_resources()
		host_site("3000")
	}
}

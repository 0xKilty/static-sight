package main

import (
	"sort"
	"time"
)

type tagPage struct {
	Title   string
	AllTags string
	Content string
}

type indexPage struct {
	RecentPosts    string
	RecentProjects string
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
		blogContent += formatRegularEntryHTML(entry)
	}
	return blogContent
}

func getAllTags(tagsMap map[string][]articleEntry) string {
	var allTags string
	for tagName := range tagsMap {
		allTags += "&nbsp;<a href=\"/tags/" + slugify(removeDotMD(tagName)) + "\">" + tagName + "</a>"
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
	writeFileWithTemplate(fileOut, "index_temp.html", index)

	if !checkExists("./build/tags") {
		createDir("./build/tags")
	}

	allTags := getAllTags(tagsMap)

	for tagName, tagEntries := range tagsMap {
		tagPageContent := getTagPageContent(tagEntries)
		tagHTMLFile := openOrCreateFile("./build/tags/" + slugify(tagName))
		tagPageData := tagPage{Title: tagName, AllTags: allTags, Content: tagPageContent}
		writeFileWithTemplate(tagHTMLFile, "tags_temp.html", tagPageData)
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
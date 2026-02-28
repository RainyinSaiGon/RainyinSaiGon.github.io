package parser

import (
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"portfolio/internal/model"
)

var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

// readTime estimates minutes to read based on a ~200 wpm average.
func readTime(htmlContent string) int {
	text := htmlTagRe.ReplaceAllString(htmlContent, " ")
	words := len(strings.Fields(text))
	minutes := int(math.Ceil(float64(words) / 200.0))
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}

// ReadPosts reads all .md files from dir and returns a slice of Posts.
func ReadPosts(dir string) ([]model.Post, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var posts []model.Post
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		slug := strings.TrimSuffix(f.Name(), ".md")
		posts = append(posts, parsePost(slug, string(raw)))
	}
	return posts, nil
}

// ReadProjects reads all .md files from dir and returns a slice of Projects.
// Returns an empty slice (no error) if the directory does not exist.
func ReadProjects(dir string) ([]model.Project, error) {
	files, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var projects []model.Project
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		slug := strings.TrimSuffix(f.Name(), ".md")
		projects = append(projects, parseProject(slug, string(raw)))
	}
	return projects, nil
}

// parsePost parses a markdown file with a simple key: value frontmatter block
// terminated by "---", followed by raw HTML content.
//
// Example:
//
//	title: My Post
//	date: 2026-02-28
//	description: A short summary
//	---
//	<p>HTML content here…</p>
func parsePost(slug, raw string) model.Post {
	post := model.Post{Slug: slug}
	lines := strings.Split(raw, "\n")
	bodyStart := len(lines)

	for i, line := range lines {
		if line == "---" {
			bodyStart = i + 1
			break
		}
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		switch key {
		case "title":
			post.Title = val
		case "date":
			post.Date = val
			if t, err := time.Parse("2006-01-02", val); err == nil {
				post.DateParsed = t
			}
		case "description":
			post.Description = val
		}
	}

	body := strings.Join(lines[bodyStart:], "\n")
	post.Content = template.HTML(body)
	post.ReadTime = readTime(body)
	return post
}

// parseProject parses a project markdown file (no body content, only frontmatter).
func parseProject(slug, raw string) model.Project {
	project := model.Project{Slug: slug}

	for _, line := range strings.Split(raw, "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		switch key {
		case "title":
			project.Title = val
		case "description":
			project.Description = val
		case "image":
			project.Image = val
		case "code":
			project.CodeURL = val
		case "demo":
			project.DemoURL = val
		case "featured":
			project.Featured = val == "true"
		}
	}
	return project
}

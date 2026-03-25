package parser

import (
	"bytes"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"portfolio/internal/model"

	"github.com/yuin/goldmark"
)

var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

// markdownToHTML converts markdown content to HTML using goldmark.
func markdownToHTML(markdown string) string {
	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return markdown // fallback to original if conversion fails
	}
	return buf.String()
}

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
	var posts []model.Post
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		slug := strings.TrimSuffix(filepath.Base(rel), ".md")
		relDir := filepath.Dir(rel)
		if relDir == "." {
			relDir = ""
		}
		relDir = filepath.ToSlash(relDir)

		posts = append(posts, parsePost(relDir, slug, string(raw)))
		return nil
	})
	if err != nil {
		return nil, err
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
func parsePost(path, slug, raw string) model.Post {
	post := model.Post{Slug: slug, Path: path}
	lines := strings.Split(raw, "\n")
	bodyStart := len(lines)

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "---" {
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
			if t, err := time.Parse("2006-01-02", val); err == nil {
				post.DateParsed = t
				post.Date = t.Format("Jan 2, 2006")
			} else {
				post.Date = val
			}
		case "description":
			post.Description = val
		case "tags":
			for _, tag := range strings.Split(val, ",") {
				if t := strings.TrimSpace(tag); t != "" {
					post.Tags = append(post.Tags, t)
				}
			}
		case "series":
			post.SeriesTag = val
		case "series_title":
			post.SeriesTitle = val
		}
	}

	body := strings.Join(lines[bodyStart:], "\n")
	htmlContent := markdownToHTML(body)
	post.Content = template.HTML(htmlContent)
	post.ReadTime = readTime(htmlContent)
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

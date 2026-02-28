package model

import (
	"html/template"
	"time"
)

// Post represents a blog post parsed from a markdown file.
type Post struct {
	Title       string
	Slug        string
	Date        string
	DateParsed  time.Time
	Description string
	ReadTime    int           // estimated minutes to read
	Content     template.HTML // raw HTML, not escaped in templates
}

// Project represents a portfolio project parsed from a markdown file.
type Project struct {
	Title       string
	Slug        string
	Description string
	Image       string
	CodeURL     string
	DemoURL     string
	Featured    bool
}

// HomeData holds the data passed to the home page template.
type HomeData struct {
	Posts    []Post
	Projects []Project
}

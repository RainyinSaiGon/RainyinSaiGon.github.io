package model

import (
	"html/template"
	"strings"
	"time"
)

// Post represents a blog post parsed from a markdown file.
type Post struct {
	Title       string
	Slug        string
	Path        string
	Date        string
	DateParsed  time.Time
	Description string
	Tags        []string      // topic tags e.g. ["Go", "AWS"]
	SeriesTag   string        // tag identifying the series
	SeriesTitle string        // title of this part in the series
	Series      *Series       // back-reference to the full series object
	SeriesPart  int           // 1-based index in the series
	SeriesNext  *Post         // next part in the series
	SeriesPrev  *Post         // previous part in the series
	ReadTime    int           // estimated minutes to read
	Content     template.HTML // raw HTML, not escaped in templates
}

// URLPath returns the canonical blog URL path for this post.
func (p Post) URLPath() string {
	if p.Path == "" {
		return "/blog/" + p.Slug + "/"
	}
	return "/blog/" + strings.Trim(p.Path, "/") + "/" + p.Slug + "/"
}

// Series represents a collection of related blog posts.
type Series struct {
	Tag   string
	Posts []*Post
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

// PostPageData holds a post and its surrounding posts for prev/next navigation.
type PostPageData struct {
	Post
	Prev *Post
	Next *Post
}

// SearchEntry is a lightweight record written into search.json for browser-side search.
type SearchEntry struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
	ReadTime    int      `json:"readTime"`
}

package renderer

import (
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"portfolio/internal/model"
)

const siteURL = "https://rainyinsaigon.github.io"

//go:embed templates
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

// Renderer renders HTML pages using embedded templates.
type Renderer struct {
	outputDir string
	tmpl      *template.Template
}

// New creates a Renderer that writes pages to outputDir.
// Templates are parsed from the embedded templates/ directory.
func New(outputDir string) (*Renderer, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Renderer{outputDir: outputDir, tmpl: tmpl}, nil
}

// CopyStaticFiles copies all files in static/ into the output directory.
func (r *Renderer) CopyStaticFiles() error {
	return fs.WalkDir(staticFS, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel("static", path)
		dst := filepath.Join(r.outputDir, rel)

		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		src, err := staticFS.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, src)
		return err
	})
}

// RenderHome renders the site home page.
func (r *Renderer) RenderHome(posts []model.Post, projects []model.Project) error {
	// Limit to 4 recent posts and 3 featured projects on the home page
	recent := posts
	if len(recent) > 4 {
		recent = recent[:4]
	}
	featured := make([]model.Project, 0, 3)
	for _, p := range projects {
		if p.Featured && len(featured) < 3 {
			featured = append(featured, p)
		}
	}

	data := model.HomeData{Posts: recent, Projects: featured}
	return r.write(filepath.Join(r.outputDir, "index.html"), "home", data)
}

// RenderBlogList renders the /blog index page.
func (r *Renderer) RenderBlogList(posts []model.Post) error {
	data := struct{ Posts []model.Post }{Posts: posts}
	return r.write(filepath.Join(r.outputDir, "blog", "index.html"), "blog_list", data)
}

// RenderPost renders an individual blog post page to /blog/<slug>/index.html.
// idx is the post's position in the sorted posts slice so prev/next can be computed.
func (r *Renderer) RenderPost(posts []model.Post, idx int) error {
	post := posts[idx]
	data := model.PostPageData{Post: post}
	if idx+1 < len(posts) {
		next := posts[idx+1]
		data.Next = &next
	}
	if idx > 0 {
		prev := posts[idx-1]
		data.Prev = &prev
	}
	dir := filepath.Join(r.outputDir, "blog", post.Slug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return r.write(filepath.Join(dir, "index.html"), "blog_post", data)
}

// RenderWorks renders the /works page.
func (r *Renderer) RenderWorks(projects []model.Project) error {
	data := struct{ Projects []model.Project }{Projects: projects}
	return r.write(filepath.Join(r.outputDir, "works", "index.html"), "works", data)
}

// RenderAbout renders the /about page.
func (r *Renderer) RenderAbout() error {
	return r.write(filepath.Join(r.outputDir, "about", "index.html"), "about", nil)
}

// Render404 renders a custom 404 error page.
func (r *Renderer) Render404() error {
	return r.write(filepath.Join(r.outputDir, "404.html"), "notfound", nil)
}

// RenderSearch renders the /search page.
func (r *Renderer) RenderSearch() error {
	return r.write(filepath.Join(r.outputDir, "search", "index.html"), "search", nil)
}

// GenerateSearchJSON writes docs/search.json for browser-side Fuse.js search.
func (r *Renderer) GenerateSearchJSON(posts []model.Post) error {
	entries := make([]model.SearchEntry, len(posts))
	for i, p := range posts {
		entries[i] = model.SearchEntry{
			Title:       p.Title,
			Slug:        p.Slug,
			Description: p.Description,
			Date:        p.Date,
			Tags:        p.Tags,
			ReadTime:    p.ReadTime,
		}
	}
	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	dst := filepath.Join(r.outputDir, "search.json")
	return os.WriteFile(dst, b, 0644)
}

// GenerateRSS writes docs/rss.xml as an RSS 2.0 feed.
func (r *Renderer) GenerateRSS(posts []model.Post) error {
	type Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		PubDate     string `xml:"pubDate"`
		Description string `xml:"description"`
	}
	type Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Items       []Item `xml:"item"`
	}
	type RSS struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		Channel Channel  `xml:"channel"`
	}

	items := make([]Item, len(posts))
	for i, p := range posts {
		items[i] = Item{
			Title:       p.Title,
			Link:        fmt.Sprintf("%s/blog/%s/", siteURL, p.Slug),
			PubDate:     p.DateParsed.UTC().Format(time.RFC1123Z),
			Description: p.Description,
		}
	}

	feed := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "RainyinSaiGon",
			Link:        siteURL,
			Description: "Software engineering, cloud, and explainable AI — by RainyinSaiGon",
			Language:    "en",
			Items:       items,
		},
	}

	out, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	content := append([]byte(xml.Header), out...)
	return os.WriteFile(filepath.Join(r.outputDir, "rss.xml"), content, 0644)
}

// GenerateSitemap writes docs/sitemap.xml.
func (r *Renderer) GenerateSitemap(posts []model.Post, projects []model.Project) error {
	type URL struct {
		Loc        string `xml:"loc"`
		LastMod    string `xml:"lastmod,omitempty"`
		ChangeFreq string `xml:"changefreq,omitempty"`
		Priority   string `xml:"priority,omitempty"`
	}
	type URLSet struct {
		XMLName xml.Name `xml:"urlset"`
		XMLNS   string   `xml:"xmlns,attr"`
		URLs    []URL    `xml:"url"`
	}

	today := time.Now().UTC().Format("2006-01-02")
	urls := []URL{
		{Loc: siteURL + "/", ChangeFreq: "weekly", Priority: "1.0", LastMod: today},
		{Loc: siteURL + "/blog/", ChangeFreq: "weekly", Priority: "0.9", LastMod: today},
		{Loc: siteURL + "/works/", ChangeFreq: "monthly", Priority: "0.8"},
		{Loc: siteURL + "/about/", ChangeFreq: "monthly", Priority: "0.7"},
		{Loc: siteURL + "/search/", ChangeFreq: "monthly", Priority: "0.5"},
	}
	for _, p := range posts {
		urls = append(urls, URL{
			Loc:        fmt.Sprintf("%s/blog/%s/", siteURL, p.Slug),
			LastMod:    p.DateParsed.UTC().Format("2006-01-02"),
			ChangeFreq: "yearly",
			Priority:   "0.7",
		})
	}

	set := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
	out, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}
	content := append([]byte(xml.Header), out...)
	return os.WriteFile(filepath.Join(r.outputDir, "sitemap.xml"), content, 0644)
}

// write creates all necessary directories and executes the named template into path.
func (r *Renderer) write(path, tmplName string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return r.tmpl.ExecuteTemplate(f, tmplName, data)
}

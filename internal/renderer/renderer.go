package renderer

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"portfolio/internal/model"
)

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
func (r *Renderer) RenderPost(post model.Post) error {
	dir := filepath.Join(r.outputDir, "blog", post.Slug)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return r.write(filepath.Join(dir, "index.html"), "blog_post", post)
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

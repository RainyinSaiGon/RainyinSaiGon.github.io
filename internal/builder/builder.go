package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"portfolio/internal/parser"
	"portfolio/internal/renderer"
)

// Config holds the input/output directories for the build.
type Config struct {
	ContentDir string // e.g. "content"
	OutputDir  string // e.g. "docs"
}

// Build parses all content, sorts it, and renders the full site.
func Build(cfg Config) error {
	// Parse content
	posts, err := parser.ReadPosts(filepath.Join(cfg.ContentDir, "posts"))
	if err != nil {
		return fmt.Errorf("reading posts: %w", err)
	}
	projects, err := parser.ReadProjects(filepath.Join(cfg.ContentDir, "projects"))
	if err != nil {
		return fmt.Errorf("reading projects: %w", err)
	}

	// Sort posts newest-first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].DateParsed.After(posts[j].DateParsed)
	})

	// Prepare output directory
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	r, err := renderer.New(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("initialising renderer: %w", err)
	}

	// Copy static assets (style.css, etc.)
	if err := r.CopyStaticFiles(); err != nil {
		return fmt.Errorf("copying static files: %w", err)
	}

	// Render pages
	if err := r.RenderHome(posts, projects); err != nil {
		return fmt.Errorf("rendering home: %w", err)
	}
	if err := r.RenderBlogList(posts); err != nil {
		return fmt.Errorf("rendering blog list: %w", err)
	}
	for i := range posts {
		if err := r.RenderPost(posts, i); err != nil {
			return fmt.Errorf("rendering post %s: %w", posts[i].Slug, err)
		}
	}
	if err := r.RenderWorks(projects); err != nil {
		return fmt.Errorf("rendering works: %w", err)
	}
	if err := r.RenderAbout(); err != nil {
		return fmt.Errorf("rendering about: %w", err)
	}
	if err := r.Render404(); err != nil {
		return fmt.Errorf("rendering 404: %w", err)
	}
	if err := r.RenderSearch(); err != nil {
		return fmt.Errorf("rendering search: %w", err)
	}
	if err := r.GenerateSearchJSON(posts); err != nil {
		return fmt.Errorf("generating search.json: %w", err)
	}
	if err := r.GenerateRSS(posts); err != nil {
		return fmt.Errorf("generating RSS: %w", err)
	}
	if err := r.GenerateSitemap(posts, projects); err != nil {
		return fmt.Errorf("generating sitemap: %w", err)
	}

	fmt.Printf("Built %d post(s), %d project(s) → %s/\n", len(posts), len(projects), cfg.OutputDir)
	return nil
}

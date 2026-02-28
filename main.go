package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

type Post struct {
	Title       string
	Slug        string
	Date        string
	DateParsed  time.Time
	Description string
	Content     string
}

type Project struct {
	Title       string
	Slug        string
	Description string
	Image       string
	CodeURL     string
	DemoURL     string
	Featured    bool
}

const outDir = "docs"

func main() {
	fmt.Println("Building portfolio website...")

	// Create docs directory (GitHub Pages serves from /docs)
	os.MkdirAll(outDir, 0755)

	// Read posts from content/posts/
	posts, err := readPosts("content/posts")
	if err != nil {
		fmt.Printf("Error reading posts: %v\n", err)
		return
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].DateParsed.After(posts[j].DateParsed)
	})

	// Read projects from content/projects/
	projects, err := readProjects("content/projects")
	if err != nil {
		fmt.Printf("Error reading projects: %v\n", err)
		return
	}

	// Generate pages
	if err := generateHomePage(posts, projects); err != nil {
		fmt.Printf("Error generating home page: %v\n", err)
		return
	}

	if err := generateBlogPage(posts); err != nil {
		fmt.Printf("Error generating blog page: %v\n", err)
		return
	}

	if err := generateWorksPage(projects); err != nil {
		fmt.Printf("Error generating works page: %v\n", err)
		return
	}

	if err := generateAboutPage(); err != nil {
		fmt.Printf("Error generating about page: %v\n", err)
		return
	}

	// Generate individual blog posts
	for _, post := range posts {
		if err := generatePost(post); err != nil {
			fmt.Printf("Error generating post %s: %v\n", post.Slug, err)
			return
		}
	}

	fmt.Printf("✓ Build complete! Generated %d posts and %d projects\n", len(posts), len(projects))
	fmt.Println("Push the 'docs' folder → GitHub Pages serves from /docs on main branch")
}

func readPosts(dir string) ([]Post, error) {
	var posts []Post

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			slug := strings.TrimSuffix(file.Name(), ".md")
			filePath := filepath.Join(dir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			post := parseMarkdown(slug, string(content))
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func readProjects(dir string) ([]Project, error) {
	var projects []Project

	files, err := os.ReadDir(dir)
	if err != nil {
		return projects, nil // Return empty if directory doesn't exist
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			slug := strings.TrimSuffix(file.Name(), ".md")
			filePath := filepath.Join(dir, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			project := parseProject(slug, string(content))
			projects = append(projects, project)
		}
	}

	return projects, nil
}

func parseMarkdown(slug string, content string) Post {
	lines := strings.Split(content, "\n")
	post := Post{Slug: slug}
	contentStart := 0

	for i, line := range lines {
		if strings.HasPrefix(line, "title:") {
			post.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
		} else if strings.HasPrefix(line, "date:") {
			dateStr := strings.TrimSpace(strings.TrimPrefix(line, "date:"))
			post.Date = dateStr
			// Parse date for sorting (expect YYYY-MM-DD format)
			if t, err := time.Parse("2006-01-02", dateStr); err == nil {
				post.DateParsed = t
			}
		} else if strings.HasPrefix(line, "description:") {
			post.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		} else if line == "---" {
			contentStart = i + 1
			break
		}
	}

	if contentStart > 0 {
		post.Content = strings.Join(lines[contentStart:], "\n")
	}

	return post
}

func parseProject(slug string, content string) Project {
	lines := strings.Split(content, "\n")
	project := Project{Slug: slug}

	for _, line := range lines {
		if strings.HasPrefix(line, "title:") {
			project.Title = strings.TrimSpace(strings.TrimPrefix(line, "title:"))
		} else if strings.HasPrefix(line, "description:") {
			project.Description = strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		} else if strings.HasPrefix(line, "image:") {
			project.Image = strings.TrimSpace(strings.TrimPrefix(line, "image:"))
		} else if strings.HasPrefix(line, "code:") {
			project.CodeURL = strings.TrimSpace(strings.TrimPrefix(line, "code:"))
		} else if strings.HasPrefix(line, "demo:") {
			project.DemoURL = strings.TrimSpace(strings.TrimPrefix(line, "demo:"))
		} else if strings.HasPrefix(line, "featured:") {
			featured := strings.TrimSpace(strings.TrimPrefix(line, "featured:"))
			project.Featured = featured == "true"
		}
	}

	return project
}

func baseCSS() string {
	return `
		*, *::before, *::after { margin: 0; padding: 0; box-sizing: border-box; }
		:root {
			--blue: #1a6eb5;
			--blue-light: #3b9eff;
			--blue-bg: #ddeeff;
			--text: #1a1a2e;
			--muted: #6b7280;
			--border: #e5e7eb;
			--bg: #ffffff;
			--max-w: 1100px;
		}
		html { scroll-behavior: smooth; }
		body { font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; color: var(--text); background: var(--bg); line-height: 1.7; font-size: 16px; }
		a { color: var(--blue); text-decoration: none; }
		a:hover { text-decoration: underline; }
		/* NAV */
		.site-nav { display: flex; align-items: center; justify-content: space-between; padding: 20px 40px; max-width: var(--max-w); margin: 0 auto; }
		.site-nav .logo { font-weight: 700; font-size: 18px; color: var(--text); }
		.site-nav .logo:hover { text-decoration: none; color: var(--blue); }
		.site-nav .nav-links { display: flex; gap: 32px; }
		.site-nav .nav-links a { color: var(--text); font-weight: 500; font-size: 15px; }
		.site-nav .nav-links a:hover { color: var(--blue); text-decoration: none; }
		/* INNER PAGE NAV (smaller pages) */
		.inner-nav { border-bottom: 1px solid var(--border); margin-bottom: 0; }
		/* FOOTER */
		.site-footer { padding: 32px 40px; max-width: var(--max-w); margin: 0 auto; color: var(--muted); font-size: 14px; border-top: 1px solid var(--border); margin-top: 80px; }
		.site-footer a { color: var(--muted); text-decoration: underline; }
		/* INNER PAGES */
		.page-wrap { max-width: 760px; margin: 0 auto; padding: 60px 24px; }
		.page-wrap h1 { font-size: 36px; font-weight: 800; margin-bottom: 40px; }
		.back-link { display: inline-flex; align-items: center; gap: 6px; color: var(--muted); font-size: 14px; margin-bottom: 28px; }
		.back-link:hover { color: var(--blue); text-decoration: none; }
`
}

func generateHomePage(posts []Post, projects []Project) error {
	// Get recent posts (limit to 4)
	recentPosts := posts
	if len(recentPosts) > 4 {
		recentPosts = recentPosts[:4]
	}

	// Get featured projects (limit to 3)
	recentProjects := []Project{}
	for _, p := range projects {
		if p.Featured && len(recentProjects) < 3 {
			recentProjects = append(recentProjects, p)
		}
	}

	homeTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        /* HERO SECTION */
        .hero-bg { background: linear-gradient(135deg, #ddeeff 0%, #c8e4ff 60%, #b8d4f0 100%); }
        .hero-nav { max-width: var(--max-w); }
        .hero-body { display: flex; align-items: center; justify-content: space-between; max-width: var(--max-w); margin: 0 auto; padding: 60px 40px 80px; gap: 40px; }
        .hero-left { flex: 1; }
        .hero-left h1 { font-size: clamp(36px, 5vw, 56px); font-weight: 900; line-height: 1.1; margin-bottom: 20px; color: var(--text); }
        .hero-left p { font-size: 17px; color: #374151; max-width: 480px; margin-bottom: 12px; }
        .hero-cta { display: flex; gap: 16px; margin-top: 36px; flex-wrap: wrap; }
        .btn-primary { display: inline-flex; align-items: center; gap: 8px; background: var(--blue); color: #fff; padding: 13px 26px; border-radius: 8px; font-weight: 600; font-size: 15px; }
        .btn-primary:hover { background: #155a9a; text-decoration: none; }
        .btn-outline { display: inline-flex; align-items: center; gap: 8px; border: 2px solid var(--blue); color: var(--blue); padding: 11px 24px; border-radius: 8px; font-weight: 600; font-size: 15px; }
        .btn-outline:hover { background: var(--blue-bg); text-decoration: none; }
        /* TILE GRID */
        .hero-right { flex-shrink: 0; position: relative; width: 360px; height: 360px; }
        .blob { position: absolute; inset: -20px; background: radial-gradient(circle at 60% 50%, #a8cff0 0%, #b0d8f8 40%, transparent 70%); border-radius: 50%; opacity: 0.7; }
        .tiles { position: relative; display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; width: 280px; margin: 40px auto 0; }
        .tile { background: var(--blue-light); border-radius: 14px; aspect-ratio: 1; opacity: 0; animation: tile-in 0.5s ease forwards; }
        .tile:nth-child(1) { grid-column: 1; grid-row: 1; animation-delay: 0.0s; }
        .tile:nth-child(2) { grid-column: 2; grid-row: 1; animation-delay: 0.08s; }
        .tile:nth-child(3) { grid-column: 1; grid-row: 2; animation-delay: 0.16s; }
        .tile:nth-child(4) { grid-column: 2; grid-row: 2; animation-delay: 0.24s; }
        .tile:nth-child(5) { grid-column: 3; grid-row: 2; animation-delay: 0.32s; }
        .tile:nth-child(6) { grid-column: 1; grid-row: 3; animation-delay: 0.40s; }
        .tile:nth-child(7) { grid-column: 2; grid-row: 3; animation-delay: 0.48s; }
        .tile:nth-child(8) { grid-column: 3; grid-row: 3; animation-delay: 0.56s; }
        @keyframes tile-in { from { opacity: 0; transform: scale(0.7); } to { opacity: 1; transform: scale(1); } }
        /* WAVE DIVIDER */
        .wave { display: block; width: 100%; line-height: 0; }
        /* WORKS SECTION */
        .works-bg { background: #f0f7ff; padding: 60px 0; }
        .works-inner { max-width: var(--max-w); margin: 0 auto; padding: 0 40px; }
        .section-header { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 32px; }
        .section-header h2 { font-size: 28px; font-weight: 800; }
        .section-header a { font-size: 14px; color: var(--blue); display: flex; align-items: center; gap: 4px; }
        .works-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 24px; }
        .work-card { background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 4px rgba(0,0,0,0.08); display: flex; flex-direction: column; }
        .work-card-img { width: 100%; height: 180px; object-fit: cover; background: var(--blue-bg); }
        .work-card-body { padding: 20px; flex: 1; display: flex; flex-direction: column; }
        .work-card-title { font-size: 17px; font-weight: 700; margin-bottom: 8px; }
        .work-card-desc { color: var(--muted); font-size: 14px; line-height: 1.6; flex: 1; }
        .work-card-links { display: flex; gap: 12px; margin-top: 16px; }
        .work-card-btn { font-size: 13px; font-weight: 600; color: var(--blue); border: 1.5px solid var(--blue); padding: 5px 14px; border-radius: 6px; }
        .work-card-btn:hover { background: var(--blue-bg); text-decoration: none; }
        /* POSTS SECTION */
        .posts-bg { background: #fff; padding: 60px 0; }
        .posts-inner { max-width: var(--max-w); margin: 0 auto; padding: 0 40px; }
        .posts-list { list-style: none; }
        .posts-list li { display: flex; align-items: baseline; gap: 20px; padding: 16px 0; border-bottom: 1px solid var(--border); }
        .posts-list li:last-child { border-bottom: none; }
        .post-date { color: var(--muted); font-size: 14px; white-space: nowrap; min-width: 95px; }
        .post-link { font-weight: 600; font-size: 16px; color: var(--text); }
        .post-link:hover { color: var(--blue); text-decoration: none; }
        /* RESPONSIVE */
        @media (max-width: 700px) {
          .hero-body { flex-direction: column; padding: 40px 20px 60px; }
          .hero-right { width: 240px; height: 240px; }
          .tiles { width: 180px; }
          .site-nav { padding: 16px 20px; }
          .works-inner, .posts-inner { padding: 0 20px; }
        }
    </style>
</head>
<body>
    <div class="hero-bg">
        <nav class="site-nav">
            <a class="logo" href="/">RainyinSaiGon</a>
            <div class="nav-links">
                <a href="/works">Works</a>
                <a href="/blog">Blog</a>
                <a href="/about">About</a>
            </div>
        </nav>
        <div class="hero-body">
            <div class="hero-left">
                <h1>Hi, I'm RainyinSaiGon</h1>
                <p>I build fast, elegant, and performant software.</p>
                <p>This is where I share my projects and thoughts, mostly about code and technology.</p>
                <div class="hero-cta">
                    <a class="btn-primary" href="/blog">My blog ✍</a>
                    <a class="btn-outline" href="/works">My projects ⚙</a>
                </div>
            </div>
            <div class="hero-right">
                <div class="blob"></div>
                <div class="tiles">
                    <div class="tile"></div>
                    <div class="tile"></div>
                    <div class="tile"></div>
                    <div class="tile"></div>
                    <div class="tile"></div>
                    <div class="tile"></div>
                    <div class="tile"></div>
                    <div class="tile"></div>
                </div>
            </div>
        </div>
    </div>

    <svg class="wave" viewBox="0 0 1440 60" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"><path d="M0,0 C360,60 1080,60 1440,0 L1440,60 L0,60 Z" fill="#f0f7ff"/></svg>

    {{if .Projects}}
    <div class="works-bg">
        <div class="works-inner">
            <div class="section-header">
                <h2>Recent works</h2>
                <a href="/works">All works →</a>
            </div>
            <div class="works-grid">
                {{range .Projects}}
                <div class="work-card">
                    {{if .Image}}<img class="work-card-img" src="{{.Image}}" alt="{{.Title}}">{{else}}<div class="work-card-img"></div>{{end}}
                    <div class="work-card-body">
                        <div class="work-card-title">{{.Title}}</div>
                        <div class="work-card-desc">{{.Description}}</div>
                        <div class="work-card-links">
                            {{if .CodeURL}}<a class="work-card-btn" href="{{.CodeURL}}">Code</a>{{end}}
                            {{if .DemoURL}}<a class="work-card-btn" href="{{.DemoURL}}">Demo</a>{{end}}
                        </div>
                    </div>
                </div>
                {{end}}
            </div>
        </div>
    </div>
    {{end}}

    <svg class="wave" viewBox="0 0 1440 60" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none"><path d="M0,60 C360,0 1080,0 1440,60 L1440,0 L0,0 Z" fill="#f0f7ff"/></svg>

    {{if .Posts}}
    <div class="posts-bg">
        <div class="posts-inner">
            <div class="section-header">
                <h2>Latest posts</h2>
                <a href="/blog">All posts →</a>
            </div>
            <ul class="posts-list">
                {{range .Posts}}
                <li>
                    <span class="post-date">{{.Date}}</span>
                    <a class="post-link" href="/blog/{{.Slug}}">{{.Title}}</a>
                </li>
                {{end}}
            </ul>
        </div>
    </div>
    {{end}}

    <footer class="site-footer">
        <p>© 2026 RainyinSaiGon. Content licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("home").Parse(homeTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Posts    []Post
		Projects []Project
	}{
		Posts:    recentPosts,
		Projects: recentProjects,
	}

	file, err := os.Create(outDir + "/index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

func generateBlogPage(posts []Post) error {
	blogTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog — RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        .posts-list { list-style: none; }
        .posts-list li { display: flex; align-items: baseline; gap: 20px; padding: 18px 0; border-bottom: 1px solid var(--border); }
        .posts-list li:last-child { border-bottom: none; }
        .post-date { color: var(--muted); font-size: 14px; white-space: nowrap; min-width: 95px; }
        .post-link { font-weight: 600; font-size: 16px; color: var(--text); }
        .post-link:hover { color: var(--blue); text-decoration: none; }
        .post-desc { color: var(--muted); font-size: 14px; display: block; margin-top: 2px; font-weight: 400; }
    </style>
</head>
<body>
    <nav class="site-nav inner-nav">
        <a class="logo" href="/">RainyinSaiGon</a>
        <div class="nav-links">
            <a href="/works">Works</a>
            <a href="/blog">Blog</a>
            <a href="/about">About</a>
        </div>
    </nav>

    <div class="page-wrap">
        <h1>Blog</h1>
        {{if .}}
        <ul class="posts-list">
            {{range .}}
            <li>
                <span class="post-date">{{.Date}}</span>
                <div>
                    <a class="post-link" href="/blog/{{.Slug}}">{{.Title}}</a>
                    {{if .Description}}<span class="post-desc">{{.Description}}</span>{{end}}
                </div>
            </li>
            {{end}}
        </ul>
        {{else}}
        <p style="color:var(--muted)">No blog posts yet. Check back soon!</p>
        {{end}}
    </div>

    <footer class="site-footer">
        <p>© 2026 RainyinSaiGon. Content licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("blog").Parse(blogTemplate)
	if err != nil {
		return err
	}

	os.MkdirAll(outDir+"/blog", 0755)
	file, err := os.Create(outDir + "/blog/index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, posts)
}

func generateWorksPage(projects []Project) error {
	worksTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Works — RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        .works-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 24px; }
        .work-card { background: #fff; border: 1px solid var(--border); border-radius: 12px; overflow: hidden; display: flex; flex-direction: column; }
        .work-card-img { width: 100%; height: 200px; object-fit: cover; background: var(--blue-bg); }
        .work-card-body { padding: 20px; flex: 1; display: flex; flex-direction: column; }
        .work-card-title { font-size: 18px; font-weight: 700; margin-bottom: 8px; }
        .work-card-desc { color: var(--muted); font-size: 14px; line-height: 1.6; flex: 1; }
        .work-card-links { display: flex; gap: 12px; margin-top: 16px; }
        .work-card-btn { font-size: 13px; font-weight: 600; color: var(--blue); border: 1.5px solid var(--blue); padding: 5px 14px; border-radius: 6px; }
        .work-card-btn:hover { background: var(--blue-bg); text-decoration: none; }
    </style>
</head>
<body>
    <nav class="site-nav inner-nav">
        <a class="logo" href="/">RainyinSaiGon</a>
        <div class="nav-links">
            <a href="/works">Works</a>
            <a href="/blog">Blog</a>
            <a href="/about">About</a>
        </div>
    </nav>

    <div class="page-wrap" style="max-width:1100px">
        <h1>Works</h1>
        {{if .}}
        <div class="works-grid">
            {{range .}}
            <div class="work-card">
                {{if .Image}}<img class="work-card-img" src="{{.Image}}" alt="{{.Title}}">{{else}}<div class="work-card-img"></div>{{end}}
                <div class="work-card-body">
                    <div class="work-card-title">{{.Title}}</div>
                    <div class="work-card-desc">{{.Description}}</div>
                    <div class="work-card-links">
                        {{if .CodeURL}}<a class="work-card-btn" href="{{.CodeURL}}">Code</a>{{end}}
                        {{if .DemoURL}}<a class="work-card-btn" href="{{.DemoURL}}">Demo</a>{{end}}
                    </div>
                </div>
            </div>
            {{end}}
        </div>
        {{else}}
        <p style="color:var(--muted)">No projects yet. Check back soon!</p>
        {{end}}
    </div>

    <footer class="site-footer">
        <p>© 2026 RainyinSaiGon. Content licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("works").Parse(worksTemplate)
	if err != nil {
		return err
	}

	os.MkdirAll(outDir+"/works", 0755)
	file, err := os.Create(outDir + "/works/index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, projects)
}

func generateAboutPage() error {
	aboutTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>About — RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        .about-content h2 { font-size: 20px; font-weight: 700; margin: 36px 0 12px; }
        .about-content p { color: #374151; margin-bottom: 14px; line-height: 1.75; }
    </style>
</head>
<body>
    <nav class="site-nav inner-nav">
        <a class="logo" href="/">RainyinSaiGon</a>
        <div class="nav-links">
            <a href="/works">Works</a>
            <a href="/blog">Blog</a>
            <a href="/about">About</a>
        </div>
    </nav>

    <div class="page-wrap">
        <h1>About</h1>
        <div class="about-content">
            <p>Hi! I'm RainyinSaiGon — a developer interested in building fast, efficient, and elegant software.</p>

            <h2>What I work on</h2>
            <p>I enjoy working with Go, web technologies, and systems programming. I like understanding how things work at a fundamental level and building tools that solve real problems.</p>

            <h2>Projects</h2>
            <p>Take a look at my <a href="/works">works page</a> to see what I've been building.</p>

            <h2>Writing</h2>
            <p>I occasionally write about things I find interesting on my <a href="/blog">blog</a>.</p>

            <h2>Contact</h2>
            <p>Find me on <a href="https://github.com/RainyinSaiGon">GitHub</a>.</p>
        </div>
    </div>

    <footer class="site-footer">
        <p>© 2026 RainyinSaiGon. Content licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("about").Parse(aboutTemplate)
	if err != nil {
		return err
	}

	os.MkdirAll(outDir+"/about", 0755)
	file, err := os.Create(outDir + "/about/index.html")
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, nil)
}

func generatePost(post Post) error {
	postTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} — RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        .post-header { margin-bottom: 40px; }
        .post-header h1 { font-size: 36px; font-weight: 900; line-height: 1.2; margin-bottom: 12px; }
        .post-meta { color: var(--muted); font-size: 14px; }
        .post-body { color: #374151; line-height: 1.8; }
        .post-body h2 { font-size: 22px; font-weight: 700; margin: 36px 0 14px; color: var(--text); }
        .post-body p { margin-bottom: 16px; }
        .post-body ul, .post-body ol { margin-left: 24px; margin-bottom: 16px; }
        .post-body li { margin-bottom: 8px; }
        .post-body code { background: #f1f5f9; padding: 2px 7px; border-radius: 4px; font-family: 'Cascadia Code', 'JetBrains Mono', monospace; font-size: 0.9em; }
        .post-body pre { background: #1e293b; color: #e2e8f0; padding: 20px; border-radius: 10px; overflow-x: auto; margin-bottom: 20px; }
        .post-body pre code { background: none; color: inherit; padding: 0; }
        .post-body a { color: var(--blue); }
    </style>
</head>
<body>
    <nav class="site-nav inner-nav">
        <a class="logo" href="/">RainyinSaiGon</a>
        <div class="nav-links">
            <a href="/works">Works</a>
            <a href="/blog">Blog</a>
            <a href="/about">About</a>
        </div>
    </nav>

    <div class="page-wrap">
        <a href="/blog" class="back-link">← Blog</a>
        <div class="post-header">
            <h1>{{.Title}}</h1>
            <div class="post-meta">{{.Date}}</div>
        </div>
        <div class="post-body">
            {{.Content}}
        </div>
    </div>

    <footer class="site-footer">
        <p>© 2026 RainyinSaiGon. Content licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return err
	}

	// Ensure blog directory exists in docs
	os.MkdirAll(outDir+"/blog", 0755)
	postDir := filepath.Join(outDir, "blog", post.Slug)
	os.MkdirAll(postDir, 0755)

	file, err := os.Create(filepath.Join(postDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, post)
}

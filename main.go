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

func main() {
	fmt.Println("Building portfolio website...")

	// Create public directory if it doesn't exist
	os.MkdirAll("public", 0755)

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
	fmt.Println("Deploy the 'public' folder to GitHub Pages")
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
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; line-height: 1.6; color: #333; background: #fff; }
.container { max-width: 900px; margin: 0 auto; padding: 40px 20px; }
header { margin-bottom: 60px; border-bottom: 1px solid #eee; padding-bottom: 40px; }
h1 { font-size: 32px; margin-bottom: 10px; }
h2 { font-size: 24px; margin: 40px 0 20px 0; }
.subtitle { color: #666; font-size: 18px; margin-bottom: 10px; }
.tagline { color: #999; font-size: 16px; margin-bottom: 30px; }
.nav { display: flex; gap: 30px; margin: 20px 0; }
.nav a { text-decoration: none; color: #0066cc; }
.nav a:hover { text-decoration: underline; }
.back-link { color: #0066cc; text-decoration: none; margin-bottom: 20px; display: inline-block; }
.back-link:hover { text-decoration: underline; }
footer { margin-top: 60px; padding-top: 40px; border-top: 1px solid #eee; color: #999; font-size: 14px; }
footer a { color: #0066cc; text-decoration: none; }
footer a:hover { text-decoration: underline; }
`
}

func generateHomePage(posts []Post, projects []Project) error {
	// Get recent posts (limit to 3)
	recentPosts := posts
	if len(recentPosts) > 3 {
		recentPosts = recentPosts[:3]
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
    <title>RainyinSaiGon - Developer & Builder</title>
    <style>
        ` + baseCSS() + `
        .hero { margin-bottom: 60px; }
        .hero h1 { font-size: 48px; margin-bottom: 20px; }
        .section { margin-bottom: 60px; }
        .project-item { margin-bottom: 40px; padding-bottom: 40px; border-bottom: 1px solid #eee; }
        .project-item:last-child { border-bottom: none; }
        .project-image { width: 100%; height: 300px; background: #f5f5f5; margin-bottom: 20px; border-radius: 8px; }
        .project-title { font-size: 22px; margin-bottom: 10px; }
        .project-title a { color: #0066cc; text-decoration: none; }
        .project-title a:hover { text-decoration: underline; }
        .project-description { color: #666; margin-bottom: 15px; }
        .project-links { display: flex; gap: 20px; }
        .project-links a { color: #0066cc; text-decoration: none; }
        .project-links a:hover { text-decoration: underline; }
        .post-item { margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px solid #eee; }
        .post-item:last-child { border-bottom: none; }
        .post-title a { color: #0066cc; text-decoration: none; }
        .post-title a:hover { text-decoration: underline; }
        .post-meta { color: #999; font-size: 14px; }
        .view-all { margin-top: 20px; }
        .view-all a { color: #0066cc; text-decoration: none; }
        .view-all a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Hi, I'm RainyinSaiGon</h1>
            <p class="subtitle">Developer & Builder</p>
            <p class="tagline">Building fast, elegant, and performant solutions.</p>
            <nav class="nav">
                <a href="/">Home</a>
                <a href="/works">Works</a>
                <a href="/blog">Blog</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            {{if .Projects}}
            <section class="section">
                <h2>Recent works</h2>
                {{range .Projects}}
                <article class="project-item">
                    {{if .Image}}
                    <div class="project-image" style="background-image: url('{{.Image}}'); background-size: cover; background-position: center;"></div>
                    {{end}}
                    <h3 class="project-title">{{.Title}}</h3>
                    <p class="project-description">{{.Description}}</p>
                    <div class="project-links">
                        {{if .CodeURL}}<a href="{{.CodeURL}}">Code</a>{{end}}
                        {{if .DemoURL}}<a href="{{.DemoURL}}">Demo</a>{{end}}
                    </div>
                </article>
                {{end}}
                <div class="view-all"><a href="/works">All works →</a></div>
            </section>
            {{end}}

            {{if .Posts}}
            <section class="section">
                <h2>Latest posts</h2>
                {{range .Posts}}
                <article class="post-item">
                    <p class="post-meta">{{.Date}}</p>
                    <h3 class="post-title"><a href="/blog/{{.Slug}}">{{.Title}}</a></h3>
                </article>
                {{end}}
                <div class="view-all"><a href="/blog">All posts →</a></div>
            </section>
            {{end}}
        </main>

        <footer>
            <p>© 2026 RainyinSaiGon. Content on this site is licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
        </footer>
    </div>
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

	file, err := os.Create("public/index.html")
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
    <title>Blog - RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        .post-item { margin-bottom: 30px; padding-bottom: 30px; border-bottom: 1px solid #eee; }
        .post-item:last-child { border-bottom: none; }
        .post-title { font-size: 20px; margin-bottom: 8px; }
        .post-title a { color: #0066cc; text-decoration: none; }
        .post-title a:hover { text-decoration: underline; }
        .post-meta { color: #999; font-size: 14px; margin-bottom: 10px; }
        .post-description { color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Blog</h1>
            <nav class="nav">
                <a href="/">Home</a>
                <a href="/works">Works</a>
                <a href="/blog">Blog</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            {{if .}}
            <div class="posts">
                {{range .}}
                <article class="post-item">
                    <h3 class="post-title"><a href="/blog/{{.Slug}}">{{.Title}}</a></h3>
                    <p class="post-meta">{{.Date}}</p>
                    <p class="post-description">{{.Description}}</p>
                </article>
                {{end}}
            </div>
            {{else}}
            <p>No blog posts yet. Check back soon!</p>
            {{end}}
        </main>

        <footer>
            <p>© 2026 RainyinSaiGon. Content on this site is licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
        </footer>
    </div>
</body>
</html>`

	tmpl, err := template.New("blog").Parse(blogTemplate)
	if err != nil {
		return err
	}

	os.MkdirAll("public/blog", 0755)
	file, err := os.Create("public/blog/index.html")
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
    <title>Works - RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        .project-item { margin-bottom: 50px; padding-bottom: 50px; border-bottom: 1px solid #eee; }
        .project-item:last-child { border-bottom: none; }
        .project-image { width: 100%; height: 400px; background: #f5f5f5; margin-bottom: 20px; border-radius: 8px; }
        .project-title { font-size: 26px; margin-bottom: 15px; }
        .project-description { color: #666; margin-bottom: 20px; }
        .project-links { display: flex; gap: 20px; }
        .project-links a { color: #0066cc; text-decoration: none; }
        .project-links a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Works</h1>
            <nav class="nav">
                <a href="/">Home</a>
                <a href="/works">Works</a>
                <a href="/blog">Blog</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            {{if .}}
            {{range .}}
            <article class="project-item">
                {{if .Image}}
                <div class="project-image" style="background-image: url('{{.Image}}'); background-size: cover; background-position: center;"></div>
                {{end}}
                <h2 class="project-title">{{.Title}}</h2>
                <p class="project-description">{{.Description}}</p>
                <div class="project-links">
                    {{if .CodeURL}}<a href="{{.CodeURL}}">Code</a>{{end}}
                    {{if .DemoURL}}<a href="{{.DemoURL}}">Demo</a>{{end}}
                </div>
            </article>
            {{end}}
            {{else}}
            <p>No projects yet. Check back soon!</p>
            {{end}}
        </main>

        <footer>
            <p>© 2026 RainyinSaiGon. Content on this site is licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
        </footer>
    </div>
</body>
</html>`

	tmpl, err := template.New("works").Parse(worksTemplate)
	if err != nil {
		return err
	}

	os.MkdirAll("public/works", 0755)
	file, err := os.Create("public/works/index.html")
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
    <title>About - RainyinSaiGon</title>
    <style>
        ` + baseCSS() + `
        article p { margin-bottom: 15px; color: #666; }
        article h2 { margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>About</h1>
            <nav class="nav">
                <a href="/">Home</a>
                <a href="/works">Works</a>
                <a href="/blog">Blog</a>
                <a href="/about">About</a>
            </nav>
        </header>

        <main>
            <article>
                <p>Hi! I'm RainyinSaiGon, a developer and builder interested in creating efficient, elegant solutions.</p>
                
                <h2>Skills</h2>
                <p>I work with various technologies including Go, web technologies, and other modern development tools.</p>
                
                <h2>Projects</h2>
                <p>Check out my <a href="/works">works page</a> to see what I've been building.</p>
                
                <h2>Contact</h2>
                <p>Feel free to reach out on <a href="https://github.com/RainyinSaiGon">GitHub</a>.</p>
            </article>
        </main>

        <footer>
            <p>© 2026 RainyinSaiGon. Content on this site is licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
        </footer>
    </div>
</body>
</html>`

	tmpl, err := template.New("about").Parse(aboutTemplate)
	if err != nil {
		return err
	}

	os.MkdirAll("public/about", 0755)
	file, err := os.Create("public/about/index.html")
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
    <title>{{.Title}} - Blog</title>
    <style>
        ` + baseCSS() + `
        article { color: #555; }
        article h2 { margin-top: 30px; }
        article p { margin-bottom: 15px; }
        article ul { margin-left: 20px; margin-bottom: 15px; }
        article li { margin-bottom: 8px; }
        article code { background: #f5f5f5; padding: 2px 6px; border-radius: 3px; font-family: 'Courier New', monospace; }
        article pre { background: #f5f5f5; padding: 15px; border-radius: 5px; overflow-x: auto; margin-bottom: 15px; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <a href="/blog" class="back-link">← Back to Blog</a>
            <h1>{{.Title}}</h1>
            <p class="post-meta">{{.Date}}</p>
        </header>

        <article>
            {{.Content}}
        </article>

        <footer>
            <p>© 2026 RainyinSaiGon. Content on this site is licensed under <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a>.</p>
        </footer>
    </div>
</body>
</html>`

	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return err
	}

	// Ensure blog directory exists in public
	os.MkdirAll("public/blog", 0755)

	file, err := os.Create(filepath.Join("public/blog", post.Slug, "index.html"))
	if err != nil {
		// Create the directory first
		os.MkdirAll(filepath.Join("public/blog", post.Slug), 0755)
		file, err = os.Create(filepath.Join("public/blog", post.Slug, "index.html"))
		if err != nil {
			return err
		}
	}
	defer file.Close()

	return tmpl.Execute(file, post)
}

# RainyinSaiGon

Personal portfolio and blog  built with a custom static site generator in Go.
Live at: **https://rainyinsaigon.github.io**

Inspired by [ziap.github.io](https://ziap.github.io)

---

## Stack

- **Generator**: Go (static site generator, no framework)
- **Styling**: Tailwind CSS (CDN) + Google Sans
- **Deployment**: GitHub Actions -> GitHub Pages (served from `/docs`)

## Project Structure

```
ch1/
 main.go                          # Entry point  calls builder.Build()
 go.mod
 Makefile
 content/
    posts/                       # Blog posts (.md with frontmatter)
    projects/                    # Portfolio projects (.md with frontmatter)
 internal/
     model/model.go               # Post, Project, HomeData types
     parser/parser.go             # Parse .md files + compute read time
     builder/builder.go           # Orchestration: parse -> sort -> render
     renderer/
         renderer.go              # html/template + embed.FS
         templates/               # HTML templates (home, blog, works, about)
         static/style.css         # Animations + post body typography
```

## Local Development

Build the site and serve it locally:

```bash
make serve
# opens at http://localhost:8080
```

Or manually:

```bash
go run .
python -m http.server 8080 --directory docs
```

## Writing a Post

Create a file in `content/posts/my-post.md`:

```
title: My Post Title
date: 2026-02-28
description: A short summary shown on the blog list page.
---
<p>Your HTML content here.</p>
<h2>A section heading</h2>
<p>More content...</p>
```

Read time is calculated automatically (~200 wpm).

## Adding a Project

Create a file in `content/projects/my-project.md`:

```
title: My Project
description: Short description shown on the works page.
image: https://example.com/screenshot.png
code: https://github.com/RainyinSaiGon/my-project
demo: https://my-project.vercel.app
featured: true
```

## Deployment

Push to `main`  GitHub Actions builds the site and deploys `docs/` automatically.

```bash
git add -A
git commit -m "your message"
git push origin main
```

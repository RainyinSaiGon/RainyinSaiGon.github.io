# RainyinSaiGon's Portfolio

A lightweight static site generator built in Go for creating a personal portfolio with blog and projects.

Inspired by [ziap.github.io](https://ziap.github.io/)

## Quick Start

```bash
# Build the generator
go build -o portfolio.exe

# Generate the site
./portfolio.exe

# Or use Makefile
make run
```

## Project Structure

```
.
├── main.go                     # Static site generator
├── go.mod                      # Go module
├── Makefile                    # Build commands
├── content/
│   ├── posts/                  # Blog posts (markdown)
│   └── projects/               # Projects/works (markdown)
├── public/                     # Generated HTML (deploy this to GitHub Pages)
│   ├── index.html              # Home page
│   ├── blog/
│   │   ├── index.html          # Blog listing
│   │   └── {slug}/index.html   # Individual blog posts
│   ├── works/index.html        # Projects listing
│   └── about/index.html        # About page
└── README.md
```

## Features

✅ **Multiple Pages**: Home, Blog, Works (Projects), About  
✅ **Fast Build**: Compiles to a single binary  
✅ **Clean Design**: Responsive HTML with no dependencies  
✅ **SEO-Friendly**: Proper meta tags and structure  
✅ **Easy Deployment**: Deploy the `public/` folder to GitHub Pages  
✅ **Date Sorting**: Posts automatically sorted by date (newest first)  
✅ **Featured Projects**: Mark projects as featured for home page display  

## Adding Content

### Blog Posts

Create `.md` files in `content/posts/` with this format:

```markdown
title: Your Post Title
date: 2026-02-28
description: Short description for the blog listing
---

<h2>Post Content</h2>
<p>Your HTML content here...</p>
```

**Note**: Posts are sorted by date in reverse order (newest first).

### Projects

Create `.md` files in `content/projects/` with this format:

```markdown
title: Project Title
description: Short description
image: /images/project.png
code: https://github.com/your-username/repo
demo: https://example.com
featured: true
---
```

- Set `featured: true` to display on the home page (limited to 3)
- Keep `featured: false` for projects only in the Works page
- `image`, `code`, and `demo` fields are optional

## Navigation Structure

- **/** → Home page (hero + featured projects + recent posts)
- **/blog** → All blog posts
- **/blog/{slug}** → Individual blog post
- **/works** → All projects
- **/about** → About page

## Deployment to GitHub Pages

1. **Create a repository** named `RainyinSaiGon.github.io` on GitHub
2. **Push to GitHub**:
   ```bash
   git init
   git add .
   git commit -m "Initial portfolio setup"
   git remote add origin https://github.com/RainyinSaiGon/RainyinSaiGon.github.io.git
   git push -u origin main
   ```
3. **Enable GitHub Pages**:
   - Go to Settings → Pages
   - Select "Deploy from a branch"
   - Choose `main` branch and `/root` directory
4. **Generate and deploy**:
   ```bash
   ./portfolio.exe
   git add public/
   git commit -m "Build: Generate portfolio site"
   git push
   ```

## Customization

All HTML templates and CSS are embedded in `main.go`. You can:

- Edit the `baseCSS()` function to change global styles
- Modify templates in `generateHomePage()`, `generateBlogPage()`, etc.
- Change colors, fonts, and spacing in the embedded CSS

## Available Make Commands

```bash
make fmt      # Format code
make vet      # Run Go vet
make build    # Build the generator (default)
make run      # Build and generate HTML
make clean    # Remove generated files
make help     # Show help
```

## Example Content

Sample files are provided:
- `content/posts/getting-started-go.md` - Intro blog post
- `content/projects/portfolio-builder.md` - Project example
- `content/projects/learning-go.md` - Another project example

## License

MIT License - Feel free to use this as a template for your own portfolio!

---

**Built with Go** 🚀


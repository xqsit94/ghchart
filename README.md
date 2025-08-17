# GitHub Charts API

> Fast, lightweight GitHub contribution charts with **light & dark theme support**

Generate beautiful GitHub contribution charts as transparent SVG images with support for both light and dark themes, custom colors, and seamless integration into any website or documentation.

![GitHub Charts Demo](https://your-domain.com/demo-image.png)

## ✨ Features

- 🎨 **Light & Dark Themes** - Automatic theme adaptation for any background
- 🌈 **Custom Colors** - Any hex color or predefined themes
- ⚡ **Lightning Fast** - Built with Go and optimized for performance  
- 🔍 **Transparent SVG** - Perfect integration with any design
- 📱 **Interactive Demo** - Beautiful homepage with live theme switching
- 🔗 **Simple API** - Clean, RESTful endpoints
- ⚙️ **Zero Config** - Works out of the box

## 🚀 Quick Start

### Basic Usage
```html
<!-- Light theme (default) -->
<img src="https://your-domain.com/username" alt="GitHub Contributions">

<!-- Dark theme -->
<img src="https://your-domain.com/dark:default/username" alt="GitHub Contributions">

<!-- Custom colors -->
<img src="https://your-domain.com/light:ff5722/username" alt="GitHub Contributions">
<img src="https://your-domain.com/dark:6366f1/username" alt="GitHub Contributions">
```

### Markdown
```markdown
![GitHub Contributions](https://your-domain.com/username)
![Dark Theme](https://your-domain.com/dark:default/username)
![Custom Orange](https://your-domain.com/light:ff5722/username)
```

## 📖 API Reference

### Light Theme (Default)
```
GET /{username}
GET /{color}/{username}
GET /{username}?theme=light
```

### Dark Theme  
```
GET /dark:default/{username}
GET /dark:{color}/{username}
GET /{username}?theme=dark
```

### Theme + Color Format
```
GET /{theme:color}/{username}
```

Where:
- `theme` - `light` or `dark`
- `color` - Hex color (without #) or special theme name
- `username` - GitHub username

### Examples

| URL | Description |
|-----|-------------|
| `/octocat` | Default light theme |
| `/dark:default/octocat` | Dark theme, default colors |
| `/light:ff5722/octocat` | Light theme, orange colors |
| `/dark:6366f1/octocat` | Dark theme, blue colors |
| `/halloween/octocat` | Light theme, halloween colors |
| `/dark:halloween/octocat` | Dark theme, halloween colors |

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd ghchart
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run .
```

The server will start on port 8080 by default, or use the `PORT` environment variable.

## Deployment

### Heroku
1. Create a Procfile:
```
web: ./ghchart
```

2. Deploy:
```bash
git add .
git commit -m "Initial commit"
heroku create your-app-name
git push heroku main
```

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ghchart .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ghchart .
EXPOSE 8080
CMD ["./ghchart"]
```

## Environment Variables

- `PORT` - Server port (default: 8080)

## 🎨 Color Schemes & Themes

### Predefined Themes
| Theme | Light Mode | Dark Mode | Description |
|-------|------------|-----------|-------------|
| `default` | GitHub green | Bright green | Classic GitHub colors |
| `halloween` | Orange/yellow | Bright orange | Halloween theme |
| `teal` | Teal/aqua | Bright teal | Ocean theme |

### Custom Colors
Use any hex color (without the `#`):
- `ff5722` - Orange
- `6366f1` - Blue  
- `8b5cf6` - Purple
- `ef4444` - Red
- `10b981` - Emerald

### Theme Combinations
```bash
# Light theme examples
/username                    # Default light
/ff5722/username            # Orange light (legacy format)
/light:ff5722/username      # Orange light (new format)
/halloween/username         # Halloween light

# Dark theme examples  
/dark:default/username      # Default dark
/dark:ff5722/username       # Orange dark
/dark:halloween/username    # Halloween dark

# Query parameter format
/username?theme=light       # Default light
/username?theme=dark        # Default dark
/ff5722/username?theme=dark # Orange dark
```

## 🌟 Advanced Usage

### Responsive Design
```html
<!-- Adapts to system theme -->
<picture>
  <source srcset="https://your-domain.com/dark:default/username" 
          media="(prefers-color-scheme: dark)">
  <img src="https://your-domain.com/username" alt="GitHub Contributions">
</picture>
```

### CSS Integration
```css
/* Light mode styling */
.github-chart-light {
  background: #ffffff;
  border-radius: 8px;
  padding: 16px;
}

/* Dark mode styling */
.github-chart-dark {
  background: #0d1117;
  border-radius: 8px;
  padding: 16px;
}
```

### JavaScript Dynamic Loading
```javascript
// Dynamic theme switching
function updateGitHubChart(username, isDark) {
  const img = document.getElementById('github-chart');
  const theme = isDark ? 'dark:default' : '';
  const url = `https://your-domain.com/${theme}${theme ? '/' : ''}${username}`;
  img.src = url;
}
```

## 🛠️ Development

### Local Development
```bash
# Run with hot reload (if using air)
air

# Or run directly
go run ./cmd/server

# Run tests
go test ./...

# Format code
go fmt ./...

# Build for production
go build -o ghchart ./cmd/server
```

### Project Structure
```
ghchart/
├── cmd/server/          # Application entry point
├── internal/
│   ├── chart/          # Chart generation logic
│   │   ├── colors.go   # Theme and color management
│   │   └── generator.go # SVG generation
│   ├── handlers/       # HTTP handlers
│   └── services/       # GitHub API integration
├── templates/          # HTML templates
└── go.mod
```

### API Response Format
All endpoints return SVG content with appropriate headers:
```
Content-Type: image/svg+xml
Cache-Control: public, max-age=86400
```

## 🚀 Deployment

### Heroku
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

```bash
# Quick deploy
heroku create your-app-name
git push heroku main
```

### Docker
```bash
# Build image
docker build -t ghchart .

# Run container
docker run -p 8080:8080 ghchart
```

### Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the original [ghchart](https://github.com/2016rshah/githubchart-api) project
- Built with [Go](https://golang.org/) and [Chi Router](https://github.com/go-chi/chi)
- GitHub contribution data fetched from public GitHub pages
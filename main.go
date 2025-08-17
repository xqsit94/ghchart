package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Cache-Control", "public, max-age=86400"))

	r.Get("/", indexHandler)
	r.Get("/{username}", chartHandler)
	r.Get("/{color}/{username}", colorChartHandler)

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, indexHTML)
}

func chartHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	username = removeExtension(username, ".svg")

	chart, err := generateChart(username, "")
	if err != nil {
		http.Error(w, "Failed to generate chart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(chart)
}

func colorChartHandler(w http.ResponseWriter, r *http.Request) {
	color := chi.URLParam(r, "color")
	username := chi.URLParam(r, "username")

	if username == "" || color == "" {
		http.Error(w, "Username and color are required", http.StatusBadRequest)
		return
	}

	username = removeExtension(username, ".svg")

	chart, err := generateChart(username, color)
	if err != nil {
		http.Error(w, "Failed to generate chart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(chart)
}

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>GitHub Chart API</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Noto Sans', Helvetica, Arial, sans-serif;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
            color: #24292f;
            background: #ffffff;
        }
        .header {
            text-align: center;
            margin-bottom: 40px;
            padding: 20px 0;
            border-bottom: 1px solid #d1d9e0;
        }
        .header h1 {
            font-size: 2.5rem;
            margin: 0 0 10px 0;
            font-weight: 600;
        }
        .header p {
            font-size: 1.2rem;
            color: #656d76;
            margin: 0;
        }
        .example {
            background: #f6f8fa;
            padding: 20px;
            margin: 20px 0;
            border-radius: 8px;
            border: 1px solid #d1d9e0;
        }
        .example h3 {
            margin-top: 0;
            color: #0969da;
            font-size: 1.1rem;
        }
        code {
            background: #afb8c133;
            padding: 3px 6px;
            border-radius: 4px;
            font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
            font-size: 0.9rem;
        }
        .chart-demo {
            margin-top: 15px;
            padding: 10px;
            background: white;
            border-radius: 6px;
            border: 1px solid #d1d9e0;
            text-align: center;
        }
        .chart-demo img {
            max-width: 100%;
            height: auto;
        }
        .examples-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
            margin: 20px 0;
        }
        .examples-list {
            background: #f6f8fa;
            padding: 20px;
            border-radius: 8px;
            border: 1px solid #d1d9e0;
        }
        .examples-list h3 {
            margin-top: 0;
            color: #0969da;
        }
        .examples-list ul {
            margin: 0;
            padding-left: 20px;
        }
        .examples-list li {
            margin: 8px 0;
        }
        .color-themes {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin-top: 10px;
        }
        .color-sample {
            display: inline-flex;
            align-items: center;
            background: white;
            padding: 8px 12px;
            border-radius: 6px;
            border: 1px solid #d1d9e0;
            font-size: 0.85rem;
        }
        .color-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            margin-right: 6px;
        }
        .footer {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #d1d9e0;
            text-align: center;
            color: #656d76;
            font-size: 0.9rem;
        }
        .api-url {
            background: #0969da;
            color: white;
            padding: 2px 6px;
            border-radius: 4px;
            font-weight: 500;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>GitHub Chart API</h1>
        <p>Embed GitHub contribution charts as SVG images</p>
    </div>

    <h2>Quick Start</h2>
    <div class="examples-grid">
        <div class="example">
            <h3>Default Chart</h3>
            <code class="api-url">GET</code> <code>/{username}</code>
            <p>Generate a contribution chart with GitHub's default green theme.</p>
            <div class="chart-demo">
                <img src="/torvalds" alt="Linus Torvalds GitHub contributions">
                <p><small>Example: Linus Torvalds' contributions</small></p>
            </div>
        </div>

        <div class="example">
            <h3>Custom Colors</h3>
            <code class="api-url">GET</code> <code>/{hex-color}/{username}</code>
            <p>Generate a chart with custom color scheme using any hex color.</p>
            <div class="chart-demo">
                <img src="/ff5722/torvalds" alt="Custom colored GitHub contributions">
                <p><small>Example: Orange theme (#ff5722)</small></p>
            </div>
        </div>
    </div>

    <div class="examples-grid">
        <div class="examples-list">
            <h3>Color Themes</h3>
            <p>Use any hex color to generate custom themes:</p>
            <div class="color-themes">
                <div class="color-sample">
                    <div class="color-dot" style="background: #2196f3;"></div>
                    <code>2196f3</code>
                </div>
                <div class="color-sample">
                    <div class="color-dot" style="background: #ff5722;"></div>
                    <code>ff5722</code>
                </div>
                <div class="color-sample">
                    <div class="color-dot" style="background: #9c27b0;"></div>
                    <code>9c27b0</code>
                </div>
                <div class="color-sample">
                    <div class="color-dot" style="background: #f44336;"></div>
                    <code>f44336</code>
                </div>
            </div>
            <p><small>Or use special themes: <code>halloween</code>, <code>teal</code></small></p>
        </div>

        <div class="examples-list">
            <h3>API Examples</h3>
            <ul>
                <li><code>/octocat</code> - GitHub mascot (empty chart)</li>
                <li><code>/torvalds</code> - Linus Torvalds</li>
                <li><code>/2196f3/torvalds</code> - Blue theme</li>
                <li><code>/ff5722/octocat</code> - Orange theme</li>
                <li><code>/halloween/torvalds</code> - Halloween theme</li>
            </ul>
        </div>
    </div>

    <h2>Usage</h2>
    <div class="example">
        <h3>Embed in Markdown</h3>
        <code>![GitHub Chart](https://your-domain.com/username)</code>
        <br><br>
        <h3>Embed in HTML</h3>
        <code>&lt;img src="https://your-domain.com/username" alt="GitHub Contributions"&gt;</code>
        <br><br>
        <h3>Direct SVG Access</h3>
        <code>https://your-domain.com/username.svg</code>
    </div>

    <div class="footer">
        <p><strong>GitHub Chart API</strong> - Built with Go & Chi Router</p>
        <p>Fetches real-time contribution data from GitHub profiles</p>
    </div>
</body>
</html>
`

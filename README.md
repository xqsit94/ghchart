# GitHub Chart API (Go)

A Go implementation of the GitHub contribution chart API. Embed GitHub contribution charts as SVG images.

## Features

- Generate GitHub contribution charts as SVG images
- Support for custom color schemes
- Built with Go and chi router
- Compatible with the original Ruby API endpoints
- Fast and lightweight

## API Endpoints

### Default Chart
```
GET /{username}
```
Returns a GitHub contribution chart with the default green color scheme.

### Custom Color Chart
```
GET /{color}/{username}
```
Returns a GitHub contribution chart with a custom color scheme based on the provided hex color.

### Examples
- `/octocat` - Default green chart for user "octocat"
- `/2196f3/octocat` - Blue themed chart for user "octocat"
- `/ff5722/octocat` - Orange themed chart for user "octocat"

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

## Color Schemes

The API supports several color schemes:

- **Default**: GitHub's standard green theme
- **Custom Hex**: Any valid hex color (e.g., `ff5722`, `2196f3`)
- **Special Themes**:
  - `halloween` - Orange/black theme
  - `teal` - Teal/aqua theme

## Development

Run tests:
```bash
go test ./...
```

Format code:
```bash
go fmt ./...
```

## License

MIT License
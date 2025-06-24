# CV Builder

A powerful AI-powered CV builder that helps users create professional, ATS-friendly resumes through an interactive chat interface. The application uses Google's Gemini AI to generate and refine CV content, and can export the final result as a professionally formatted PDF.

![CV Builder Screenshot](https://via.placeholder.com/800x400/4a6cf7/ffffff?text=CV+Builder+Screenshot)

## Features

- 🤖 AI-powered CV generation using Google's Gemini AI
- 💬 Interactive chat interface for CV creation and refinement
- 📝 ATS-friendly resume formatting
- 🌐 Web scraping capability for importing content from URLs
- 📄 PDF export with professional formatting
- ⚡ Fast and responsive API endpoints
- 🐳 Docker support for easy deployment

## Screenshots

### Main Interface
![Main Interface](https://via.placeholder.com/600x400/4a6cf7/ffffff?text=Main+Interface)

*Interactive chat interface with AI-powered CV generation*

### Chat Interface
![Chat Interface](https://via.placeholder.com/600x400/4a6cf7/ffffff?text=Chat+Interface)

*Real-time chat with AI to build your CV step by step*

### PDF Export
![PDF Export](https://via.placeholder.com/600x400/4a6cf7/ffffff?text=PDF+Export)

*Professional PDF export with ATS-friendly formatting*

### Docker Deployment
![Docker Deployment](https://via.placeholder.com/600x400/4a6cf7/ffffff?text=Docker+Deployment)

*Easy deployment with Docker and Docker Compose*

## Prerequisites

- Go 1.23.0 or higher (for local development)
- Google Cloud API key for Gemini AI
- Chrome/Chromium (for PDF generation)
- Docker and Docker Compose (for containerized deployment)

## Installation

### Local Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/cv-builder.git
cd cv-builder
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory and add your Google API key:
```
GOOGLE_API_KEY=your_api_key_here
```

4. (Optional) Set the port (default is 9090):
```
PORT=9090
```

### Docker Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/cv-builder.git
cd cv-builder
```

2. Create a `.env` file with your Google API key and port:
```
GOOGLE_API_KEY=your_api_key_here
PORT=9090
```

3. Build and run with Docker Compose:
```bash
docker-compose up --build
```

The application will be available at [http://localhost:9090](http://localhost:9090).

## Usage

### Local Usage

Start the server:
```bash
go run main.go
```

The server will start on the port specified by the `PORT` environment variable (default: 9090).

### Docker Usage

The application will be available at `http://localhost:9090` after running `docker-compose up`.

## API Endpoints

### Chat Endpoint
```http
POST /chat
Content-Type: application/json

{
    "message": "Your message here",
    "history": [
        {
            "role": "user",
            "parts": ["Previous message"]
        }
    ]
}
```

### PDF Generation Endpoint
```http
POST /pdf
Content-Type: application/json

{
    "html": "Your CV HTML content"
}
```

## Docker & Deployment Notes

- The Go application must listen on the port specified by the `PORT` environment variable (default: 9090).
- All static files (such as `templates/index.html`) must be present in the Docker image and referenced with the correct relative path (relative to `/app` in the container).
- If you get a `404 Not Found` at `/`, ensure your Go code includes a handler for `/` that serves `templates/index.html`.
- Example Go handler for serving the frontend:

```go
r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/index.html")
})
```

- If you change the port, update both the `.env` file and the `docker-compose.yml` accordingly.

## Troubleshooting

- **404 Not Found at `/` in Docker:**
  - Make sure your Go code has a handler for `/` that serves the frontend HTML file.
  - Check that `templates/index.html` exists in the Docker image (`/app/templates/index.html`).
  - Ensure the Go app is listening on the correct port (`PORT` environment variable).
  - Check container logs with `docker-compose logs` for errors.
  - Test inside the container: `docker-compose exec cv-builder sh` then `curl localhost:9090/`.

- **Environment variables not working:**
  - Ensure `.env` is present and variables are referenced in both Go code and `docker-compose.yml`.

- **Static files not loading:**
  - Check file paths in your Go code and Docker image.

## Dependencies

- [github.com/google/generative-ai-go](https://github.com/google/generative-ai-go) - Google's Gemini AI client
- [github.com/chromedp/chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol for PDF generation
- [github.com/gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing
- [github.com/joho/godotenv](https://github.com/joho/godotenv) - Environment variable management

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
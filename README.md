# CV Builder

A powerful AI-powered CV builder that helps users create professional, ATS-friendly resumes through an interactive chat interface. The application uses Google's Gemini AI to generate and refine CV content, and can export the final result as a professionally formatted PDF.

## Features

- 🤖 AI-powered CV generation using Google's Gemini AI
- 💬 Interactive chat interface for CV creation and refinement
- 📝 ATS-friendly resume formatting
- 🌐 Web scraping capability for importing content from URLs
- 📄 PDF export with professional formatting
- ⚡ Fast and responsive API endpoints
- 🐳 Docker support for easy deployment

## Prerequisites

- Go 1.23.0 or higher
- Google Cloud API key for Gemini AI
- Chrome/Chromium (for PDF generation)
- Docker and Docker Compose (optional, for containerized deployment)

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

### Docker Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/cv-builder.git
cd cv-builder
```

2. Create a `.env` file with your Google API key:
```
GOOGLE_API_KEY=your_api_key_here
```

3. Build and run with Docker Compose:
```bash
docker-compose up --build
```

## Usage

### Local Usage

1. Start the server:
```bash
go run main.go
```

2. The server will start on the default port. You can interact with the API using the following endpoints:

- `POST /chat` - Chat with the AI to build your CV
- `POST /pdf` - Generate a PDF from your CV HTML

### Docker Usage

The application will be available at `http://localhost:8080` after running `docker-compose up`.

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
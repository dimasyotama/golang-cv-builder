package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/generative-ai-go/genai"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// This version includes aggressive timeouts and the most comprehensive logging possible.

type ChatMessage struct{ Role string; Parts []string `json:"parts"`}
type ChatRequest struct{ Message string; History []ChatMessage }
type AIResponse struct{ AIResponse string `json:"ai_response"`; CVHTML string `json:"cv_html"`}

type FinalResponse struct {
    AIResponse string        `json:"ai_response"`
    CVHTML     string        `json:"cv_html"`
    History    []ChatMessage `json:"history"`
}
type PDFRequest struct{ HTML string `json:"html"`}

func scrapeWebsite(url string) (string, error) {
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	res, err := client.Get(url)
	if err != nil {return "", err}
	defer res.Body.Close()

	if res.StatusCode != 200 {return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {return "", err}
	doc.Find("script, style").Each(func(i int, s *goquery.Selection) {s.Remove()})
	text := strings.Join(strings.Fields(doc.Text()), " ")
	if len(text) > 4000 {text = text[:4000]}
	return text, nil
}

func getCVParsingPrompt(history []ChatMessage) string {
	var conversationStr strings.Builder
	for _, msg := range history {
		conversationStr.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Parts[0]))
	}
	return fmt.Sprintf(`
    You are an expert ATS-friendly CV writer... (prompt unchanged)
    CONVERSATION HISTORY:
    %s
    ---
    EXAMPLE OUTPUT FORMAT:
    {
      "ai_response": "That's a great start, Budi! I've added your name to the CV. What would you like to add next, such as your work experience or education?",
      "cv_html": "<h1>Budi Santoso</h1>"
    }
    ---
    `, conversationStr.String())
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("----------------- NEW CHAT REQUEST -----------------")
	log.Println("[DEBUG] /chat endpoint hit. Starting request processing...")

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[FATAL ERROR] Failed to decode request JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("[DEBUG] Successfully decoded request body. User message: '%s'", req.Message)

	userMessage := req.Message
	if strings.Contains(userMessage, "http") {
		log.Printf("[DEBUG] URL detected. Scraping website: %s", userMessage)
		scrapedContent, err := scrapeWebsite(userMessage)
		if err != nil {
            log.Printf("[ERROR] Failed to scrape website: %v", err)
			userMessage += fmt.Sprintf("\n[Could not access the URL: %s.]", userMessage)
		} else {
			userMessage += fmt.Sprintf("\n[Scraped content from %s: %s]", userMessage, scrapedContent)
		}
	}
	req.History = append(req.History, ChatMessage{Role: "user", Parts: []string{userMessage}})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
	if err != nil {
		log.Printf("[FATAL ERROR] Failed to create Gemini client. Check API Key and .env file. Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()
	log.Println("[DEBUG] Gemini client created successfully.")

	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := getCVParsingPrompt(req.History)
	log.Printf("[DEBUG] Sending this prompt to Gemini:\n---\n%s\n---", prompt)
	
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	
	if err != nil {
		log.Printf("[FATAL ERROR] AI content generation failed. This could be a timeout or an API issue. Error: %v", err)
		http.Error(w, "AI generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[DEBUG] Successfully received response from Gemini API.")

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("[FATAL ERROR] Received empty response from AI.")
		http.Error(w, "Received empty response from AI", http.StatusInternalServerError)
		return
	}
	aiResponseText := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	log.Printf("[DEBUG] Raw response from AI:\n---\n%s\n---", aiResponseText)

	// Sanitize the response to remove markdown
	sanitizedJSON := strings.TrimSpace(aiResponseText)
	sanitizedJSON = strings.TrimPrefix(sanitizedJSON, "```json")
	sanitizedJSON = strings.TrimSuffix(sanitizedJSON, "```")
	sanitizedJSON = strings.TrimSpace(sanitizedJSON)

	var aiResp AIResponse
	if err := json.Unmarshal([]byte(sanitizedJSON), &aiResp); err != nil {
		log.Printf("[FATAL ERROR] Failed to parse SANITIZED JSON from AI response. Error: %v", err)
		http.Error(w, "Failed to parse AI response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("[DEBUG] Successfully parsed AI response.")

	req.History = append(req.History, ChatMessage{Role: "model", Parts: []string{aiResp.AIResponse}})
	finalResp := FinalResponse{AIResponse: aiResp.AIResponse,CVHTML: aiResp.CVHTML,History: req.History}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResp)
	log.Println("[DEBUG] Successfully sent final response to client.")
	log.Println("----------------- CHAT REQUEST COMPLETE -----------------")
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
    var req PDFRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pdfStyles := `
	<style>
		body { 
			font-family: 'Helvetica', 'Arial', sans-serif; 
			font-size: 10.5pt; 
			line-height: 1.5; 
			color: #333;
			margin: 30px;
		}
		h1 {
			font-size: 26pt; text-align: center; margin-bottom: 4px;
			font-weight: 600; color: #000;
		}
		/* FIX: This selector now targets any <p> immediately after <h1> */
		h1 + p {
			text-align: center; font-size: 10pt; margin-top: 0;
			margin-bottom: 20px; color: #222;
		}
		h2 {
			font-size: 11pt; font-weight: bold; color: #000;
			border-bottom: 1.5px solid #000; padding-bottom: 4px;
			margin-top: 20px; margin-bottom: 12px; text-transform: uppercase;
		}
		h3 {
			font-size: 10.5pt; font-weight: bold; color: #111;
			margin-bottom: 2px; margin-top: 12px;
		}
		p { margin: 0 0 8px 0; color: #333; }
		ul { padding-left: 18px; margin: 5px 0; list-style-type: '•  ';}
		li { margin-bottom: 6px; padding-left: 4px; }
	</style>
	`

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 25*time.Second) 
	defer cancel()

	var pdfBuffer []byte
	fullHTML := pdfStyles + req.HTML
	
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, fullHTML).Do(ctx)
		}),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(false).
				WithMarginTop(0.4). WithMarginBottom(0.4).
				WithMarginLeft(0.4). WithMarginRight(0.4).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuffer = buf
			return nil
		}),
	)

	if err != nil {
		log.Printf("Failed to generate PDF with chromedp: %v", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"AI-Architected-CV.pdf\"")
	w.Write(pdfBuffer)
	log.Println("[DEBUG] Successfully generated and sent PDF.")
}


func main() {
    log.Println("--- Application Starting ---")
	if err := godotenv.Load(); err != nil {log.Println("[WARN] No .env file found. Make sure it exists if API key is not set elsewhere.")}

	r := mux.NewRouter()
	r.HandleFunc("/chat", chatHandler).Methods("POST")
	r.HandleFunc("/download_pdf", pdfHandler).Methods("POST")
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {http.ServeFile(w, r, "templates/index.html")})

	port := "8000"
	log.Printf("✅ Go backend server starting on http://localhost:%s.", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {log.Fatal(err)}
}
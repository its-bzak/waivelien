package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-pdf/fpdf"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type InvoiceData struct {
	ProjectName       string  `json:"project_name"`
	GCName            string  `json:"gc_name"`
	SubcontractorName string  `json:"subcontractor_name"`
	InvoiceAmount     float64 `json:"invoice_amount"`
	ThroughDate       string  `json:"through_date"`
}

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Warning: No .env file found")
	}
}

// --- MIDDLEWARE ---
// This wrapper tells the browser it is safe for React to talk to this Go server
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests from the browser
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// --- ENDPOINTS ---

// GET /api/invoice: Reads the JSON and sends it to React
func getInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	fileData, err := os.ReadFile("../mock/invoice_data.json")
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(fileData)
}

// POST /api/sign: Triggers the PDF generation
func signHandler(w http.ResponseWriter, r *http.Request) {
	fileData, err := os.ReadFile("../mock/invoice_data.json")
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	var invoice InvoiceData
	json.Unmarshal(fileData, &invoice)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "CONDITIONAL WAIVER AND RELEASE")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)

	bodyText := fmt.Sprintf("Project: %s\nGC: %s\nSubcontractor: %s\nAmount: $%.2f",
		invoice.ProjectName, invoice.GCName, invoice.SubcontractorName, invoice.InvoiceAmount)

	pdf.MultiCell(0, 8, bodyText, "", "L", false)
	pdf.Ln(12)
	pdf.Cell(0, 10, "Digital Signature: ___________________________")

	err = pdf.OutputFileAndClose("waiver_output.pdf")
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "PDF Generated"}`))
}

// POST /api/send-sms: Sends the PDF via Twilio SMS
func sendSMSHandler(w http.ResponseWriter, r *http.Request) {

	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	apiKey := os.Getenv("TWILIO_API_KEY")
	apiSecret := os.Getenv("TWILIO_API_SECRET")
	verifySid := os.Getenv("TWILIO_VERIFY_SERVICE_SID")

	// Initialize Twilio client
	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username:   apiKey,
			Password:   apiSecret,
			AccountSid: accountSid,
		},
	)

	// Create verification request
	params := &verify.CreateVerificationParams{}
	params.SetTo(os.Getenv("MY_CELL_NUMBER"))
	params.SetChannel("sms")

	// Send verification code
	resp, err := client.VerifyV2.CreateVerification(verifySid, params)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send verification: %v", err),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Verification code sent!",
		"sid":     *resp.Sid,
	})
}

func main() {
	// Register the endpoints and wrap them in our CORS middleware
	http.HandleFunc("/api/invoice", enableCORS(getInvoiceHandler))
	http.HandleFunc("/api/sign", enableCORS(signHandler))
	http.HandleFunc("/api/send-code", enableCORS(sendSMSHandler))

	fmt.Println("Go server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

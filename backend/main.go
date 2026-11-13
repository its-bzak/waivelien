package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-pdf/fpdf"
)

// 1. Define a "struct" (the Go version of a class/object) that exactly matches our JSON shape.
type InvoiceData struct {
	ProjectName       string  `json:"project_name"`
	GCName            string  `json:"gc_name"`
	SubcontractorName string  `json:"subcontractor_name"`
	InvoiceAmount     float64 `json:"invoice_amount"`
	ThroughDate       string  `json:"through_date"`
}

func main() {
	// 2. Read the raw bytes from the JSON file
	fileData, err := os.ReadFile("../mock/invoice_data.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// 3. Parse those bytes into our Go struct
	var invoice InvoiceData
	err = json.Unmarshal(fileData, &invoice)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// 4. Initialize a new PDF (Portrait, millimeters, A4 size)
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// 5. Write the Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "CONDITIONAL WAIVER AND RELEASE")
	pdf.Ln(12) // Line break

	// 6. Write the body text, injecting the data from our JSON
	pdf.SetFont("Arial", "", 12)

	// fmt.Sprintf allows us to format a string with variables easily
	bodyText := fmt.Sprintf(
		"Project Name: %s\nGeneral Contractor: %s\nSubcontractor: %s\nPayment Amount: $%.2f\nThrough Date: %s",
		invoice.ProjectName,
		invoice.GCName,
		invoice.SubcontractorName,
		invoice.InvoiceAmount,
		invoice.ThroughDate,
	)

	// MultiCell allows text to wrap to the next line if it's long
	pdf.MultiCell(0, 8, bodyText, "", "L", false)
	pdf.Ln(12)

	// 7. Add a placeholder signature line
	pdf.SetFont("Arial", "I", 12)
	pdf.Cell(0, 10, "Digital Signature: ___________________________")

	// 8. Save the PDF to your hard drive
	err = pdf.OutputFileAndClose("waiver_output.pdf")
	if err != nil {
		log.Fatalf("Failed to save PDF: %v", err)
	}

	fmt.Println("Success! waiver_output.pdf has been generated in your folder.")
}

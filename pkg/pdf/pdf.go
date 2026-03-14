package pdf

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CategoryBreakdownRow represents one category breakdown row for PDF export.
type CategoryBreakdownRow struct {
	Category    string
	Subcategory string
	Amount      float64
	Percentage  float64
	Type        string
}

// TransactionRow represents one transaction row for PDF export.
type TransactionRow struct {
	Date        string
	Name        string
	Category    string
	Subcategory string
	Amount      float64
	Notes       string
}

// StatsPDFData contains all data required to generate a stats PDF report.
type StatsPDFData struct {
	PeriodLabel   string
	AccountName   string
	CurrencyCode  string
	TotalIncome   float64
	TotalExpense  float64
	Net           float64
	BreakdownRows []CategoryBreakdownRow
	Transactions  []TransactionRow
}

// GenerateStatsPDF creates a PDF report for stats export.
func GenerateStatsPDF(data StatsPDFData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 15)

	generatedAt := time.Now().Format("2006-01-02 15:04:05")

	pdf.SetFooterFunc(func() {
		pdf.SetY(-10)
		pdf.SetFont("Helvetica", "", 8)
		footer := fmt.Sprintf("Generated %s | Page %d/{nb}", generatedAt, pdf.PageNo())
		pdf.CellFormat(0, 6, footer, "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("{nb}")

	addHeader(pdf, data, generatedAt)
	addSummaryTable(pdf, data)
	addBreakdownTable(pdf, data)
	addTransactions(pdf, data)

	var out bytes.Buffer
	if err := pdf.Output(&out); err != nil {
		return nil, fmt.Errorf("pdf.GenerateStatsPDF.Output: %w", err)
	}
	return out.Bytes(), nil
}

func addHeader(pdf *fpdf.Fpdf, data StatsPDFData, generatedAt string) {
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(0, 10, "Expenditure Tracker", "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(0, 7, "Generated on "+generatedAt, "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Period: %s | Account: %s", data.PeriodLabel, data.AccountName), "", 1, "L", false, 0, "")
	pdf.Ln(3)
}

func addSummaryTable(pdf *fpdf.Fpdf, data StatsPDFData) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "Summary", "", 1, "L", false, 0, "")

	headers := []string{"Total Income", "Total Expense", "Net Balance"}
	values := []string{
		formatMoney(data.CurrencyCode, data.TotalIncome),
		formatMoney(data.CurrencyCode, data.TotalExpense),
		formatMoney(data.CurrencyCode, data.Net),
	}

	pdf.SetFont("Helvetica", "B", 10)
	for _, h := range headers {
		pdf.CellFormat(63, 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 10)
	for _, v := range values {
		pdf.CellFormat(63, 8, v, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(12)
}

func addBreakdownTable(pdf *fpdf.Fpdf, data StatsPDFData) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "Category Breakdown", "", 1, "L", false, 0, "")

	headers := []string{"Type", "Category", "Subcategory", "Amount", "% of Total"}
	widths := []float64{24, 40, 48, 38, 30}

	pdf.SetFont("Helvetica", "B", 9)
	for idx, h := range headers {
		pdf.CellFormat(widths[idx], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 9)
	titleCaser := cases.Title(language.Und)
	for _, row := range data.BreakdownRows {
		pdf.CellFormat(widths[0], 7, titleCaser.String(row.Type), "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 7, safeText(row.Category), "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 7, safeText(row.Subcategory), "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 7, formatMoney(data.CurrencyCode, row.Amount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[4], 7, fmt.Sprintf("%.2f%%", row.Percentage), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(8)
}

func addTransactions(pdf *fpdf.Fpdf, data StatsPDFData) {
	if len(data.Transactions) == 0 {
		pdf.SetFont("Helvetica", "I", 10)
		pdf.CellFormat(0, 8, "No transactions available for selected filters.", "", 1, "L", false, 0, "")
		return
	}

	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(0, 8, "Transaction History", "", 1, "L", false, 0, "")

	headers := []string{"Date", "Name", "Category", "Subcategory", "Amount"}
	widths := []float64{24, 48, 38, 40, 30}

	pdf.SetFont("Helvetica", "B", 9)
	for idx, h := range headers {
		pdf.CellFormat(widths[idx], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 9)
	for idx, row := range data.Transactions {
		if idx%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.CellFormat(widths[0], 7, safeText(row.Date), "1", 0, "L", true, 0, "")
		pdf.CellFormat(widths[1], 7, safeText(row.Name), "1", 0, "L", true, 0, "")
		pdf.CellFormat(widths[2], 7, safeText(row.Category), "1", 0, "L", true, 0, "")
		pdf.CellFormat(widths[3], 7, safeText(row.Subcategory), "1", 0, "L", true, 0, "")
		pdf.CellFormat(widths[4], 7, formatMoney(data.CurrencyCode, row.Amount), "1", 0, "R", true, 0, "")
		pdf.Ln(-1)

		if row.Notes != "" {
			pdf.SetFont("Helvetica", "I", 8)
			notes := row.Notes
			if len(notes) > 50 {
				notes = notes[:50] + "..."
			}
			pdf.CellFormat(0, 5, "  Note: "+safeText(notes), "LRB", 1, "L", true, 0, "")
			pdf.SetFont("Helvetica", "", 9)
		}
	}
}

func safeText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func formatMoney(currency string, amount float64) string {
	if currency == "" {
		currency = "INR"
	}
	return fmt.Sprintf("%s %.2f", currency, amount)
}

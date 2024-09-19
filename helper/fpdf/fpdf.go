package fpdf

import (
	"bytes"
	"fmt"

	"github.com/unidoc/unipdf/v3/model"
)

// MergePDFBytes merges two PDF files provided as []byte and returns the merged result as []byte.
func MergePDFBytes(pdf1, pdf2 []byte) ([]byte, error) {
	// Load first PDF from []byte
	reader1 := bytes.NewReader(pdf1)
	pdfReader1, err := model.NewPdfReader(reader1)
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf1: %v", err)
	}

	// Load second PDF from []byte
	reader2 := bytes.NewReader(pdf2)
	pdfReader2, err := model.NewPdfReader(reader2)
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf2: %v", err)
	}

	// Create a new PDF writer to merge into
	pdfWriter := model.NewPdfWriter()

	// Append all pages from the first PDF
	numPages1, err := pdfReader1.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get number of pages in pdf1: %v", err)
	}
	for i := 1; i <= numPages1; i++ {
		page, err := pdfReader1.GetPage(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d from pdf1: %v", i, err)
		}
		pdfWriter.AddPage(page)
	}

	// Append all pages from the second PDF
	numPages2, err := pdfReader2.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get number of pages in pdf2: %v", err)
	}
	for i := 1; i <= numPages2; i++ {
		page, err := pdfReader2.GetPage(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d from pdf2: %v", i, err)
		}
		pdfWriter.AddPage(page)
	}

	// Write the merged PDF to a buffer
	var buf bytes.Buffer
	err = pdfWriter.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write merged PDF: %v", err)
	}

	return buf.Bytes(), nil
}

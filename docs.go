package main

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	"github.com/extrame/xls"
	"github.com/fumiama/go-docx"
	"github.com/ledongthuc/pdf"
	"github.com/xuri/excelize/v2"
)

var projectRoot string

func readFileContent(path string) (string, error) {
	resolved := path
	if !filepath.IsAbs(path) {
		resolved = filepath.Join(projectRoot, path)
	}
	cleaned := filepath.Clean(resolved)
	if !strings.HasPrefix(cleaned, projectRoot) {
		return "", os.ErrPermission
	}

	ext := strings.ToLower(filepath.Ext(cleaned))
	switch ext {
	case ".eml":
		return readEml(cleaned)
	case ".pdf":
		return readPdf(cleaned)
	case ".docx":
		return readDocx(cleaned)
	case ".doc":
		return readDoc(cleaned)
	case ".xlsx":
		return readXlsx(cleaned)
	case ".xls":
		return readXls(cleaned)
	default:
		data, err := os.ReadFile(cleaned)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}

// readEml parses a .eml file and extracts headers + text content,
// skipping base64 attachments.
func readEml(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	msg, err := mail.ReadMessage(f)
	if err != nil {
		return "", fmt.Errorf("failed to parse email: %w", err)
	}

	var out strings.Builder
	// Write key headers
	for _, h := range []string{"From", "To", "Cc", "Subject", "Date"} {
		if v := msg.Header.Get(h); v != "" {
			fmt.Fprintf(&out, "%s: %s\n", h, v)
		}
	}
	out.WriteString("\n")

	// Extract text parts from the body
	contentType := msg.Header.Get("Content-Type")
	body, err := extractTextParts(contentType, msg.Body)
	if err != nil {
		return "", fmt.Errorf("failed to extract email body: %w", err)
	}
	if body == "" {
		body = "[No text content found in email]"
	}
	out.WriteString(body)
	return out.String(), nil
}

// extractTextParts recursively walks MIME parts and returns text/plain content.
// Falls back to text/html if no plain text is found.
func extractTextParts(contentType string, body io.Reader) (string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		// Not MIME — treat as plain text
		data, err := io.ReadAll(body)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	if strings.HasPrefix(mediaType, "text/plain") {
		data, err := io.ReadAll(body)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	if strings.HasPrefix(mediaType, "text/html") {
		data, err := io.ReadAll(body)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return "", nil // skip non-text parts (images, attachments)
	}

	boundary := params["boundary"]
	if boundary == "" {
		return "", nil
	}

	reader := multipart.NewReader(body, boundary)
	var plainTexts, htmlTexts []string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		partCT := part.Header.Get("Content-Type")
		if partCT == "" {
			partCT = "text/plain"
		}
		text, err := extractTextParts(partCT, part)
		if err != nil {
			continue
		}
		if text == "" {
			continue
		}
		partMedia, _, _ := mime.ParseMediaType(partCT)
		if strings.HasPrefix(partMedia, "text/plain") {
			plainTexts = append(plainTexts, text)
		} else {
			htmlTexts = append(htmlTexts, text)
		}
	}

	if len(plainTexts) > 0 {
		return strings.Join(plainTexts, "\n\n"), nil
	}
	if len(htmlTexts) > 0 {
		return strings.Join(htmlTexts, "\n\n"), nil
	}

	// If we got multipart but found text in recursive calls without tracking type
	// try a raw re-read
	return "", nil
}


func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var out strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		out.WriteString(text)
		out.WriteString("\n")
	}

	content := strings.TrimSpace(out.String())
	if content == "" {
		return "[No extractable text in PDF — may be scanned/image-based]", nil
	}
	return content, nil
}

func extractRunText(r *docx.Run) string {
	var s strings.Builder
	for _, child := range r.Children {
		if t, ok := child.(*docx.Text); ok {
			s.WriteString(t.Text)
		}
	}
	return s.String()
}

func readDocx(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return "", err
	}

	doc, err := docx.Parse(f, fi.Size())
	if err != nil {
		return "", fmt.Errorf("failed to parse docx: %w", err)
	}

	var out strings.Builder
	for _, item := range doc.Document.Body.Items {
		switch v := item.(type) {
		case *docx.Paragraph:
			for _, child := range v.Children {
				if r, ok := child.(*docx.Run); ok {
					out.WriteString(extractRunText(r))
				}
			}
			out.WriteString("\n")
		case *docx.Table:
			for _, row := range v.TableRows {
				for ci, cell := range row.TableCells {
					if ci > 0 {
						out.WriteString("\t")
					}
					for _, p := range cell.Paragraphs {
						for _, child := range p.Children {
							if r, ok := child.(*docx.Run); ok {
								out.WriteString(extractRunText(r))
							}
						}
					}
				}
				out.WriteString("\n")
			}
		}
	}

	content := strings.TrimSpace(out.String())
	if content == "" {
		return "[No extractable text in docx]", nil
	}
	return content, nil
}

func readDoc(path string) (string, error) {
	// Legacy .doc is a complex binary format. Read raw bytes and extract
	// printable text as a best-effort approach.
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var out strings.Builder
	inText := false
	for _, b := range data {
		if b >= 0x20 && b < 0x7F || b == '\n' || b == '\r' || b == '\t' {
			out.WriteByte(b)
			inText = true
		} else {
			if inText {
				out.WriteByte(' ')
				inText = false
			}
		}
	}

	content := strings.TrimSpace(out.String())
	if content == "" {
		return "[No extractable text in .doc file]", nil
	}
	return content, nil
}

func readXlsx(path string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to open xlsx: %w", err)
	}
	defer f.Close()

	var out strings.Builder
	for _, sheet := range f.GetSheetList() {
		fmt.Fprintf(&out, "=== %s ===\n", sheet)
		rows, err := f.GetRows(sheet)
		if err != nil {
			continue
		}
		for _, row := range rows {
			out.WriteString(strings.Join(row, "\t"))
			out.WriteString("\n")
		}
		out.WriteString("\n")
	}

	content := strings.TrimSpace(out.String())
	if content == "" {
		return "[No data in xlsx]", nil
	}
	return content, nil
}

func readXls(path string) (string, error) {
	wb, err := xls.Open(path, "utf-8")
	if err != nil {
		return "", fmt.Errorf("failed to open xls: %w", err)
	}

	var out strings.Builder
	for i := 0; i < wb.NumSheets(); i++ {
		sheet := wb.GetSheet(i)
		if sheet == nil {
			continue
		}
		fmt.Fprintf(&out, "=== %s ===\n", sheet.Name)
		for r := 0; r <= int(sheet.MaxRow); r++ {
			row := sheet.Row(r)
			if row == nil {
				continue
			}
			for c := row.FirstCol(); c < row.LastCol(); c++ {
				if c > row.FirstCol() {
					out.WriteString("\t")
				}
				out.WriteString(row.Col(c))
			}
			out.WriteString("\n")
		}
		out.WriteString("\n")
	}

	content := strings.TrimSpace(out.String())
	if content == "" {
		return "[No data in xls]", nil
	}
	return content, nil
}

func listDir(dir string) ([]string, error) {
	root := filepath.Join(projectRoot, dir)
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			rel, _ := filepath.Rel(projectRoot, path)
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}

func writeFile(dir, filename, content string) error {
	dest := filepath.Join(projectRoot, dir, filepath.Base(filename))
	cleaned := filepath.Clean(dest)
	if !strings.HasPrefix(cleaned, filepath.Join(projectRoot, dir)) {
		return os.ErrPermission
	}
	return os.WriteFile(cleaned, []byte(content), 0644)
}

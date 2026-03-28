package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/saintfish/chardet"
	"github.com/suapapa/go_hangul/encoding/cp949"
	subtitle "github.com/suapapa/go_subtitle"
)

// convertSmiToSrt reads SMI from r and writes SRT to w.
func convertSmiToSrt(r io.Reader, w io.Writer) error {
	// Detect charset.
	buff, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	det := chardet.NewHtmlDetector()
	d, err := det.DetectBest(buff)
	if err != nil {
		// Fallback or error?
		return fmt.Errorf("failed to detect charset: %w", err)
	}

	// CP949 or UTF-8 to UTF-8.
	var utf8Buff []byte
	switch d.Charset {
	case "EUC-KR":
		utf8Buff, err = cp949.From(buff)
		if err != nil {
			return fmt.Errorf("failed to convert from CP949: %w", err)
		}
	case "UTF-8":
		utf8Buff = buff
	default:
		// Try UTF-8 anyway if it looks like it? Or just error.
		utf8Buff = buff
	}

	utf8Src := bytes.NewReader(utf8Buff)

	// Convert it to subtitle.Book.
	b, err := subtitle.ReadSmi(utf8Src)
	if err != nil {
		return fmt.Errorf("failed to read SMI: %w", err)
	}

	// Export to SRT.
	if err := subtitle.ExportToSrtFile(b, w); err != nil {
		return fmt.Errorf("failed to export SRT: %w", err)
	}

	return nil
}

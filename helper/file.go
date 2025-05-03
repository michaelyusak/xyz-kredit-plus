package helper

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/michaelyusak/xyz-kredit-plus/entity"
)

var allowedMimeTypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"application/pdf": ".pdf",
	"text/plain":      ".txt",
}

// Naive check for embedded script content
func hasMaliciousContent(data []byte) bool {
	content := strings.ToLower(string(data))

	return strings.Contains(content, "<script>") || strings.Contains(content, "<?php") || strings.Contains(content, "eval(")
}

// ValidateFile checks the MIME type and optionally scans for embedded scripts. Takes bytes, returns extension and error.
func ValidateFileBytes(data []byte, allowedExt []string) (string, error) {
	contentType := http.DetectContentType(data)

	ext, ok := allowedMimeTypes[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	if !slices.Contains(allowedExt, ext) {
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	if strings.HasPrefix(contentType, "text/") || contentType == "application/pdf" {
		if hasMaliciousContent(data) {
			return "", fmt.Errorf("file contains potentially malicious content")
		}
	}

	return ext, nil
}

// ValidateFile checks the MIME type and optionally scans for embedded scripts. Takes multipart, returns extension and error.
func ValidateFileMultipart(media entity.Media, allowedExt []string) (string, error) {
	buf := make([]byte, 512)

	n, err := media.File.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	contentType := http.DetectContentType(buf[:n])

	ext, ok := allowedMimeTypes[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported file type: %s", contentType)
	}

	if !slices.Contains(allowedExt, ext) {
		return "", fmt.Errorf("file extension '%s' is not allowed", ext)
	}

	if _, err := media.File.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("failed to rewind file: %w", err)
	}

	if strings.HasPrefix(contentType, "text/") || contentType == "application/pdf" {
		data, err := io.ReadAll(media.File)
		if err != nil {
			return "", fmt.Errorf("failed to read file for scan: %w", err)
		}

		if hasMaliciousContent(data) {
			return "", fmt.Errorf("file contains potentially malicious content")
		}
	}

	return ext, nil
}

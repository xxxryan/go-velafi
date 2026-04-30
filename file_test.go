package velafi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-id.jpg")
	if err := os.WriteFile(testFile, []byte("fake-image-data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/files" {
			t.Errorf("path = %q, want /v1/files", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("Content-Type = %q, want multipart/form-data", r.Header.Get("Content-Type"))
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatalf("ParseMultipartForm error: %v", err)
		}
		if r.FormValue("purpose") != "id_document" {
			t.Errorf("purpose = %q, want id_document", r.FormValue("purpose"))
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile error: %v", err)
		}
		defer file.Close()
		if header.Filename != "test-id.jpg" {
			t.Errorf("filename = %q, want test-id.jpg", header.Filename)
		}
		data, _ := io.ReadAll(file)
		if string(data) != "fake-image-data" {
			t.Errorf("file content = %q", string(data))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data": map[string]any{
				"fileId":    "file-abc-123",
				"filename":  "test-id.jpg",
				"purpose":   "id_document",
				"size":      15,
				"createdAt": "2026-04-30T10:00:00Z",
			},
		})
	})

	result, err := c.UploadFile(context.Background(), &UploadFileParams{
		FilePath: testFile,
		Purpose:  "id_document",
	})
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
	if result.FileID != "file-abc-123" {
		t.Errorf("FileID = %q, want %q", result.FileID, "file-abc-123")
	}
}

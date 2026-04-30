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
	testFile1 := filepath.Join(tmpDir, "id-front.jpg")
	testFile2 := filepath.Join(tmpDir, "id-back.jpg")
	if err := os.WriteFile(testFile1, []byte("fake-front-data"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile2, []byte("fake-back-data"), 0644); err != nil {
		t.Fatal(err)
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/file/upload" {
			t.Errorf("path = %q, want /v2/base/file/upload", r.URL.Path)
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
		if r.FormValue("businessType") != "kyc" {
			t.Errorf("businessType = %q, want kyc", r.FormValue("businessType"))
		}

		files := r.MultipartForm.File["files"]
		if len(files) != 2 {
			t.Fatalf("len(files) = %d, want 2", len(files))
		}
		if files[0].Filename != "id-front.jpg" {
			t.Errorf("files[0].Filename = %q, want id-front.jpg", files[0].Filename)
		}
		if files[1].Filename != "id-back.jpg" {
			t.Errorf("files[1].Filename = %q, want id-back.jpg", files[1].Filename)
		}

		f, err := files[0].Open()
		if err != nil {
			t.Fatalf("open file error: %v", err)
		}
		defer f.Close()
		data, _ := io.ReadAll(f)
		if string(data) != "fake-front-data" {
			t.Errorf("file content = %q, want fake-front-data", string(data))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": []map[string]any{
				{
					"fileName":    "id-front.jpg",
					"fileType":    "jpg",
					"fileUrl":     "https://cdn.velafi.com/id-front.jpg",
					"tempFileUrl": "https://cdn.velafi.com/tmp/id-front.jpg",
				},
				{
					"fileName":    "id-back.jpg",
					"fileType":    "jpg",
					"fileUrl":     "https://cdn.velafi.com/id-back.jpg",
					"tempFileUrl": "https://cdn.velafi.com/tmp/id-back.jpg",
				},
			},
		})
	})

	result, err := c.UploadFile(context.Background(), &UploadFileParams{
		BusinessType: "kyc",
		FilePaths:    []string{testFile1, testFile2},
	})
	if err != nil {
		t.Fatalf("UploadFile() error = %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("len(result) = %d, want 2", len(result))
	}
	if result[0].FileName != "id-front.jpg" {
		t.Errorf("FileName = %q, want %q", result[0].FileName, "id-front.jpg")
	}
	if result[0].FileURL != "https://cdn.velafi.com/id-front.jpg" {
		t.Errorf("FileURL = %q, want %q", result[0].FileURL, "https://cdn.velafi.com/id-front.jpg")
	}
	if result[1].FileName != "id-back.jpg" {
		t.Errorf("FileName = %q, want %q", result[1].FileName, "id-back.jpg")
	}
}

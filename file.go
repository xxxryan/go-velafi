package velafi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (c *Client) UploadFile(ctx context.Context, params *UploadFileParams) (*File, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	f, err := os.Open(params.FilePath)
	if err != nil {
		return nil, fmt.Errorf("velafi: open file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(params.FilePath))
	if err != nil {
		return nil, fmt.Errorf("velafi: create form file: %w", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, fmt.Errorf("velafi: copy file data: %w", err)
	}

	if err := writer.WriteField("purpose", params.Purpose); err != nil {
		return nil, fmt.Errorf("velafi: write purpose field: %w", err)
	}
	writer.Close()

	url := c.baseURL + "/v1/files"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("velafi: create upload request: %w", err)
	}
	req.Header.Set("X-BH-TOKEN", c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("velafi: read upload response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("velafi: unmarshal upload response: %w", err)
	}

	if resp.StatusCode >= 400 || apiResp.Code != 0 {
		return nil, &Error{
			HTTPStatus: resp.StatusCode,
			Code:       apiResp.Code,
			Message:    apiResp.Message,
		}
	}

	var result File
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("velafi: unmarshal file data: %w", err)
	}
	return &result, nil
}

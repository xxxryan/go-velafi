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

func (c *Client) UploadFile(ctx context.Context, params *UploadFileParams) ([]UploadedFile, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("businessType", params.BusinessType); err != nil {
		return nil, fmt.Errorf("velafi: write businessType field: %w", err)
	}

	for _, fp := range params.FilePaths {
		f, err := os.Open(fp)
		if err != nil {
			return nil, fmt.Errorf("velafi: open file %s: %w", fp, err)
		}
		part, err := writer.CreateFormFile("files", filepath.Base(fp))
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("velafi: create form file: %w", err)
		}
		if _, err := io.Copy(part, f); err != nil {
			f.Close()
			return nil, fmt.Errorf("velafi: copy file data: %w", err)
		}
		f.Close()
	}
	writer.Close()

	reqURL := c.baseURL + "/v2/base/file/upload"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, &buf)
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

	if resp.StatusCode >= 400 || apiResp.Code != 200 {
		return nil, &Error{
			HTTPStatus: resp.StatusCode,
			Code:       apiResp.Code,
			Message:    apiResp.Msg,
		}
	}

	var result []UploadedFile
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("velafi: unmarshal file data: %w", err)
	}
	return result, nil
}

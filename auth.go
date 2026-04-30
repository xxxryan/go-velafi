package velafi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type tokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Token      string `json:"token"`
		ExpireTime string `json:"expireTime"`
	} `json:"data"`
}

func (c *Client) generateSignature() (timestamp string, signature string) {
	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	message := "timestamp=" + ts
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(message))
	sig := hex.EncodeToString(mac.Sum(nil))
	return ts, sig
}

func (c *Client) ensureToken(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Add(c.tokenRefreshBuffer).Before(c.expiresAt) {
		return nil
	}

	return c.refreshToken(ctx)
}

func (c *Client) refreshToken(ctx context.Context) error {
	ts, sig := c.generateSignature()

	url := c.baseURL + "/v2/token/generate?timestamp=" + ts + "&signature=" + sig
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("velafi: create token request: %w", err)
	}
	req.Header.Set("X-BH-APIKEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("velafi: token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("velafi: read token response: %w", err)
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("velafi: unmarshal token response: %w", err)
	}

	if tokenResp.Code != 200 {
		return &Error{
			HTTPStatus: resp.StatusCode,
			Code:       tokenResp.Code,
			Message:    tokenResp.Msg,
		}
	}

	expireMs, err := strconv.ParseInt(tokenResp.Data.ExpireTime, 10, 64)
	if err != nil {
		return fmt.Errorf("velafi: parse expireTime %q: %w", tokenResp.Data.ExpireTime, err)
	}

	c.token = tokenResp.Data.Token
	c.expiresAt = time.UnixMilli(expireMs)

	return nil
}

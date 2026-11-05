package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxDownloadSize = 25 << 20 // 25 MB

// Client interacts with the Telegram Bot HTTP API.
type Client struct {
	BotToken string
	HTTP     *http.Client
}

// NewClient returns a new Telegram API client.
func NewClient(token string) *Client {
	return &Client{
		BotToken: token,
		HTTP:     &http.Client{Timeout: 60 * time.Second},
	}
}

// GetFilePath resolves a Telegram file_id to its relative file_path.
func (c *Client) GetFilePath(ctx context.Context, fileID string) (string, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", c.BotToken, fileID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return "", fmt.Errorf("call getFile: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("getFile returned status %d: %s", resp.StatusCode, string(body))
	}

	var res FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("decode getFile response: %w", err)
	}

	if !res.OK || res.Result == nil || res.Result.FilePath == "" {
		return "", fmt.Errorf("telegram getFile error: %s", res.Description)
	}

	return res.Result.FilePath, nil
}

// Download fetches the raw binary file from Telegram file servers.
func (c *Client) Download(ctx context.Context, filePath string) ([]byte, error) {
	url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", c.BotToken, filePath)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("download file stream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxDownloadSize+1))
	if err != nil {
		return nil, fmt.Errorf("read file data: %w", err)
	}

	if int64(len(data)) > maxDownloadSize {
		return nil, fmt.Errorf("file exceeds maximum allowed size of %d bytes", maxDownloadSize)
	}

	return data, nil
}

// SendMessage sends an acknowledgment message back to the sender on Telegram.
func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.BotToken)
	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}


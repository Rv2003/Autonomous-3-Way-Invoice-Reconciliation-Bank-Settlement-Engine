package whatsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maxFileBytes = 25 << 20 // 25 MB limit

// MediaClient downloads files from the WhatsApp Graph API.
type MediaClient struct {
	BaseURL string // default: "https://graph.facebook.com/v21.0"
	Token   string // WA_ACCESS_TOKEN
	HTTP    *http.Client
}

// MediaInfo represents the metadata returned by Meta when querying a media ID.
type MediaInfo struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	FileSize int64  `json:"file_size"`
}

// Download resolves a media ID to its temporary URL and downloads the binary payload.
func (c *MediaClient) Download(ctx context.Context, mediaID string) (MediaInfo, []byte, error) {
	var info MediaInfo

	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}

	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = "https://graph.facebook.com/v21.0"
	}

	resp, err := c.get(ctx, baseURL+"/"+mediaID)
	if err != nil {
		return info, nil, fmt.Errorf("fetch media info: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return info, nil, fmt.Errorf("decode media info: %w", err)
	}

	if info.URL == "" {
		return info, nil, fmt.Errorf("media info did not contain a download URL")
	}

	fileResp, err := c.get(ctx, info.URL)
	if err != nil {
		return info, nil, fmt.Errorf("download file stream: %w", err)
	}
	defer fileResp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(fileResp.Body, maxFileBytes+1))
	if err != nil {
		return info, nil, fmt.Errorf("read media bytes: %w", err)
	}
	if int64(len(data)) > maxFileBytes {
		return info, nil, fmt.Errorf("media %s exceeds maximum allowed limit of %d bytes", mediaID, maxFileBytes)
	}

	return info, data, nil
}

func (c *MediaClient) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s returned status %d", url, resp.StatusCode)
	}

	return resp, nil
}

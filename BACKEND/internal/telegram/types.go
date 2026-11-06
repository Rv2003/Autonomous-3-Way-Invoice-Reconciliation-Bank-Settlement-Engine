package telegram

import "time"

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int64       `json:"message_id"`
	From      *User       `json:"from,omitempty"`
	Date      int64       `json:"date"`
	Document  *Document   `json:"document,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Caption   string      `json:"caption,omitempty"`
	Text      string      `json:"text,omitempty"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type Document struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id,omitempty"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type FileResponse struct {
	OK          bool        `json:"ok"`
	Result      *FileResult `json:"result,omitempty"`
	Description string      `json:"description,omitempty"`
}

type FileResult struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size,omitempty"`
}

type InboundDocument struct {
	ID                int64     `json:"id"`
	Source            string    `json:"source"`
	ExternalMessageID string    `json:"external_message_id"`
	Sender            string    `json:"sender"`
	MediaID           string    `json:"media_id"`
	MimeType          string    `json:"mime_type"`
	OriginalFilename  string    `json:"original_filename"`
	StoragePath       string    `json:"storage_path"`
	SHA256            string    `json:"sha256"`
	SizeBytes         int64     `json:"size_bytes"`
	Status            string    `json:"status"`
	ReceivedAt        time.Time `json:"received_at"`
}

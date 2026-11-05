package whatsapp

import "time"

// WebhookPayload is the body Meta POSTs to our webhook.
type WebhookPayload struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id,omitempty"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Field string      `json:"field"`
	Value ChangeValue `json:"value"`
}

type ChangeValue struct {
	MessagingProduct string    `json:"messaging_product,omitempty"`
	Metadata         Metadata  `json:"metadata,omitempty"`
	Contacts         []Contact `json:"contacts,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number,omitempty"`
	PhoneNumberID      string `json:"phone_number_id,omitempty"`
}

type Contact struct {
	Profile Profile `json:"profile,omitempty"`
	WaID    string  `json:"wa_id,omitempty"`
}

type Profile struct {
	Name string `json:"name,omitempty"`
}

type Message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp,omitempty"`
	Type      string `json:"type"`
	Document  *Media `json:"document,omitempty"`
	Image     *Media `json:"image,omitempty"`
}

// Media is the reference Meta gives us. It contains an ID, not the file itself.
type Media struct {
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	Filename string `json:"filename"`
}

// InboundDocument represents a received document recorded in the database.
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

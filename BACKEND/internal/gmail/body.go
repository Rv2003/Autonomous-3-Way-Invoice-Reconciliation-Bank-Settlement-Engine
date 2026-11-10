package gmail

import (
    "encoding/base64"
    "fmt"
    "strings"

    gmailapi "google.golang.org/api/gmail/v1"
)

func ExtractBody(message *gmailapi.Message) (string, error) {
    if message == nil {
        return "", fmt.Errorf("message must not be nil")
    }

    body := findBodyPart(message.Payload)

    if body == "" {
        return "", fmt.Errorf("message body not found")
    }

    return body, nil
}

func findBodyPart(part *gmailapi.MessagePart) string {
    if part == nil {
        return ""
    }

    if part.MimeType == "text/plain" && part.Body != nil && part.Body.Data != "" {
        return decodeBody(part.Body.Data)
    }

    for _, child := range part.Parts {
        if body := findBodyPart(child); body != "" {
            return body
        }
    }

    return ""
}

func decodeBody(data string) string {
    data = strings.ReplaceAll(data, "-", "+")
    data = strings.ReplaceAll(data, "_", "/")

    decoded, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return ""
    }

    return string(decoded)
}

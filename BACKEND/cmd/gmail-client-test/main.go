package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Rv2003/Autonomous-3-Way-Invoice-Reconciliation-Bank-Settlement-Engine/BACKEND/internal/gmail"
)

func main() {
	ctx := context.Background()

	client, err := gmail.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create Gmail client: %v", err)
	}

	messages, err := client.Service().Users.Messages.
		List("me").
		MaxResults(10).
		Do()
	if err != nil {
		log.Fatalf("failed to list Gmail messages: %v", err)
	}

	fmt.Println("Inspecting Gmail messages...")
	fmt.Println("Messages found:", len(messages.Messages))

	for i, messageRef := range messages.Messages {
		message, err := client.GetMessage(ctx, messageRef.Id)
		if err != nil {
			log.Printf(
				"Message %s: failed to retrieve: %v",
				messageRef.Id,
				err,
			)
			continue
		}

		attachments, err := gmail.ListAttachments(message)
		if err != nil {
			log.Printf(
				"Message %s: failed to inspect attachments: %v",
				message.Id,
				err,
			)
			continue
		}

		fmt.Printf(
			"\nMessage %d\n",
			i+1,
		)
		fmt.Println("  ID:", message.Id)
		fmt.Println("  Attachments:", len(attachments))
		for _, attachment := range attachments {
			fmt.Println("    Filename:", attachment.Filename)
			fmt.Println("    MIME Type:", attachment.MimeType)
			fmt.Println("    Attachment ID:", attachment.AttachmentID)
			fmt.Println("    Size:", attachment.Size)

			data, err := client.DownloadAttachment(
				ctx,
				message.Id,
				attachment.AttachmentID,
			)
			if err != nil {
				log.Fatalf(
					"failed to download attachment %q: %v",
					attachment.Filename,
					err,
				)
			}

			fmt.Println("    Downloaded Bytes:", len(data))
		}
	}
}

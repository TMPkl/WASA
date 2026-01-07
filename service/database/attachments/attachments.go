package attachments

import (
	"net/http"
)

type Attachment struct {
	Type    string
	Content []byte
}

type AttachmentsPack struct {
	Attachments []*Attachment
}

func (ap *AttachmentsPack) Quantity() uint {
	return uint(len(ap.Attachments))
}

func NewAttachment(content []byte) *Attachment {
	return &Attachment{
		Type:    http.DetectContentType(content),
		Content: content,
	}

}

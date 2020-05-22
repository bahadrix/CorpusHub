package index

import (
	"errors"
	"strings"
	"time"
)

type Document struct {
	URI          string    `json:"uri"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"last_modified"`
}

func (d *Document) Type() string {
	return "document"
}

func (d *Document) Validate() error {
	d.URI = strings.TrimSpace(d.URI)
	d.Content = strings.TrimSpace(d.Content)
	d.Title = strings.TrimSpace(d.Title)

	if len(d.URI) == 0 {
		return errors.New("URI is empty")
	}

	if len(d.Content) == 0 {
		return errors.New("content is empty")
	}

	return nil
}
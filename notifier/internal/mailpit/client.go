package mailpit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"
)

const timeout = time.Second * 10

type MessagesSummary struct {
	Total         int              `json:"total"`
	Unread        int              `json:"unread"`
	Count         int              `json:"count"`
	MessagesCount int              `json:"messages_count"`
	Start         int              `json:"start"`
	Tags          []string         `json:"tags"`
	Messages      []MessageSummary `json:"messages"`
}

type MessageSummary struct {
	ID          string          `json:"id"`
	MessageID   string          `json:"message_id"`
	Read        bool            `json:"read"`
	From        *mail.Address   `json:"from"`
	To          []*mail.Address `json:"to"`
	CC          []*mail.Address `json:"cc"`
	BCC         []*mail.Address `json:"bcc"`
	Subject     string          `json:"subject"`
	Created     time.Time       `json:"created"`
	Tags        []string        `json:"tags"`
	Size        int             `json:"size"`
	Attachments int             `json:"attachments"`
}

type Message struct {
	ID          string
	MessageID   string
	Read        bool
	From        *mail.Address
	To          []*mail.Address
	CC          []*mail.Address
	BCC         []*mail.Address
	ReplyTo     []*mail.Address
	ReturnPath  string
	Subject     string
	Date        time.Time
	Tags        []string
	Text        string
	HTML        string
	Size        int
	Inline      []Attachment
	Attachments []Attachment
}

type Attachment struct {
	PartID      string
	FileName    string
	ContentType string
	ContentID   string
	Size        int
}

// Client represents an HTTP client to Mailpit API
// which is used for retrieving messages sent via SMTP during tests.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		strings.TrimSuffix(baseURL, "/"),
		&http.Client{Timeout: timeout},
	}
}

func (c *Client) MessagesBySubject(ctx context.Context, subj string) (*MessagesSummary, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/search", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	query := url.Values{
		"query": {fmt.Sprintf("subject:%s", subj)},
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var msgs *MessagesSummary
	if err := json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
		return nil, fmt.Errorf("can't decode response body: %w", err)
	}

	return msgs, nil
}

func (c *Client) MessageByID(ctx context.Context, id string) (*Message, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/message/%s", c.baseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var msg *Message
	if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
		return nil, fmt.Errorf("can't decode response body: %w", err)
	}

	return msg, nil
}

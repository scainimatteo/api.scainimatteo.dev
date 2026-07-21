package firefly

import "time"

type FireflyWebhookResponse struct {
	UUID        string  `json:"uuid"`
	UserID      int     `json:"user_id"`
	UserGroupID int     `json:"user_group_id"`
	Trigger     string  `json:"trigger"`
	Response    string  `json:"response"`
	URL         string  `json:"url"`
	Version     string  `json:"version"`
	Content     Content `json:"content"`
}

type Content struct {
	ID           int           `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	User         int           `json:"user"`
	GroupTitle   interface{}   `json:"group_title"`
	Transactions []Transaction `json:"transactions"`
	Links        []Link        `json:"links"`
}

type Link struct {
	Rel string `json:"rel"`
	URI string `json:"uri"`
}

type Payload struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Type          string      `json:"type"`
	Date          string      `json:"date"`
	Amount        string      `json:"amount"`
	Description   string      `json:"description"`
	CategoryName  string      `json:"category_name"`
	SourceID      string      `json:"source_id,omitempty"`
	SourceName    string      `json:"source_name,omitempty"`
	DestinationID string      `json:"destination_id,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	RecurrenceID  interface{} `json:"recurrence_id,omitempty"`
}

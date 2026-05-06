package vikunja

import "time"

type VikunjaWebhookResponse struct {
	Data      Data      `json:"data"`
	EventName string    `json:"event_name"`
	Time      time.Time `json:"time"`
}

type Data struct {
	Project  Project  `json:"project"`
	Reminder Reminder `json:"reminder"` // Usato nei webhook di tipo reminder
	Task     Task     `json:"task"`
	User     User     `json:"user"`
}

type Project struct {
	BackgroundBlurHash    string      `json:"background_blur_hash"`
	BackgroundInformation interface{} `json:"background_information"`
	Created               time.Time   `json:"created"`
	Description           string      `json:"description"`
	HexColor              string      `json:"hex_color"`
	ID                    int         `json:"id"`
	Identifier            string      `json:"identifier"`
	IsArchived            bool        `json:"is_archived"`
	IsFavorite            bool        `json:"is_favorite"`
	MaxPermission         int         `json:"max_permission"`
	Owner                 *User       `json:"owner"` // Riutilizzo User
	ParentProjectID       int         `json:"parent_project_id"`
	Position              float64     `json:"position"`
	Title                 string      `json:"title"`
	Updated               time.Time   `json:"updated"`
	Views                 interface{} `json:"views"`
}

type Reminder struct {
	RelativePeriod int       `json:"relative_period"`
	RelativeTo     string    `json:"relative_to"`
	Reminder       time.Time `json:"reminder"`
}

type Label struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	HexColor    string    `json:"hex_color"`
	CreatedBy   *User     `json:"created_by"` // Riutilizzo User
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

type Task struct {
	ID                     int         `json:"id"`
	Title                  string      `json:"title"`
	Description            string      `json:"description"`
	Done                   bool        `json:"done"`
	DoneAt                 time.Time   `json:"done_at"`
	DueDate                time.Time   `json:"due_date"`
	StartDate              time.Time   `json:"start_date"`
	EndDate                time.Time   `json:"end_date"`
	Reminders              []*Reminder `json:"reminders"` // Ora è una slice di Reminder
	ProjectID              int         `json:"project_id"`
	RepeatAfter            int         `json:"repeat_after"`
	RepeatMode             int         `json:"repeat_mode"`
	Priority               int         `json:"priority"`
	Labels                 []*Label    `json:"labels"` // Ora è una slice di Label
	HexColor               string      `json:"hex_color"`
	PercentDone            float64     `json:"percent_done"`
	Identifier             string      `json:"identifier"`
	Index                  int         `json:"index"`
	IsFavorite             bool        `json:"is_favorite"`
	Created                time.Time   `json:"created"`
	Updated                time.Time   `json:"updated"`
	BucketID               int         `json:"bucket_id"`
	Position               float64     `json:"position"`
	CreatedBy              *User       `json:"created_by"`
	Assignees              interface{} `json:"assignees"`
	Attachments            interface{} `json:"attachments"`
	RelatedTasks           interface{} `json:"related_tasks"`
	Reactions              interface{} `json:"reactions"`
	CoverImageAttachmentID int         `json:"cover_image_attachment_id"`
}

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

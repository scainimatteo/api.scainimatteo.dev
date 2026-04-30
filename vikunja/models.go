package vikunja

import "time"

type VikunjaWebhookResponse struct {
	Data      Data      `json:"data"`
	EventName string    `json:"event_name"`
	Time      time.Time `json:"time"`
}

type Data struct {
	Project  Project  `json:"project"`
	Reminder Reminder `json:"reminder"`
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
	Owner                 interface{} `json:"owner"`
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

type Task struct {
	Assignees              interface{} `json:"assignees"`
	Attachments            interface{} `json:"attachments"`
	BucketID               int         `json:"bucket_id"`
	CoverImageAttachmentID int         `json:"cover_image_attachment_id"`
	Created                time.Time   `json:"created"`
	CreatedBy              interface{} `json:"created_by"`
	Description            string      `json:"description"`
	Done                   bool        `json:"done"`
	DoneAt                 time.Time   `json:"done_at"`
	DueDate                time.Time   `json:"due_date"`
	EndDate                time.Time   `json:"end_date"`
	HexColor               string      `json:"hex_color"`
	ID                     int         `json:"id"`
	Identifier             string      `json:"identifier"`
	Index                  int         `json:"index"`
	IsFavorite             bool        `json:"is_favorite"`
	Labels                 interface{} `json:"labels"`
	PercentDone            float64     `json:"percent_done"`
	Position               float64     `json:"position"`
	Priority               int         `json:"priority"`
	ProjectID              int         `json:"project_id"`
	Reactions              interface{} `json:"reactions"`
	RelatedTasks           interface{} `json:"related_tasks"`
	Reminders              interface{} `json:"reminders"`
	RepeatAfter            int         `json:"repeat_after"`
	RepeatMode             int         `json:"repeat_mode"`
	StartDate              time.Time   `json:"start_date"`
	Title                  string      `json:"title"`
	Updated                time.Time   `json:"updated"`
}

type User struct {
	Created  time.Time `json:"created"`
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Updated  time.Time `json:"updated"`
	Username string    `json:"username"`
}

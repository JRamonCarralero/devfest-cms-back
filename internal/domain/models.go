package domain

import "time"

type Speaker struct {
	ID       string `json:"id"`
	PersonID string `json:"person_id"`
	EventID  string `json:"event_id"`
	Bio      string `json:"bio"`
	Company  string `json:"company"`
	Audit
}

type Developer struct {
	ID              string `json:"id"`
	PersonID        string `json:"person_id"`
	EventID         string `json:"event_id"`
	RoleDescription string `json:"role_description"`
	Audit
}

type Collaborator struct {
	ID       string `json:"id"`
	PersonID string `json:"person_id"`
	EventID  string `json:"event_id"`
	Area     string `json:"area"`
	Audit
}

type Sponsor struct {
	ID            string `json:"id"`
	EventID       string `json:"event_id"`
	Name          string `json:"name"`
	LogoURL       string `json:"logo_url"`
	WebsiteURL    string `json:"website_url"`
	Tier          string `json:"tier"`
	OrderPriority int    `json:"order_priority"`
	Audit
}

type Talk struct {
	ID          string   `json:"id"`
	EventID     string   `json:"event_id"`
	SpeakerID   string   `json:"speaker_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Audit
}

type Schedule struct {
	ID        string     `json:"id"`
	EventID   string     `json:"event_id"`
	TalkID    *string    `json:"talk_id,omitempty"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Duration  string     `json:"duration"`
	Room      string     `json:"room"`
	Audit
}

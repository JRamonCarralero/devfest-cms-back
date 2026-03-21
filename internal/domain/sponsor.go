package domain

import "github.com/google/uuid"

type Sponsor struct {
	ID            uuid.UUID
	EventID       uuid.UUID
	Name          string
	LogoURL       string
	WebsiteURL    string
	Tier          int
	OrderPriority *int
	Audit
}

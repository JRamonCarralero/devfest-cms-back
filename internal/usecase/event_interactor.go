package usecase

import (
	"context"
	"devfest/internal/domain"
	"errors"
)

// eventInteractor implements domain.EventUsecase
type eventInteractor struct {
	repo domain.EventRepository
}

// NewEventInteractor is a constructor for eventInteractor
func NewEventInteractor(repo domain.EventRepository) domain.EventUsecase {
	return &eventInteractor{
		repo: repo,
	}
}

// GetEvents validates the parameters and calls the repository
func (i *eventInteractor) GetEvents(ctx context.Context, search string, page, pageSize int32, orderBy string) ([]domain.Event, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	validOrders := map[string]bool{
		"name_asc":        true,
		"name_desc":       true,
		"created_at_asc":  true,
		"created_at_desc": true,
	}
	if !validOrders[orderBy] {
		orderBy = "created_at_desc"
	}

	return i.repo.ListPaged(ctx, search, page, pageSize, orderBy)
}

// GetEventBySlug validates the slug and calls the repository
func (i *eventInteractor) GetEventBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	if slug == "" {
		return nil, errors.New("event slug is required")
	}

	event, err := i.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New("event not found")
	}

	return event, nil
}

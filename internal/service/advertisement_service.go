package service

import (
	"advertisement-service/internal/domain"
	"advertisement-service/internal/repository"
	"context"
)

// Service defines the advertisement service interface.
type Advertisement interface {
	Create(ctx context.Context, ad domain.Advertisement) (*domain.Advertisement, error)
	Get(ctx context.Context, id int64) (*domain.Advertisement, error)
	List(ctx context.Context, offset, limit int64, sortBy, sortDir string) ([]*domain.Advertisement, error)
	Update(ctx context.Context, id int64, ad *domain.Advertisement) (*domain.Advertisement, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

// service implements the Service interface.
type service struct {
	repo repository.AdvertisementRepository
}

// NewService creates a new advertisement service.
func NewService(repo repository.AdvertisementRepository) Advertisement {
	return &service{repo: repo}
}

// CreateAdvertisement creates a new advertisement.
func (s *service) Create(ctx context.Context, ad domain.Advertisement) (*domain.Advertisement, error) {
	// Implement the logic to create a new advertisement.
	// This usually involves calling the repository layer to insert the advertisement into the database.
	return s.repo.Create(ctx, &ad)
}

// GetAdvertisement retrieves an advertisement by its ID.
func (s *service) Get(ctx context.Context, id int64) (*domain.Advertisement, error) {
	// Implement the logic to retrieve an advertisement by its ID.
	return s.repo.FindByID(ctx, id)
}

// ListAdvertisements lists all advertisements.
func (s *service) List(ctx context.Context, offset, limit int64, sortBy, sortDir string) ([]*domain.Advertisement, error) {
	// Implement the logic to list all advertisements.
	return s.repo.Find(ctx, offset, limit, sortBy, sortDir)
}

// UpdateAdvertisement updates an existing advertisement.
func (s *service) Update(ctx context.Context, id int64, ad *domain.Advertisement) (*domain.Advertisement, error) {
	// Implement the logic to update an advertisement.
	return s.repo.Update(ctx, id, ad)
}

// DeleteAdvertisement deletes an advertisement by its ID.
func (s *service) Delete(ctx context.Context, id int64) (int64, error) {
	// Implement the logic to delete an advertisement.
	return s.repo.Delete(ctx, id)
}

package repository

import (
	"advertisement-service/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"
)

// AdvertisementRepository handles the CRUD operations with MySQL for Advertisement entity
type AdvertisementRepository struct {
	DB *sql.DB
}

// NewAdvertisementRepository creates a new AdvertisementRepository
func NewAdvertisementRepository(db *sql.DB) *AdvertisementRepository {
	return &AdvertisementRepository{DB: db}
}

// Create inserts a new Advertisement into the database
func (r *AdvertisementRepository) Create(ctx context.Context, ad *domain.Advertisement) (*domain.Advertisement, error) {
	query := `INSERT INTO advertisements (name, description, price, created_at, is_active) VALUES (?, ?, ?, ?, ?)`
	result, err := r.DB.ExecContext(ctx, query, ad.Name, ad.Description, ad.Price, time.Now(), ad.IsActive)
	if err != nil {
		return nil, err
	}
	insertedID, _ := result.LastInsertId()
	insertedAd, _ := r.FindByID(ctx, insertedID)
	return insertedAd, nil
}

// FindByID finds an Advertisement by its ID
func (r *AdvertisementRepository) FindByID(ctx context.Context, id int64) (*domain.Advertisement, error) {
	query := `SELECT id, name, description, price, created_at, is_active FROM advertisements WHERE id = ?`
	ad := domain.Advertisement{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&ad.ID, &ad.Name, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive)
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

// Find retrieves all advertisements from the database, with optional pagination and sorting
func (r *AdvertisementRepository) Find(ctx context.Context, offset, limit int64, sortBy string, sortDir string) ([]*domain.Advertisement, error) {
	query := fmt.Sprintf(`SELECT id, name, description, price, created_at, is_active FROM advertisements ORDER BY %s %s LIMIT ? OFFSET ?`, sortBy, sortDir)
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []*domain.Advertisement
	for rows.Next() {
		var ad domain.Advertisement
		if err := rows.Scan(&ad.ID, &ad.Name, &ad.Description, &ad.Price, &ad.CreatedAt, &ad.IsActive); err != nil {
			return nil, err
		}
		ads = append(ads, &ad)
	}
	return ads, nil
}

// Update modifies an existing Advertisement in the database
func (r *AdvertisementRepository) Update(ctx context.Context, id int64, ad *domain.Advertisement) (*domain.Advertisement, error) {
	query := `UPDATE advertisements SET name = ?, description = ?, price = ?, is_active = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, ad.Name, ad.Description, ad.Price, ad.IsActive, id)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

// Delete removes an Advertisement from the database
func (r *AdvertisementRepository) Delete(ctx context.Context, id int64) (int64, error) {
	query := `DELETE FROM advertisements WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, id)

	return id, err
}

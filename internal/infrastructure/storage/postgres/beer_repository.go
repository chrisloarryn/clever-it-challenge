package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"beers-challenge/internal/core/domain/beers"
	"beers-challenge/internal/core/ports/secondary"
	"beers-challenge/internal/infrastructure/config"
)

// Repository implements the secondary.BeerRepository interface
type Repository struct {
	db     *sql.DB
	config *config.ConfigProvider
}

// NewRepository creates a new PostgreSQL repository
func NewRepository(configProvider *config.ConfigProvider) (secondary.BeerRepository, error) {
	connStr := configProvider.GetDatabaseConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Repository{
		db:     db,
		config: configProvider,
	}, nil
}

// Save saves a beer to the database
func (r *Repository) Save(ctx context.Context, beer *beers.Beer) error {
	query := `
		INSERT INTO beer (id, name, brewery, country, price, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			brewery = EXCLUDED.brewery,
			country = EXCLUDED.country,
			price = EXCLUDED.price,
			currency = EXCLUDED.currency,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		beer.ID,
		beer.Name,
		beer.Brewery,
		beer.Country,
		beer.Price,
		beer.Currency,
		beer.CreatedAt,
		beer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save beer: %w", err)
	}

	return nil
}

// FindByID finds a beer by its ID
func (r *Repository) FindByID(ctx context.Context, id int) (*beers.Beer, error) {
	query := `
		SELECT id, name, brewery, country, price, currency, created_at, updated_at
		FROM beer
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var beer beers.Beer
	err := row.Scan(
		&beer.ID,
		&beer.Name,
		&beer.Brewery,
		&beer.Country,
		&beer.Price,
		&beer.Currency,
		&beer.CreatedAt,
		&beer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, beers.NewDomainError("BEER_NOT_FOUND", "Beer not found", err)
		}
		return nil, fmt.Errorf("failed to find beer: %w", err)
	}

	return &beer, nil
}

// FindAll finds all beers
func (r *Repository) FindAll(ctx context.Context) ([]beers.Beer, error) {
	query := `
		SELECT id, name, brewery, country, price, currency, created_at, updated_at
		FROM beer
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query beers: %w", err)
	}
	defer rows.Close()

	var result []beers.Beer
	for rows.Next() {
		var beer beers.Beer
		err := rows.Scan(
			&beer.ID,
			&beer.Name,
			&beer.Brewery,
			&beer.Country,
			&beer.Price,
			&beer.Currency,
			&beer.CreatedAt,
			&beer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan beer: %w", err)
		}
		result = append(result, beer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return result, nil
}

// ExistsByID checks if a beer exists by its ID
func (r *Repository) ExistsByID(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM beer WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check beer existence: %w", err)
	}

	return exists, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

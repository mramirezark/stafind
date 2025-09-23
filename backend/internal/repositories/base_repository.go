package repositories

import (
	"database/sql"
	"stafind-backend/internal/queries"
)

// BaseRepository provides common functionality for all repositories
type BaseRepository struct {
	db           *sql.DB
	queryManager *queries.YAMLQueryManager
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB) (*BaseRepository, error) {
	queryManager, err := queries.NewYAMLQueryManager()
	if err != nil {
		return nil, err
	}

	return &BaseRepository{
		db:           db,
		queryManager: queryManager,
	}, nil
}

// GetQuery safely retrieves a query by name
func (br *BaseRepository) GetQuery(name string) (string, error) {
	return br.queryManager.GetQuery(name)
}

// MustGetQuery retrieves a query by name, panics if not found
func (br *BaseRepository) MustGetQuery(name string) string {
	return br.queryManager.MustGetQuery(name)
}

// GetDB returns the database connection
func (br *BaseRepository) GetDB() *sql.DB {
	return br.db
}

// GetQueryManager returns the query manager
func (br *BaseRepository) GetQueryManager() *queries.YAMLQueryManager {
	return br.queryManager
}

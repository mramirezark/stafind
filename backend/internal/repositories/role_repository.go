package repositories

import (
	"database/sql"
	"fmt"

	"stafind-backend/internal/models"
)

// RoleRepository interface defines methods for role data access
type RoleRepository interface {
	CreateRole(role *models.Role) error
	GetRoleByID(id int) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	UpdateRole(role *models.Role) error
	DeleteRole(id int) error
	ListRoles() ([]*models.Role, error)
}

// roleRepository implements RoleRepository interface
type roleRepository struct {
	*BaseRepository
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *sql.DB) (RoleRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &roleRepository{
		BaseRepository: baseRepo,
	}, nil
}

// CreateRole creates a new role
func (r *roleRepository) CreateRole(role *models.Role) error {
	query := `
		INSERT INTO roles (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, role.Name, role.Description).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
	return err
}

// GetRoleByID retrieves a role by ID
func (r *roleRepository) GetRoleByID(id int) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, id).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByName retrieves a role by name
func (r *roleRepository) GetRoleByName(name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return role, nil
}

// UpdateRole updates an existing role
func (r *roleRepository) UpdateRole(role *models.Role) error {
	query := `
		UPDATE roles 
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.Exec(query, role.Name, role.Description, role.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// DeleteRole deletes a role by ID
func (r *roleRepository) DeleteRole(id int) error {
	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// ListRoles retrieves all roles
func (r *roleRepository) ListRoles() ([]*models.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

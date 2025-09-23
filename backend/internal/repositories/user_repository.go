package repositories

import (
	"database/sql"
	"fmt"

	"stafind-backend/internal/models"
)

// UserRepository interface defines methods for user data access
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	ListUsers(limit, offset int) ([]*models.User, error)
	GetUserCount() (int, error)
	CreateUserSession(session *models.UserSession) error
	GetUserSession(tokenHash string) (*models.UserSession, error)
	RevokeUserSession(tokenHash string) error
	CleanupExpiredSessions() error
	GetUserRoles(userID int) ([]models.Role, error)
	AssignRoleToUser(userID, roleID int) error
	RemoveRoleFromUser(userID, roleID int) error
}

// userRepository implements UserRepository interface
type userRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) (UserRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		BaseRepository: baseRepo,
	}, nil
}

// CreateUser creates a new user
func (r *userRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.RoleID,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetUserByID retrieves a user by ID
func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, 
		       u.role_id, u.is_active, u.last_login, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.created_at, r.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`

	user := &models.User{}
	var role models.Role
	var roleID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&roleID, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if roleID.Valid {
		user.RoleID = &role.ID
		user.Role = &role
	}

	// Get user roles
	roles, err := r.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, 
		       u.role_id, u.is_active, u.last_login, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.created_at, r.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`

	user := &models.User{}
	var role models.Role
	var roleID sql.NullInt64

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&roleID, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if roleID.Valid {
		user.RoleID = &role.ID
		user.Role = &role
	}

	// Get user roles
	roles, err := r.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.first_name, u.last_name, 
		       u.role_id, u.is_active, u.last_login, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.created_at, r.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = $1
	`

	user := &models.User{}
	var role models.Role
	var roleID sql.NullInt64

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&roleID, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if roleID.Valid {
		user.RoleID = &role.ID
		user.Role = &role
	}

	// Get user roles
	roles, err := r.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return user, nil
}

// UpdateUser updates an existing user
func (r *userRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, first_name = $3, last_name = $4, role_id = $5, 
		    is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	result, err := r.db.Exec(query, user.Username, user.Email, user.FirstName, user.LastName, user.RoleID, user.IsActive, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user by ID
func (r *userRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ListUsers retrieves a list of users with pagination
func (r *userRepository) ListUsers(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT u.id, u.email, u.first_name, u.last_name, u.role_id, u.is_active, 
		       u.last_login, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.created_at, r.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var role models.Role
		var roleID sql.NullInt64

		err := rows.Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName,
			&roleID, &user.IsActive, &user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
			&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if roleID.Valid {
			user.RoleID = &role.ID
			user.Role = &role
		}

		// Get user roles
		roles, err := r.GetUserRoles(user.ID)
		if err != nil {
			return nil, err
		}
		user.Roles = roles

		users = append(users, user)
	}

	return users, nil
}

// GetUserCount returns the total number of users
func (r *userRepository) GetUserCount() (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// CreateUserSession creates a new user session
func (r *userRepository) CreateUserSession(session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, session.UserID, session.TokenHash, session.ExpiresAt).Scan(&session.ID, &session.CreatedAt)
	return err
}

// GetUserSession retrieves a user session by token hash
func (r *userRepository) GetUserSession(tokenHash string) (*models.UserSession, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, is_revoked
		FROM user_sessions
		WHERE token_hash = $1 AND is_revoked = FALSE
	`

	session := &models.UserSession{}
	err := r.db.QueryRow(query, tokenHash).Scan(
		&session.ID, &session.UserID, &session.TokenHash,
		&session.ExpiresAt, &session.CreatedAt, &session.IsRevoked,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// RevokeUserSession revokes a user session
func (r *userRepository) RevokeUserSession(tokenHash string) error {
	query := `UPDATE user_sessions SET is_revoked = TRUE WHERE token_hash = $1`

	_, err := r.db.Exec(query, tokenHash)
	return err
}

// CleanupExpiredSessions removes expired sessions
func (r *userRepository) CleanupExpiredSessions() error {
	query := `DELETE FROM user_sessions WHERE expires_at < CURRENT_TIMESTAMP OR is_revoked = TRUE`

	_, err := r.db.Exec(query)
	return err
}

// GetUserRoles retrieves all roles for a user
func (r *userRepository) GetUserRoles(userID int) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// AssignRoleToUser assigns a role to a user
func (r *userRepository) AssignRoleToUser(userID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT (user_id, role_id) DO NOTHING`

	_, err := r.db.Exec(query, userID, roleID)
	return err
}

// RemoveRoleFromUser removes a role from a user
func (r *userRepository) RemoveRoleFromUser(userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	_, err := r.db.Exec(query, userID, roleID)
	return err
}

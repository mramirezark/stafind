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
	query := r.MustGetQuery("create_user")

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
	query := r.MustGetQuery("get_user_by_id")

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
	query := r.MustGetQuery("get_user_by_email")

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
	query := r.MustGetQuery("get_user_by_username")

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
	query := r.MustGetQuery("update_user")

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
	query := r.MustGetQuery("delete_user")

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
	query := r.MustGetQuery("list_users")

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
	query := r.MustGetQuery("get_user_count")

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// CreateUserSession creates a new user session
func (r *userRepository) CreateUserSession(session *models.UserSession) error {
	query := r.MustGetQuery("create_user_session")

	err := r.db.QueryRow(query, session.UserID, session.TokenHash, session.ExpiresAt).Scan(&session.ID, &session.CreatedAt)
	return err
}

// GetUserSession retrieves a user session by token hash
func (r *userRepository) GetUserSession(tokenHash string) (*models.UserSession, error) {
	query := r.MustGetQuery("get_user_session")

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
	query := r.MustGetQuery("revoke_user_session")

	_, err := r.db.Exec(query, tokenHash)
	return err
}

// CleanupExpiredSessions removes expired sessions
func (r *userRepository) CleanupExpiredSessions() error {
	query := r.MustGetQuery("cleanup_expired_sessions")

	_, err := r.db.Exec(query)
	return err
}

// GetUserRoles retrieves all roles for a user
func (r *userRepository) GetUserRoles(userID int) ([]models.Role, error) {
	query := r.MustGetQuery("get_user_roles")

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
	query := r.MustGetQuery("assign_role_to_user")

	_, err := r.db.Exec(query, userID, roleID)
	return err
}

// RemoveRoleFromUser removes a role from a user
func (r *userRepository) RemoveRoleFromUser(userID, roleID int) error {
	query := r.MustGetQuery("remove_role_from_user")

	_, err := r.db.Exec(query, userID, roleID)
	return err
}

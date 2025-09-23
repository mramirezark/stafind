package services

import (
	"fmt"
	"time"

	"stafind-backend/internal/auth"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

// UserService interface defines business logic for user operations
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(id int) (*models.UserResponse, error)
	GetUserByEmail(email string) (*models.UserResponse, error)
	UpdateUser(id int, req *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(id int) error
	ListUsers(page, limit int) ([]*models.UserResponse, int, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(tokenHash string) error
	ChangePassword(userID int, req *models.ChangePasswordRequest) error
	RefreshToken(token string) (string, error)
	GetUserProfile(userID int) (*models.UserResponse, error)
	UpdateUserProfile(userID int, req *models.UpdateUserRequest) (*models.UserResponse, error)
}

// userService implements UserService interface
type userService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, NewConflictError("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		RoleID:       req.RoleID,
		IsActive:     true,
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the created user with roles
	createdUser, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created user: %w", err)
	}

	response := createdUser.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, NewNotFoundError("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByEmail retrieves a user by email
func (s *userService) GetUserByEmail(email string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, NewNotFoundError("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, NewNotFoundError("user not found")
	}

	// Update fields if provided
	if req.Username != nil {
		// Check if username is already taken by another user
		existingUser, err := s.userRepo.GetUserByUsername(*req.Username)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, NewConflictError("username already taken by another user")
		}
		user.Username = *req.Username
	}

	if req.Email != nil {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetUserByEmail(*req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, NewConflictError("email already taken by another user")
		}
		user.Email = *req.Email
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if req.RoleID != nil {
		user.RoleID = req.RoleID
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update user in database
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	response := updatedUser.ToResponse()
	return &response, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id int) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return NewNotFoundError("user not found")
	}

	err = s.userRepo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users with pagination
func (s *userService) ListUsers(page, limit int) ([]*models.UserResponse, int, error) {
	offset := (page - 1) * limit

	users, err := s.userRepo.ListUsers(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.userRepo.GetUserCount()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user count: %w", err)
	}

	var responses []*models.UserResponse
	for _, user := range users {
		response := user.ToResponse()
		responses = append(responses, &response)
	}

	return responses, total, nil
}

// Login authenticates a user and returns a JWT token
func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username only
	user, err := s.userRepo.GetUserByUsername(req.Username)

	if err != nil {
		return nil, NewValidationError("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, NewValidationError("account is deactivated")
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, NewValidationError("invalid username or password")
	}

	// Generate JWT token
	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	token, err := auth.GenerateJWT(user.ID, user.Email, user.FirstName, user.LastName, roleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session
	tokenHash := auth.GenerateTokenHash(token)
	session := &models.UserSession{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.userRepo.CreateUserSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login: %v\n", err)
	}

	response := &models.LoginResponse{
		User:  *user,
		Token: token,
	}

	return response, nil
}

// Logout revokes a user session
func (s *userService) Logout(tokenHash string) error {
	err := s.userRepo.RevokeUserSession(tokenHash)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

// ChangePassword changes a user's password
func (s *userService) ChangePassword(userID int, req *models.ChangePasswordRequest) error {
	// Get user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return NewNotFoundError("user not found")
	}

	// Verify current password
	if !auth.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		return NewValidationError("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = hashedPassword
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// RefreshToken generates a new token with extended expiration
func (s *userService) RefreshToken(token string) (string, error) {
	// Validate current token
	claims, err := auth.ValidateJWT(token)
	if err != nil {
		return "", NewValidationError("invalid token")
	}

	// Generate new token
	newToken, err := auth.RefreshToken(token)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	// Create new session
	tokenHash := auth.GenerateTokenHash(newToken)
	session := &models.UserSession{
		UserID:    claims.UserID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.userRepo.CreateUserSession(session)
	if err != nil {
		return "", fmt.Errorf("failed to create new session: %w", err)
	}

	return newToken, nil
}

// GetUserProfile retrieves a user's profile
func (s *userService) GetUserProfile(userID int) (*models.UserResponse, error) {
	return s.GetUserByID(userID)
}

// UpdateUserProfile updates a user's profile
func (s *userService) UpdateUserProfile(userID int, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Remove role_id from update request for profile updates
	if req.RoleID != nil {
		req.RoleID = nil // Users cannot change their own role
	}

	return s.UpdateUser(userID, req)
}

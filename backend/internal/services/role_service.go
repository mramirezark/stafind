package services

import (
	"fmt"

	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

// RoleService interface defines business logic for role operations
type RoleService interface {
	CreateRole(role *models.Role) (*models.Role, error)
	GetRoleByID(id int) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	UpdateRole(id int, role *models.Role) (*models.Role, error)
	DeleteRole(id int) error
	ListRoles() ([]*models.Role, error)
}

// roleService implements RoleService interface
type roleService struct {
	roleRepo repositories.RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(roleRepo repositories.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}

// CreateRole creates a new role
func (s *roleService) CreateRole(role *models.Role) (*models.Role, error) {
	// Check if role already exists
	existingRole, err := s.roleRepo.GetRoleByName(role.Name)
	if err == nil && existingRole != nil {
		return nil, NewConflictError("role with this name already exists")
	}

	err = s.roleRepo.CreateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Get the created role
	createdRole, err := s.roleRepo.GetRoleByID(role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created role: %w", err)
	}

	return createdRole, nil
}

// GetRoleByID retrieves a role by ID
func (s *roleService) GetRoleByID(id int) (*models.Role, error) {
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, NewNotFoundError("role not found")
	}

	return role, nil
}

// GetRoleByName retrieves a role by name
func (s *roleService) GetRoleByName(name string) (*models.Role, error) {
	role, err := s.roleRepo.GetRoleByName(name)
	if err != nil {
		return nil, NewNotFoundError("role not found")
	}

	return role, nil
}

// UpdateRole updates an existing role
func (s *roleService) UpdateRole(id int, role *models.Role) (*models.Role, error) {
	// Check if role exists
	_, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, NewNotFoundError("role not found")
	}

	// Check if name is already taken by another role
	existingRole, err := s.roleRepo.GetRoleByName(role.Name)
	if err == nil && existingRole != nil && existingRole.ID != id {
		return nil, NewConflictError("role name already taken")
	}

	// Update role
	role.ID = id
	err = s.roleRepo.UpdateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Get updated role
	updatedRole, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated role: %w", err)
	}

	return updatedRole, nil
}

// DeleteRole deletes a role
func (s *roleService) DeleteRole(id int) error {
	// Check if role exists
	_, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return NewNotFoundError("role not found")
	}

	err = s.roleRepo.DeleteRole(id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// ListRoles retrieves all roles
func (s *roleService) ListRoles() ([]*models.Role, error) {
	roles, err := s.roleRepo.ListRoles()
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	return roles, nil
}

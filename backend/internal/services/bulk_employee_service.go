package services

import (
	"encoding/json"
	"fmt"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"time"
)

// BulkEmployeeService handles bulk operations for employees
type BulkEmployeeService interface {
	BulkCreateEmployees(employees []models.Employee) (*BulkOperationResult, error)
	BulkUpdateEmployees(employees []models.Employee) (*BulkOperationResult, error)
	BulkDeleteEmployees(employeeIDs []int) (*BulkOperationResult, error)
	BulkUpsertEmployees(employees []models.Employee) (*BulkOperationResult, error)
}

type bulkEmployeeService struct {
	employeeRepo repositories.EmployeeRepository
	skillRepo    repositories.SkillRepository
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	TotalProcessed int                  `json:"total_processed"`
	Successful     int                  `json:"successful"`
	Failed         int                  `json:"failed"`
	Errors         []BulkOperationError `json:"errors"`
	CreatedAt      time.Time            `json:"created_at"`
	Duration       time.Duration        `json:"duration"`
}

// BulkOperationError represents an error in a bulk operation
type BulkOperationError struct {
	Index     int    `json:"index"`
	Employee  string `json:"employee"`
	Field     string `json:"field,omitempty"`
	Message   string `json:"message"`
	ErrorType string `json:"error_type"`
}

// NewBulkEmployeeService creates a new bulk employee service
func NewBulkEmployeeService(employeeRepo repositories.EmployeeRepository, skillRepo repositories.SkillRepository) BulkEmployeeService {
	return &bulkEmployeeService{
		employeeRepo: employeeRepo,
		skillRepo:    skillRepo,
	}
}

// BulkCreateEmployees creates multiple employees in a single operation
func (s *bulkEmployeeService) BulkCreateEmployees(employees []models.Employee) (*BulkOperationResult, error) {
	startTime := time.Now()
	result := &BulkOperationResult{
		TotalProcessed: len(employees),
		Errors:         []BulkOperationError{},
		CreatedAt:      startTime,
	}

	// Validate all employees first
	validEmployees := []models.Employee{}
	for i, employee := range employees {
		if err := s.validateEmployee(employee); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  employee.Name,
				Message:   err.Error(),
				ErrorType: "validation",
			})
			continue
		}
		validEmployees = append(validEmployees, employee)
	}

	// Process valid employees
	for i, employee := range validEmployees {
		if err := s.createEmployeeWithSkills(employee); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  employee.Name,
				Message:   err.Error(),
				ErrorType: "creation",
			})
		} else {
			result.Successful++
		}
	}

	result.Failed = len(result.Errors)
	result.Duration = time.Since(startTime)
	return result, nil
}

// BulkUpdateEmployees updates multiple employees
func (s *bulkEmployeeService) BulkUpdateEmployees(employees []models.Employee) (*BulkOperationResult, error) {
	startTime := time.Now()
	result := &BulkOperationResult{
		TotalProcessed: len(employees),
		Errors:         []BulkOperationError{},
		CreatedAt:      startTime,
	}

	for i, employee := range employees {
		if err := s.validateEmployee(employee); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  employee.Name,
				Message:   err.Error(),
				ErrorType: "validation",
			})
			continue
		}

		if err := s.updateEmployeeWithSkills(employee); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  employee.Name,
				Message:   err.Error(),
				ErrorType: "update",
			})
		} else {
			result.Successful++
		}
	}

	result.Failed = len(result.Errors)
	result.Duration = time.Since(startTime)
	return result, nil
}

// BulkDeleteEmployees deletes multiple employees
func (s *bulkEmployeeService) BulkDeleteEmployees(employeeIDs []int) (*BulkOperationResult, error) {
	startTime := time.Now()
	result := &BulkOperationResult{
		TotalProcessed: len(employeeIDs),
		Errors:         []BulkOperationError{},
		CreatedAt:      startTime,
	}

	for i, id := range employeeIDs {
		if err := s.employeeRepo.Delete(id); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  fmt.Sprintf("ID: %d", id),
				Message:   err.Error(),
				ErrorType: "deletion",
			})
		} else {
			result.Successful++
		}
	}

	result.Failed = len(result.Errors)
	result.Duration = time.Since(startTime)
	return result, nil
}

// BulkUpsertEmployees creates or updates employees based on email
func (s *bulkEmployeeService) BulkUpsertEmployees(employees []models.Employee) (*BulkOperationResult, error) {
	startTime := time.Now()
	result := &BulkOperationResult{
		TotalProcessed: len(employees),
		Errors:         []BulkOperationError{},
		CreatedAt:      startTime,
	}

	for i, employee := range employees {
		if err := s.validateEmployee(employee); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  employee.Name,
				Message:   err.Error(),
				ErrorType: "validation",
			})
			continue
		}

		// For now, we'll always create new employees
		// TODO: Implement email-based duplicate checking
		if err := s.createEmployeeWithSkills(employee); err != nil {
			result.Errors = append(result.Errors, BulkOperationError{
				Index:     i,
				Employee:  employee.Name,
				Message:   err.Error(),
				ErrorType: "creation",
			})
		} else {
			result.Successful++
		}
	}

	result.Failed = len(result.Errors)
	result.Duration = time.Since(startTime)
	return result, nil
}

// validateEmployee validates an employee before processing
func (s *bulkEmployeeService) validateEmployee(employee models.Employee) error {
	if employee.Name == "" {
		return fmt.Errorf("name is required")
	}
	if employee.Email == "" {
		return fmt.Errorf("email is required")
	}
	if employee.Department == "" {
		return fmt.Errorf("department is required")
	}
	if employee.Level == "" {
		return fmt.Errorf("level is required")
	}
	return nil
}

// createEmployeeWithSkills creates an employee and their skills
func (s *bulkEmployeeService) createEmployeeWithSkills(employee models.Employee) error {
	// Create the employee
	createReq := &models.CreateEmployeeRequest{
		Name:       employee.Name,
		Email:      employee.Email,
		Department: employee.Department,
		Level:      employee.Level,
		Location:   employee.Location,
		Bio:        employee.Bio,
	}
	createdEmployee, err := s.employeeRepo.Create(createReq)
	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	// Add skills if provided
	if len(employee.Skills) > 0 {
		for _, skill := range employee.Skills {
			// Check if skill exists, create if not
			existingSkill, err := s.skillRepo.GetByName(skill.Name)
			if err != nil {
				// Skill doesn't exist, create it
				newSkill := &models.Skill{
					Name:     skill.Name,
					Category: skill.Category,
				}
				createdSkill, err := s.skillRepo.Create(newSkill)
				if err != nil {
					return fmt.Errorf("failed to create skill %s: %w", skill.Name, err)
				}
				skill.ID = createdSkill.ID
			} else {
				skill.ID = existingSkill.ID
			}

			// Add employee skill relationship
			employeeSkillReq := &models.EmployeeSkillReq{
				SkillName:        skill.Name,
				ProficiencyLevel: 3,   // Default proficiency level
				YearsExperience:  1.0, // Default years experience
			}
			if err := s.employeeRepo.AddSkill(createdEmployee.ID, employeeSkillReq); err != nil {
				return fmt.Errorf("failed to add skill %s: %w", skill.Name, err)
			}
		}
	}

	return nil
}

// updateEmployeeWithSkills updates an employee and their skills
func (s *bulkEmployeeService) updateEmployeeWithSkills(employee models.Employee) error {
	// Update the employee
	updateReq := &models.CreateEmployeeRequest{
		Name:       employee.Name,
		Email:      employee.Email,
		Department: employee.Department,
		Level:      employee.Level,
		Location:   employee.Location,
		Bio:        employee.Bio,
	}
	_, err := s.employeeRepo.Update(employee.ID, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	// Update skills if provided
	if len(employee.Skills) > 0 {
		// Remove existing skills
		if err := s.employeeRepo.RemoveSkills(employee.ID); err != nil {
			return fmt.Errorf("failed to remove existing skills: %w", err)
		}

		// Add new skills
		for _, skill := range employee.Skills {
			// Check if skill exists, create if not
			existingSkill, err := s.skillRepo.GetByName(skill.Name)
			if err != nil {
				// Skill doesn't exist, create it
				newSkill := &models.Skill{
					Name:     skill.Name,
					Category: skill.Category,
				}
				createdSkill, err := s.skillRepo.Create(newSkill)
				if err != nil {
					return fmt.Errorf("failed to create skill %s: %w", skill.Name, err)
				}
				skill.ID = createdSkill.ID
			} else {
				skill.ID = existingSkill.ID
			}

			// Add employee skill relationship
			employeeSkillReq := &models.EmployeeSkillReq{
				SkillName:        skill.Name,
				ProficiencyLevel: 3,   // Default proficiency level
				YearsExperience:  1.0, // Default years experience
			}
			if err := s.employeeRepo.AddSkill(employee.ID, employeeSkillReq); err != nil {
				return fmt.Errorf("failed to add skill %s: %w", skill.Name, err)
			}
		}
	}

	return nil
}

// ParseResumeData parses resume data from JSON
func ParseResumeData(jsonData []byte) ([]models.Employee, error) {
	var employees []models.Employee
	if err := json.Unmarshal(jsonData, &employees); err != nil {
		return nil, fmt.Errorf("failed to parse resume data: %w", err)
	}
	return employees, nil
}

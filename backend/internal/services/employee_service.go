package services

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type employeeService struct {
	employeeRepo repositories.EmployeeRepository
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(employeeRepo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
	}
}

func (s *employeeService) GetAllEmployees() ([]models.Employee, error) {
	return s.employeeRepo.GetAll()
}

func (s *employeeService) GetEmployeeByID(id int) (*models.Employee, error) {
	return s.employeeRepo.GetByID(id)
}

func (s *employeeService) CreateEmployee(req *models.CreateEmployeeRequest) (*models.Employee, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, &ValidationError{Field: "name", Message: "Name is required"}
	}
	if req.Email == "" {
		return nil, &ValidationError{Field: "email", Message: "Email is required"}
	}

	// Validate skills
	for i, skill := range req.Skills {
		if skill.SkillName == "" {
			return nil, &ValidationError{Field: "skills", Message: "Skill name is required"}
		}
		if skill.ProficiencyLevel < 1 || skill.ProficiencyLevel > 5 {
			return nil, &ValidationError{
				Field:   "skills",
				Message: "Proficiency level must be between 1 and 5",
			}
		}
		if skill.YearsExperience < 0 {
			return nil, &ValidationError{
				Field:   "skills",
				Message: "Years of experience cannot be negative",
			}
		}
		if skill.YearsExperience > 50 {
			return nil, &ValidationError{
				Field:   "skills",
				Message: "Years of experience seems unrealistic",
			}
		}
		// Check for duplicate skills
		for j := i + 1; j < len(req.Skills); j++ {
			if req.Skills[j].SkillName == skill.SkillName {
				return nil, &ValidationError{
					Field:   "skills",
					Message: "Duplicate skill: " + skill.SkillName,
				}
			}
		}
	}

	return s.employeeRepo.Create(req)
}

func (s *employeeService) UpdateEmployee(id int, req *models.CreateEmployeeRequest) (*models.Employee, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, &ValidationError{Field: "name", Message: "Name is required"}
	}
	if req.Email == "" {
		return nil, &ValidationError{Field: "email", Message: "Email is required"}
	}

	// Validate skills
	for i, skill := range req.Skills {
		if skill.SkillName == "" {
			return nil, &ValidationError{Field: "skills", Message: "Skill name is required"}
		}
		if skill.ProficiencyLevel < 1 || skill.ProficiencyLevel > 5 {
			return nil, &ValidationError{
				Field:   "skills",
				Message: "Proficiency level must be between 1 and 5",
			}
		}
		if skill.YearsExperience < 0 {
			return nil, &ValidationError{
				Field:   "skills",
				Message: "Years of experience cannot be negative",
			}
		}
		if skill.YearsExperience > 50 {
			return nil, &ValidationError{
				Field:   "skills",
				Message: "Years of experience seems unrealistic",
			}
		}
		// Check for duplicate skills
		for j := i + 1; j < len(req.Skills); j++ {
			if req.Skills[j].SkillName == skill.SkillName {
				return nil, &ValidationError{
					Field:   "skills",
					Message: "Duplicate skill: " + skill.SkillName,
				}
			}
		}
	}

	return s.employeeRepo.Update(id, req)
}

func (s *employeeService) DeleteEmployee(id int) error {
	return s.employeeRepo.Delete(id)
}

package services

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type skillService struct {
	skillRepo repositories.SkillRepository
}

// NewSkillService creates a new skill service
func NewSkillService(skillRepo repositories.SkillRepository) SkillService {
	return &skillService{
		skillRepo: skillRepo,
	}
}

func (s *skillService) GetAllSkills() ([]models.Skill, error) {
	return s.skillRepo.GetAll()
}

func (s *skillService) GetSkillByID(id int) (*models.Skill, error) {
	return s.skillRepo.GetByID(id)
}

func (s *skillService) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	// Validate required fields
	if skill.Name == "" {
		return nil, &ValidationError{Field: "name", Message: "Skill name is required"}
	}

	// Check if skill already exists
	existingSkill, err := s.skillRepo.GetByName(skill.Name)
	if err == nil && existingSkill != nil {
		return nil, &ConflictError{Resource: "skill", Message: "Skill already exists"}
	}

	return s.skillRepo.Create(skill)
}

func (s *skillService) UpdateSkill(id int, skill *models.Skill) (*models.Skill, error) {
	// Check if skill exists
	_, err := s.skillRepo.GetByID(id)
	if err != nil {
		return nil, &NotFoundError{Resource: "skill", ID: id}
	}

	// Validate required fields
	if skill.Name == "" {
		return nil, &ValidationError{Field: "name", Message: "Skill name is required"}
	}

	// Check if another skill with the same name exists
	existingSkill, err := s.skillRepo.GetByName(skill.Name)
	if err == nil && existingSkill != nil && existingSkill.ID != id {
		return nil, &ConflictError{Resource: "skill", Message: "Another skill with this name already exists"}
	}

	return s.skillRepo.Update(id, skill)
}

func (s *skillService) DeleteSkill(id int) error {
	// Check if skill exists
	_, err := s.skillRepo.GetByID(id)
	if err != nil {
		return &NotFoundError{Resource: "skill", ID: id}
	}

	return s.skillRepo.Delete(id)
}

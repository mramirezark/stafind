package services

import (
	"database/sql"
	"fmt"
	"regexp"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *categoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetCategoryByName(name string) (*models.Category, error) {
	return s.categoryRepo.GetByName(name)
}

func (s *categoryService) GetCategoriesWithSkillCount() ([]models.CategoryWithSkillCount, error) {
	return s.categoryRepo.GetCategoriesWithSkillCount()
}

func (s *categoryService) GetSkillsByCategoryID(categoryID int) ([]models.Skill, error) {
	return s.categoryRepo.GetSkillsByCategoryID(categoryID)
}

func (s *categoryService) CreateCategory(category *models.Category) (*models.Category, error) {
	// Validate category
	if err := s.ValidateCategory(category); err != nil {
		return nil, err
	}

	// Check if category already exists
	existingCategory, err := s.categoryRepo.GetByName(category.Name)
	if err == nil && existingCategory != nil {
		return nil, &ConflictError{Resource: "category", Message: "Category already exists"}
	}

	return s.categoryRepo.Create(category)
}

func (s *categoryService) CreateCategoriesBatch(categories []models.Category) ([]models.Category, error) {
	if len(categories) == 0 {
		return []models.Category{}, nil
	}

	if len(categories) > 100 {
		return nil, &ValidationError{Field: "categories", Message: "Cannot create more than 100 categories at once"}
	}

	// Validate all categories
	for i, category := range categories {
		if err := s.ValidateCategory(&category); err != nil {
			return nil, fmt.Errorf("category at index %d: %w", i, err)
		}
	}

	return s.categoryRepo.CreateBatch(categories)
}

func (s *categoryService) UpdateCategory(id int, category *models.Category) (*models.Category, error) {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Resource: "category", ID: id}
		}
		return nil, err
	}

	// Validate category
	if err := s.ValidateCategory(category); err != nil {
		return nil, err
	}

	// Check if another category with the same name exists
	existingCategory, err := s.categoryRepo.GetByName(category.Name)
	if err == nil && existingCategory != nil && existingCategory.ID != id {
		return nil, &ConflictError{Resource: "category", Message: "Another category with this name already exists"}
	}

	return s.categoryRepo.Update(id, category)
}

func (s *categoryService) DeleteCategory(id int) error {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{Resource: "category", ID: id}
		}
		return err
	}

	return s.categoryRepo.Delete(id)
}

func (s *categoryService) DeleteCategoriesBatch(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	if len(ids) > 100 {
		return &ValidationError{Field: "ids", Message: "Cannot delete more than 100 categories at once"}
	}

	// Remove duplicates and validate IDs
	uniqueIDs := make(map[int]bool)
	validIDs := []int{}
	for _, id := range ids {
		if id > 0 && !uniqueIDs[id] {
			uniqueIDs[id] = true
			validIDs = append(validIDs, id)
		}
	}

	return s.categoryRepo.DeleteBatch(validIDs)
}

func (s *categoryService) GetCategoryStats() (*models.SkillStats, error) {
	return s.categoryRepo.GetCategoryStats()
}

func (s *categoryService) ValidateCategory(category *models.Category) error {
	if category == nil {
		return &ValidationError{Field: "category", Message: "Category cannot be nil"}
	}

	if category.Name == "" {
		return &ValidationError{Field: "name", Message: "Category name is required"}
	}

	if len(category.Name) > 100 {
		return &ValidationError{Field: "name", Message: "Category name cannot exceed 100 characters"}
	}

	// Validate name format (alphanumeric, spaces, hyphens, underscores, ampersands)
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_&]+$`)
	if !nameRegex.MatchString(category.Name) {
		return &ValidationError{Field: "name", Message: "Category name contains invalid characters"}
	}

	return nil
}

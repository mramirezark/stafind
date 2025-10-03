package repositories

import (
	"database/sql"
	"fmt"
	"stafind-backend/internal/models"
	"strings"
)

type categoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sql.DB) (CategoryRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &categoryRepository{BaseRepository: baseRepo}, nil
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	query := r.MustGetQuery("get_all_categories")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) GetByID(id int) (*models.Category, error) {
	query := r.MustGetQuery("get_category_by_id")
	var category models.Category
	err := r.db.QueryRow(query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByName(name string) (*models.Category, error) {
	query := r.MustGetQuery("get_category_by_name")
	var category models.Category
	err := r.db.QueryRow(query, name).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetCategoriesWithSkillCount() ([]models.CategoryWithSkillCount, error) {
	query := r.MustGetQuery("get_categories_with_skill_count")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryWithSkillCount
	for rows.Next() {
		var category models.CategoryWithSkillCount
		err := rows.Scan(&category.ID, &category.Name, &category.SkillCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) GetSkillsByCategoryID(categoryID int) ([]models.Skill, error) {
	query := r.MustGetQuery("get_skills_by_category_id")
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var skill models.Skill
		err := rows.Scan(&skill.ID, &skill.Name)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (r *categoryRepository) Create(category *models.Category) (*models.Category, error) {
	query := r.MustGetQuery("create_category")
	var createdCategory models.Category
	err := r.db.QueryRow(query, category.Name).Scan(&createdCategory.ID, &createdCategory.Name)
	if err != nil {
		return nil, err
	}

	return &createdCategory, nil
}

func (r *categoryRepository) CreateBatch(categories []models.Category) ([]models.Category, error) {
	if len(categories) == 0 {
		return []models.Category{}, nil
	}

	// Build dynamic query for batch insert
	valueStrings := make([]string, len(categories))
	valueArgs := make([]interface{}, len(categories))

	for i, category := range categories {
		valueStrings[i] = fmt.Sprintf("($%d)", i+1)
		valueArgs[i] = category.Name
	}

	query := fmt.Sprintf("INSERT INTO categories (name) VALUES %s RETURNING id, name",
		strings.Join(valueStrings, ","))

	rows, err := r.db.Query(query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdCategories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		createdCategories = append(createdCategories, category)
	}

	return createdCategories, nil
}

func (r *categoryRepository) Update(id int, category *models.Category) (*models.Category, error) {
	query := r.MustGetQuery("update_category")
	var updatedCategory models.Category
	err := r.db.QueryRow(query, category.Name, id).Scan(&updatedCategory.ID, &updatedCategory.Name)
	if err != nil {
		return nil, err
	}

	return &updatedCategory, nil
}

func (r *categoryRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_category")
	_, err := r.db.Exec(query, id)
	return err
}

func (r *categoryRepository) DeleteBatch(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := r.MustGetQuery("delete_categories_batch")
	_, err := r.db.Exec(query, fmt.Sprintf("{%s}", strings.Join(strings.Split(fmt.Sprintf("%d", ids[0]), ""), ",")))
	return err
}

func (r *categoryRepository) GetCategoryStats() (*models.SkillStats, error) {
	query := r.MustGetQuery("get_category_stats")
	var stats models.SkillStats
	var mostPopularID sql.NullInt64
	var mostPopularName sql.NullString

	err := r.db.QueryRow(query).Scan(
		&stats.TotalSkills,
		&stats.TotalCategories,
		&mostPopularID,
		&mostPopularName,
	)
	if err != nil {
		return nil, err
	}

	// Set most popular skill if exists
	if mostPopularID.Valid && mostPopularName.Valid {
		stats.MostPopularSkill = &models.Skill{
			ID:   int(mostPopularID.Int64),
			Name: mostPopularName.String,
		}
	}

	// Get category stats
	categoryQuery := `
		SELECT 
			c.name,
			COUNT(DISTINCT sc.skill_id) as skill_count,
			COUNT(DISTINCT es.employee_id) as employee_count
		FROM categories c
		LEFT JOIN skills_categories sc ON c.id = sc.category_id
		LEFT JOIN employee_skills es ON sc.skill_id = es.skill_id
		GROUP BY c.id, c.name
		ORDER BY c.name
	`

	rows, err := r.db.Query(categoryQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryStats models.CategoryStats
		err := rows.Scan(&categoryStats.Category, &categoryStats.SkillCount, &categoryStats.EmployeeCount)
		if err != nil {
			return nil, err
		}
		stats.CategoryStats = append(stats.CategoryStats, categoryStats)
	}

	return &stats, nil
}

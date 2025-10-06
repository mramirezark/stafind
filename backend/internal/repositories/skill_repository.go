package repositories

import (
	"database/sql"
	"fmt"
	"stafind-backend/internal/models"
	"strings"
)

type skillRepository struct {
	*BaseRepository
}

// NewSkillRepository creates a new skill repository
func NewSkillRepository(db *sql.DB) (SkillRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &skillRepository{BaseRepository: baseRepo}, nil
}

func (r *skillRepository) GetAll() ([]models.Skill, error) {
	query := r.MustGetQuery("get_all_skills")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skillMap := make(map[int]*models.Skill)
	for rows.Next() {
		var skillID int
		var skillName string
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		err := rows.Scan(&skillID, &skillName, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}

		// Get or create skill
		skill, exists := skillMap[skillID]
		if !exists {
			skill = &models.Skill{
				ID:         skillID,
				Name:       skillName,
				Categories: []models.Category{},
			}
			skillMap[skillID] = skill
		}

		// Add category if it exists
		if categoryID.Valid && categoryName.Valid {
			category := models.Category{
				ID:   int(categoryID.Int64),
				Name: categoryName.String,
			}
			skill.Categories = append(skill.Categories, category)
		}
	}

	// Convert map to slice
	var skills []models.Skill
	for _, skill := range skillMap {
		skills = append(skills, *skill)
	}

	return skills, nil
}

func (r *skillRepository) GetByID(id int) (*models.Skill, error) {
	query := r.MustGetQuery("get_skill_by_id")
	var skill models.Skill
	err := r.db.QueryRow(query, id).Scan(&skill.ID, &skill.Name)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) GetByName(name string) (*models.Skill, error) {
	query := r.MustGetQuery("get_skill_by_name")
	var skill models.Skill
	err := r.db.QueryRow(query, name).Scan(&skill.ID, &skill.Name)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) Create(skill *models.Skill) (*models.Skill, error) {
	query := r.MustGetQuery("create_skill")
	var createdSkill models.Skill
	err := r.db.QueryRow(query, skill.Name).
		Scan(&createdSkill.ID, &createdSkill.Name)
	if err != nil {
		return nil, err
	}

	return &createdSkill, nil
}

func (r *skillRepository) Update(id int, skill *models.Skill) (*models.Skill, error) {
	query := r.MustGetQuery("update_skill")
	var updatedSkill models.Skill
	err := r.db.QueryRow(query, skill.Name, id).
		Scan(&updatedSkill.ID, &updatedSkill.Name)
	if err != nil {
		return nil, err
	}

	return &updatedSkill, nil
}

func (r *skillRepository) SearchSkills(query string) ([]models.Skill, error) {
	sqlQuery := r.MustGetQuery("search_skills")
	rows, err := r.db.Query(sqlQuery, query)
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

func (r *skillRepository) GetPopularSkills(limit int) ([]models.Skill, error) {
	query := r.MustGetQuery("get_popular_skills")
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var skill models.Skill
		var employeeCount int
		err := rows.Scan(&skill.ID, &skill.Name, &employeeCount)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (r *skillRepository) GetSkillsByIDs(ids []int) ([]models.Skill, error) {
	if len(ids) == 0 {
		return []models.Skill{}, nil
	}

	query := r.MustGetQuery("get_skills_by_ids")
	rows, err := r.db.Query(query, fmt.Sprintf("{%s}", strings.Join(strings.Split(fmt.Sprintf("%d", ids[0]), ""), ",")))
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

func (r *skillRepository) GetSkillsWithEmployeeCount() ([]models.SkillWithCount, error) {
	query := r.MustGetQuery("get_skills_with_employee_count")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.SkillWithCount
	for rows.Next() {
		var skill models.SkillWithCount
		err := rows.Scan(&skill.ID, &skill.Name, &skill.EmployeeCount)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (r *skillRepository) CreateBatch(skills []models.Skill) ([]models.Skill, error) {
	if len(skills) == 0 {
		return []models.Skill{}, nil
	}

	// Build dynamic query for batch insert
	valueStrings := make([]string, len(skills))
	valueArgs := make([]interface{}, len(skills))

	for i, skill := range skills {
		valueStrings[i] = fmt.Sprintf("($%d)", i+1)
		valueArgs[i] = skill.Name
	}

	query := fmt.Sprintf("INSERT INTO skills (name) VALUES %s RETURNING id, name",
		strings.Join(valueStrings, ","))

	rows, err := r.db.Query(query, valueArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var createdSkills []models.Skill
	for rows.Next() {
		var skill models.Skill
		err := rows.Scan(&skill.ID, &skill.Name)
		if err != nil {
			return nil, err
		}
		createdSkills = append(createdSkills, skill)
	}

	return createdSkills, nil
}

func (r *skillRepository) UpdateBatch(updates []models.SkillUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, update := range updates {
		var query string
		var args []interface{}

		if update.Name != nil {
			query = "UPDATE skills SET name = $1 WHERE id = $2"
			args = []interface{}{*update.Name, update.ID}
		} else {
			continue // Skip if no fields to update
		}

		_, err := tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *skillRepository) DeleteBatch(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := r.MustGetQuery("delete_skills_batch")
	_, err := r.db.Exec(query, fmt.Sprintf("{%s}", strings.Join(strings.Split(fmt.Sprintf("%d", ids[0]), ""), ",")))
	return err
}

func (r *skillRepository) GetSkillStats() (*models.SkillStats, error) {
	query := r.MustGetQuery("get_skill_stats")
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

	// Get category stats using normalized schema
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

func (r *skillRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_skill")
	_, err := r.db.Exec(query, id)
	return err
}

func (r *skillRepository) GetByCategoryID(categoryID int) ([]models.Skill, error) {
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

func (r *skillRepository) GetSkillsWithCategories() ([]models.Skill, error) {
	query := r.MustGetQuery("get_skills_with_categories")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	skillMap := make(map[int]*models.Skill)
	for rows.Next() {
		var skillID int
		var skillName string
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		err := rows.Scan(&skillID, &skillName, &categoryID, &categoryName)
		if err != nil {
			return nil, err
		}

		// Get or create skill
		skill, exists := skillMap[skillID]
		if !exists {
			skill = &models.Skill{
				ID:         skillID,
				Name:       skillName,
				Categories: []models.Category{},
			}
			skillMap[skillID] = skill
		}

		// Add category if it exists
		if categoryID.Valid && categoryName.Valid {
			category := models.Category{
				ID:   int(categoryID.Int64),
				Name: categoryName.String,
			}
			skill.Categories = append(skill.Categories, category)
		}
	}

	// Convert map to slice
	var skills []models.Skill
	for _, skill := range skillMap {
		skills = append(skills, *skill)
	}

	return skills, nil
}

func (r *skillRepository) AddSkillToCategory(skillID, categoryID int) error {
	query := r.MustGetQuery("add_skill_to_category")
	_, err := r.db.Exec(query, skillID, categoryID)
	return err
}

func (r *skillRepository) RemoveSkillFromCategory(skillID, categoryID int) error {
	query := r.MustGetQuery("remove_skill_from_category")
	_, err := r.db.Exec(query, skillID, categoryID)
	return err
}

func (r *skillRepository) GetSkillCategories(skillID int) ([]models.Category, error) {
	query := r.MustGetQuery("get_skill_categories")
	rows, err := r.db.Query(query, skillID)
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

func (r *skillRepository) AssociateCategories(skillID int, categoryIDs []int) error {
	if len(categoryIDs) == 0 {
		return nil
	}

	// First, remove existing associations
	query := r.MustGetQuery("remove_all_skill_categories")
	_, err := r.db.Exec(query, skillID)
	if err != nil {
		return err
	}

	// Then add new associations
	query = r.MustGetQuery("add_skill_to_category")
	for _, categoryID := range categoryIDs {
		_, err := r.db.Exec(query, skillID, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

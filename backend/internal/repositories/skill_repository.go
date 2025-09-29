package repositories

import (
	"database/sql"
	"stafind-backend/internal/models"
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

	var skills []models.Skill
	for rows.Next() {
		var skill models.Skill
		err := rows.Scan(&skill.ID, &skill.Name, &skill.Category)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func (r *skillRepository) GetByID(id int) (*models.Skill, error) {
	query := r.MustGetQuery("get_skill_by_id")
	var skill models.Skill
	err := r.db.QueryRow(query, id).Scan(&skill.ID, &skill.Name, &skill.Category)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) GetByName(name string) (*models.Skill, error) {
	query := r.MustGetQuery("get_skill_by_name")
	var skill models.Skill
	err := r.db.QueryRow(query, name).Scan(&skill.ID, &skill.Name, &skill.Category)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) Create(skill *models.Skill) (*models.Skill, error) {
	query := r.MustGetQuery("create_skill")
	var createdSkill models.Skill
	err := r.db.QueryRow(query, skill.Name, skill.Category).
		Scan(&createdSkill.ID, &createdSkill.Name, &createdSkill.Category)
	if err != nil {
		return nil, err
	}

	return &createdSkill, nil
}

func (r *skillRepository) Update(id int, skill *models.Skill) (*models.Skill, error) {
	query := r.MustGetQuery("update_skill")
	var updatedSkill models.Skill
	err := r.db.QueryRow(query, skill.Name, skill.Category, id).
		Scan(&updatedSkill.ID, &updatedSkill.Name, &updatedSkill.Category)
	if err != nil {
		return nil, err
	}

	return &updatedSkill, nil
}

func (r *skillRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_skill")
	_, err := r.db.Exec(query, id)
	return err
}

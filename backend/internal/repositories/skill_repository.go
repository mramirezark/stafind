package repositories

import (
	"database/sql"
	"stafind-backend/internal/models"
	"stafind-backend/internal/queries"
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
	rows, err := r.db.Query(queries.GetAllSkills)
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
	var skill models.Skill
	err := r.db.QueryRow(queries.GetSkillByID, id).Scan(&skill.ID, &skill.Name, &skill.Category)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) GetByName(name string) (*models.Skill, error) {
	var skill models.Skill
	err := r.db.QueryRow(queries.GetSkillByName, name).Scan(&skill.ID, &skill.Name, &skill.Category)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) Create(skill *models.Skill) (*models.Skill, error) {
	var createdSkill models.Skill
	err := r.db.QueryRow(queries.CreateSkill, skill.Name, skill.Category).
		Scan(&createdSkill.ID, &createdSkill.Name, &createdSkill.Category)
	if err != nil {
		return nil, err
	}

	return &createdSkill, nil
}

func (r *skillRepository) Update(id int, skill *models.Skill) (*models.Skill, error) {
	var updatedSkill models.Skill
	err := r.db.QueryRow(queries.UpdateSkill, skill.Name, skill.Category, id).
		Scan(&updatedSkill.ID, &updatedSkill.Name, &updatedSkill.Category)
	if err != nil {
		return nil, err
	}

	return &updatedSkill, nil
}

func (r *skillRepository) Delete(id int) error {
	_, err := r.db.Exec(queries.DeleteSkill, id)
	return err
}

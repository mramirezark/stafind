package repositories

import (
	"database/sql"
	"encoding/json"
	"stafind-backend/internal/models"
)

type matchRepository struct {
	*BaseRepository
}

// NewMatchRepository creates a new match repository
func NewMatchRepository(db *sql.DB) (MatchRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &matchRepository{BaseRepository: baseRepo}, nil
}

func (r *matchRepository) GetByEmployeeID(employeeID int) ([]models.Match, error) {
	// For now, return empty slice since we don't have a query for this
	// This can be implemented later if needed
	return []models.Match{}, nil
}

func (r *matchRepository) GetAll() ([]models.Match, error) {
	query := r.MustGetQuery("get_all_matches")

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var matchingSkillsJSON string
		var employee models.Employee
		var currentProject sql.NullString

		err := rows.Scan(
			&match.ID, &match.EmployeeID, &match.MatchScore, &matchingSkillsJSON, &match.Notes, &match.CreatedAt,
			&employee.ID, &employee.Name, &employee.Email, &employee.Department, &employee.Level,
			&employee.Location, &employee.Bio, &currentProject, &employee.CreatedAt, &employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse matching skills JSON
		if matchingSkillsJSON != "" {
			json.Unmarshal([]byte(matchingSkillsJSON), &match.MatchingSkills)
		}

		// Handle current project
		if currentProject.Valid {
			employee.CurrentProject = &currentProject.String
		}

		match.Employee = employee
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *matchRepository) Create(match *models.Match) (*models.Match, error) {
	matchingSkillsJSON, _ := json.Marshal(match.MatchingSkills)
	query := r.MustGetQuery("create_match")

	var createdMatch models.Match
	err := r.db.QueryRow(query, match.EmployeeID, match.MatchScore,
		matchingSkillsJSON, match.Notes).
		Scan(&createdMatch.ID, &createdMatch.CreatedAt)
	if err != nil {
		return nil, err
	}

	createdMatch.EmployeeID = match.EmployeeID
	createdMatch.MatchScore = match.MatchScore
	createdMatch.MatchingSkills = match.MatchingSkills
	createdMatch.Notes = match.Notes
	createdMatch.Employee = match.Employee

	return &createdMatch, nil
}

func (r *matchRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_match")
	_, err := r.db.Exec(query, id)
	return err
}

func (r *matchRepository) getEmployeeSkills(employeeID int) ([]models.Skill, error) {
	query := r.MustGetQuery("get_employee_skills_for_match")
	rows, err := r.db.Query(query, employeeID)
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

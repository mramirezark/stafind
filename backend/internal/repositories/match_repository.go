package repositories

import (
	"database/sql"
	"encoding/json"
	"stafind-backend/internal/models"
	"stafind-backend/internal/queries"
)

type matchRepository struct {
	db *sql.DB
}

// NewMatchRepository creates a new match repository
func NewMatchRepository(db *sql.DB) MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) GetByJobRequestID(jobRequestID int) ([]models.Match, error) {
	rows, err := r.db.Query(queries.GetMatchesByJobRequestID, jobRequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var match models.Match
		var matchingSkillsJSON string
		var employee models.Employee

		err := rows.Scan(
			&match.ID, &match.JobRequestID, &match.EmployeeID, &match.MatchScore,
			&matchingSkillsJSON, &match.Notes, &match.CreatedAt,
			&employee.ID, &employee.Name, &employee.Email, &employee.Department,
			&employee.Level, &employee.Location, &employee.Bio,
			&employee.CreatedAt, &employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse matching skills JSON
		json.Unmarshal([]byte(matchingSkillsJSON), &match.MatchingSkills)

		// Get employee skills
		skills, err := r.getEmployeeSkills(employee.ID)
		if err != nil {
			return nil, err
		}
		employee.Skills = skills

		match.Employee = employee
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *matchRepository) Create(match *models.Match) (*models.Match, error) {
	matchingSkillsJSON, _ := json.Marshal(match.MatchingSkills)

	var createdMatch models.Match
	err := r.db.QueryRow(queries.CreateMatch, match.JobRequestID, match.EmployeeID, match.MatchScore,
		matchingSkillsJSON, match.Notes).
		Scan(&createdMatch.ID, &createdMatch.CreatedAt)
	if err != nil {
		return nil, err
	}

	createdMatch.JobRequestID = match.JobRequestID
	createdMatch.EmployeeID = match.EmployeeID
	createdMatch.MatchScore = match.MatchScore
	createdMatch.MatchingSkills = match.MatchingSkills
	createdMatch.Notes = match.Notes
	createdMatch.Employee = match.Employee

	return &createdMatch, nil
}

func (r *matchRepository) DeleteByJobRequestID(jobRequestID int) error {
	_, err := r.db.Exec(queries.DeleteMatchesByJobRequestID, jobRequestID)
	return err
}

func (r *matchRepository) getEmployeeSkills(employeeID int) ([]models.Skill, error) {
	rows, err := r.db.Query(queries.GetEmployeeSkillsForMatch, employeeID)
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

package repositories

import (
	"database/sql"
	"encoding/json"
	"stafind-backend/internal/models"
	"stafind-backend/internal/queries"
)

type jobRequestRepository struct {
	*BaseRepository
}

// NewJobRequestRepository creates a new job request repository
func NewJobRequestRepository(db *sql.DB) (JobRequestRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &jobRequestRepository{BaseRepository: baseRepo}, nil
}

func (r *jobRequestRepository) GetAll() ([]models.JobRequest, error) {
	rows, err := r.db.Query(queries.GetAllJobRequests)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.JobRequest
	for rows.Next() {
		var req models.JobRequest
		var requiredSkillsJSON, preferredSkillsJSON string

		err := rows.Scan(
			&req.ID, &req.Title, &req.Description, &req.Department,
			&requiredSkillsJSON, &preferredSkillsJSON, &req.ExperienceLevel,
			&req.Location, &req.Priority, &req.Status, &req.CreatedBy,
			&req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON arrays
		json.Unmarshal([]byte(requiredSkillsJSON), &req.RequiredSkills)
		json.Unmarshal([]byte(preferredSkillsJSON), &req.PreferredSkills)

		requests = append(requests, req)
	}

	return requests, nil
}

func (r *jobRequestRepository) GetByID(id int) (*models.JobRequest, error) {
	var req models.JobRequest
	var requiredSkillsJSON, preferredSkillsJSON string

	err := r.db.QueryRow(queries.GetJobRequestByID, id).Scan(
		&req.ID, &req.Title, &req.Description, &req.Department,
		&requiredSkillsJSON, &preferredSkillsJSON, &req.ExperienceLevel,
		&req.Location, &req.Priority, &req.Status, &req.CreatedBy,
		&req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse JSON arrays
	json.Unmarshal([]byte(requiredSkillsJSON), &req.RequiredSkills)
	json.Unmarshal([]byte(preferredSkillsJSON), &req.PreferredSkills)

	return &req, nil
}

func (r *jobRequestRepository) Create(req *models.CreateJobRequestRequest) (*models.JobRequest, error) {
	requiredSkillsJSON, _ := json.Marshal(req.RequiredSkills)
	preferredSkillsJSON, _ := json.Marshal(req.PreferredSkills)

	var jobRequest models.JobRequest
	err := r.db.QueryRow(queries.CreateJobRequest, req.Title, req.Description, req.Department,
		requiredSkillsJSON, preferredSkillsJSON, req.ExperienceLevel,
		req.Location, req.Priority, req.CreatedBy).
		Scan(&jobRequest.ID, &jobRequest.CreatedAt, &jobRequest.UpdatedAt)
	if err != nil {
		return nil, err
	}

	jobRequest.Title = req.Title
	jobRequest.Description = req.Description
	jobRequest.Department = req.Department
	jobRequest.RequiredSkills = req.RequiredSkills
	jobRequest.PreferredSkills = req.PreferredSkills
	jobRequest.ExperienceLevel = req.ExperienceLevel
	jobRequest.Location = req.Location
	jobRequest.Priority = req.Priority
	jobRequest.CreatedBy = req.CreatedBy
	jobRequest.Status = "open"

	return &jobRequest, nil
}

func (r *jobRequestRepository) Update(id int, req *models.CreateJobRequestRequest) (*models.JobRequest, error) {
	requiredSkillsJSON, _ := json.Marshal(req.RequiredSkills)
	preferredSkillsJSON, _ := json.Marshal(req.PreferredSkills)

	_, err := r.db.Exec(queries.UpdateJobRequest, req.Title, req.Description, req.Department,
		requiredSkillsJSON, preferredSkillsJSON, req.ExperienceLevel,
		req.Location, req.Priority, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *jobRequestRepository) Delete(id int) error {
	_, err := r.db.Exec(queries.DeleteJobRequest, id)
	return err
}

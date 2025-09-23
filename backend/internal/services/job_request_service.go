package services

import (
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type jobRequestService struct {
	jobRequestRepo repositories.JobRequestRepository
}

// NewJobRequestService creates a new job request service
func NewJobRequestService(jobRequestRepo repositories.JobRequestRepository) JobRequestService {
	return &jobRequestService{
		jobRequestRepo: jobRequestRepo,
	}
}

func (s *jobRequestService) GetAllJobRequests() ([]models.JobRequest, error) {
	return s.jobRequestRepo.GetAll()
}

func (s *jobRequestService) GetJobRequestByID(id int) (*models.JobRequest, error) {
	return s.jobRequestRepo.GetByID(id)
}

func (s *jobRequestService) CreateJobRequest(req *models.CreateJobRequestRequest) (*models.JobRequest, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, &ValidationError{Field: "title", Message: "Title is required"}
	}

	// Validate priority
	validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
	if req.Priority != "" && !validPriorities[req.Priority] {
		return nil, &ValidationError{
			Field:   "priority",
			Message: "Priority must be one of: low, medium, high",
		}
	}

	// Validate experience level
	validLevels := map[string]bool{"junior": true, "mid": true, "senior": true, "staff": true, "principal": true}
	if req.ExperienceLevel != "" && !validLevels[req.ExperienceLevel] {
		return nil, &ValidationError{
			Field:   "experience_level",
			Message: "Experience level must be one of: junior, mid, senior, staff, principal",
		}
	}

	// Set default priority if not provided
	if req.Priority == "" {
		req.Priority = "medium"
	}

	return s.jobRequestRepo.Create(req)
}

func (s *jobRequestService) UpdateJobRequest(id int, req *models.CreateJobRequestRequest) (*models.JobRequest, error) {
	// Check if job request exists
	_, err := s.jobRequestRepo.GetByID(id)
	if err != nil {
		return nil, &NotFoundError{Resource: "job_request", ID: id}
	}

	// Validate required fields (same as create)
	if req.Title == "" {
		return nil, &ValidationError{Field: "title", Message: "Title is required"}
	}

	// Validate priority
	validPriorities := map[string]bool{"low": true, "medium": true, "high": true}
	if req.Priority != "" && !validPriorities[req.Priority] {
		return nil, &ValidationError{
			Field:   "priority",
			Message: "Priority must be one of: low, medium, high",
		}
	}

	// Validate experience level
	validLevels := map[string]bool{"junior": true, "mid": true, "senior": true, "staff": true, "principal": true}
	if req.ExperienceLevel != "" && !validLevels[req.ExperienceLevel] {
		return nil, &ValidationError{
			Field:   "experience_level",
			Message: "Experience level must be one of: junior, mid, senior, staff, principal",
		}
	}

	// Set default priority if not provided
	if req.Priority == "" {
		req.Priority = "medium"
	}

	return s.jobRequestRepo.Update(id, req)
}

func (s *jobRequestService) DeleteJobRequest(id int) error {
	// Check if job request exists
	_, err := s.jobRequestRepo.GetByID(id)
	if err != nil {
		return &NotFoundError{Resource: "job_request", ID: id}
	}

	return s.jobRequestRepo.Delete(id)
}

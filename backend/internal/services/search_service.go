package services

import (
	"stafind-backend/internal/matching"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type searchService struct {
	employeeRepo   repositories.EmployeeRepository
	jobRequestRepo repositories.JobRequestRepository
	matchEngine    *matching.MatchEngine
}

// NewSearchService creates a new search service
func NewSearchService(
	employeeRepo repositories.EmployeeRepository,
	jobRequestRepo repositories.JobRequestRepository,
) SearchService {
	return &searchService{
		employeeRepo:   employeeRepo,
		jobRequestRepo: jobRequestRepo,
		matchEngine:    matching.NewMatchEngine(),
	}
}

func (s *searchService) SearchEmployees(searchReq *models.SearchRequest) ([]models.Match, error) {
	// Get all employees
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Use the matching engine to find matches
	matches := s.matchEngine.SearchEmployees(searchReq, employees)

	return matches, nil
}

func (s *searchService) FindMatchesForJobRequest(jobRequestID int) ([]models.Match, error) {
	// Get the job request
	jobRequest, err := s.jobRequestRepo.GetByID(jobRequestID)
	if err != nil {
		return nil, &NotFoundError{Resource: "job_request", ID: jobRequestID}
	}

	// Get all employees
	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Use the matching engine to find matches
	matches := s.matchEngine.FindMatches(jobRequest, employees)

	return matches, nil
}

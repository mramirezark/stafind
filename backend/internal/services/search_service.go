package services

import (
	"stafind-backend/internal/matching"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
)

type searchService struct {
	employeeRepo repositories.EmployeeRepository
	matchEngine  *matching.MatchEngine
}

// NewSearchService creates a new search service
func NewSearchService(
	employeeRepo repositories.EmployeeRepository,
) SearchService {
	return &searchService{
		employeeRepo: employeeRepo,
		matchEngine:  matching.NewMatchEngine(),
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

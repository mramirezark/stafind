package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"stafind-backend/internal/database"
	"stafind-backend/internal/models"
	"strings"
	"time"

	"github.com/lib/pq"
)

type aiAgentRepository struct {
	*BaseRepository
	db *database.DB
}

// NewAIAgentRepository creates a new AI agent repository
func NewAIAgentRepository(db *database.DB) (AIAgentRepository, error) {
	baseRepo, err := NewBaseRepository(db.DB)
	if err != nil {
		return nil, err
	}

	return &aiAgentRepository{
		BaseRepository: baseRepo,
		db:             db,
	}, nil
}

func (r *aiAgentRepository) Create(req *models.AIAgentRequest) (*models.AIAgentRequest, error) {
	query := r.MustGetQuery("create_ai_agent_request")

	var id int
	var createdAt time.Time
	err := r.db.QueryRow(query,
		req.TeamsMessageID,
		req.ChannelID,
		req.UserID,
		req.UserName,
		req.MessageText,
		req.AttachmentURL,
		req.Status,
		req.CreatedAt,
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, err
	}
	req.ID = id
	req.CreatedAt = createdAt
	return req, nil
}

func (r *aiAgentRepository) GetByID(id int) (*models.AIAgentRequest, error) {
	query := r.MustGetQuery("get_ai_agent_request_by_id")

	req := &models.AIAgentRequest{}
	var extractedSkillsJSON sql.NullString
	var processedAt sql.NullTime
	var attachmentURL, extractedText, error sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&req.ID,
		&req.TeamsMessageID,
		&req.ChannelID,
		&req.UserID,
		&req.UserName,
		&req.MessageText,
		&attachmentURL,
		&extractedText,
		&extractedSkillsJSON,
		&req.Status,
		&error,
		&req.CreatedAt,
		&processedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("AI agent request not found")
		}
		return nil, err
	}

	// Handle extracted skills as PostgreSQL array
	if extractedSkillsJSON.Valid {
		// PostgreSQL arrays are returned as a string like "{skill1,skill2,skill3}"
		// We need to parse this format
		skillsStr := extractedSkillsJSON.String
		if skillsStr != "{}" && skillsStr != "" {
			// Remove the curly braces and split by comma
			skillsStr = skillsStr[1 : len(skillsStr)-1] // Remove { and }
			if skillsStr != "" {
				skills := []string{}
				for _, skill := range strings.Split(skillsStr, ",") {
					skill = strings.TrimSpace(skill)
					if skill != "" {
						skills = append(skills, skill)
					}
				}
				req.ExtractedSkills = skills
			}
		}
	}

	if processedAt.Valid {
		req.ProcessedAt = &processedAt.Time
	}

	return req, nil
}

func (r *aiAgentRepository) GetByTeamsMessageID(teamsMessageID string) (*models.AIAgentRequest, error) {
	query := r.MustGetQuery("get_ai_agent_request_by_teams_message_id")

	req := &models.AIAgentRequest{}
	var extractedSkillsJSON sql.NullString
	var processedAt sql.NullTime

	err := r.db.QueryRow(query, teamsMessageID).Scan(
		&req.ID,
		&req.TeamsMessageID,
		&req.ChannelID,
		&req.UserID,
		&req.UserName,
		&req.MessageText,
		&req.AttachmentURL,
		&req.ExtractedText,
		&extractedSkillsJSON,
		&req.Status,
		&req.Error,
		&req.CreatedAt,
		&processedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("AI agent request not found")
		}
		return nil, err
	}

	// Handle extracted skills as PostgreSQL array
	if extractedSkillsJSON.Valid {
		// PostgreSQL arrays are returned as a string like "{skill1,skill2,skill3}"
		// We need to parse this format
		skillsStr := extractedSkillsJSON.String
		if skillsStr != "{}" && skillsStr != "" {
			// Remove the curly braces and split by comma
			skillsStr = skillsStr[1 : len(skillsStr)-1] // Remove { and }
			if skillsStr != "" {
				skills := []string{}
				for _, skill := range strings.Split(skillsStr, ",") {
					skill = strings.TrimSpace(skill)
					if skill != "" {
						skills = append(skills, skill)
					}
				}
				req.ExtractedSkills = skills
			}
		}
	}

	if processedAt.Valid {
		req.ProcessedAt = &processedAt.Time
	}

	return req, nil
}

func (r *aiAgentRepository) Update(id int, req *models.AIAgentRequest) error {
	query := r.MustGetQuery("update_ai_agent_request")

	// Convert extracted skills to PostgreSQL array format using pq.Array
	var extractedSkillsArray interface{}
	if req.ExtractedSkills != nil && len(req.ExtractedSkills) > 0 {
		extractedSkillsArray = pq.Array(req.ExtractedSkills)
	} else {
		// Use NULL for empty arrays to avoid PostgreSQL array literal issues
		extractedSkillsArray = nil
	}

	_, err := r.db.Exec(query,
		req.ExtractedText,
		extractedSkillsArray,
		req.Status,
		req.Error,
		req.ProcessedAt,
		id,
	)

	return err
}

func (r *aiAgentRepository) GetAll(limit int, offset int) ([]models.AIAgentRequest, error) {
	query := r.MustGetQuery("get_all_ai_agent_requests")

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.AIAgentRequest
	for rows.Next() {
		req := models.AIAgentRequest{}
		var extractedSkillsJSON sql.NullString
		var processedAt sql.NullTime

		err := rows.Scan(
			&req.ID,
			&req.TeamsMessageID,
			&req.ChannelID,
			&req.UserID,
			&req.UserName,
			&req.MessageText,
			&req.AttachmentURL,
			&req.ExtractedText,
			&extractedSkillsJSON,
			&req.Status,
			&req.Error,
			&req.CreatedAt,
			&processedAt,
		)

		if err != nil {
			return nil, err
		}

		// Handle extracted skills as PostgreSQL array
		if extractedSkillsJSON.Valid {
			// PostgreSQL arrays are returned as a string like "{skill1,skill2,skill3}"
			// We need to parse this format
			skillsStr := extractedSkillsJSON.String
			if skillsStr != "{}" && skillsStr != "" {
				// Remove the curly braces and split by comma
				skillsStr = skillsStr[1 : len(skillsStr)-1] // Remove { and }
				if skillsStr != "" {
					skills := []string{}
					for _, skill := range strings.Split(skillsStr, ",") {
						skill = strings.TrimSpace(skill)
						if skill != "" {
							skills = append(skills, skill)
						}
					}
					req.ExtractedSkills = skills
				}
			}
		}

		if processedAt.Valid {
			req.ProcessedAt = &processedAt.Time
		}

		requests = append(requests, req)
	}

	return requests, nil
}

func (r *aiAgentRepository) SaveResponse(response *models.AIAgentResponse) error {
	query := r.MustGetQuery("save_ai_agent_response")

	// Convert matches to JSON
	matchesJSON, err := json.Marshal(response.Matches)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query,
		response.RequestID,
		string(matchesJSON),
		response.Summary,
		response.ProcessingTime,
		response.Status,
		response.Error,
		time.Now(),
	)

	return err
}

func (r *aiAgentRepository) GetResponseByRequestID(requestID int) (*models.AIAgentResponse, error) {
	query := r.MustGetQuery("get_ai_agent_response_by_request_id")

	var response models.AIAgentResponse
	var matchesJSON string
	var error sql.NullString

	err := r.db.QueryRow(query, requestID).Scan(
		&response.RequestID,
		&matchesJSON,
		&response.Summary,
		&response.ProcessingTime,
		&response.Status,
		&error,
		&response.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("AI agent response not found for request ID %d", requestID)
		}
		return nil, err
	}

	// Parse matches JSON
	if matchesJSON != "" {
		err = json.Unmarshal([]byte(matchesJSON), &response.Matches)
		if err != nil {
			return nil, fmt.Errorf("failed to parse matches JSON: %v", err)
		}
	}

	if error.Valid {
		response.Error = error.String
	}

	return &response, nil
}

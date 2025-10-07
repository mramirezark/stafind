package models

import "time"

// HuggingFaceSkillExtractionRequest represents a request for skill extraction using Hugging Face
type HuggingFaceSkillExtractionRequest struct {
	Text                string                 `json:"text" validate:"required"`
	ModelName           string                 `json:"model_name,omitempty"`           // Optional: specific model to use
	ConfidenceThreshold float64                `json:"confidence_threshold,omitempty"` // Optional: minimum confidence score
	MaxSkills           int                    `json:"max_skills,omitempty"`           // Optional: maximum number of skills to extract
	Categories          []string               `json:"categories,omitempty"`           // Optional: specific categories to focus on
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// HuggingFaceSkillExtractionResponse represents the response from Hugging Face skill extraction
type HuggingFaceSkillExtractionResponse struct {
	Success         bool                `json:"success"`
	Skills          []HuggingFaceSkill  `json:"skills"`
	Categories      map[string][]string `json:"categories"`
	TotalSkills     int                 `json:"total_skills"`
	ConfidenceScore float64             `json:"confidence_score"`
	ModelUsed       string              `json:"model_used"`
	ProcessingTime  time.Duration       `json:"processing_time"`
	RawResponse     interface{}         `json:"raw_response,omitempty"`
	Error           string              `json:"error,omitempty"`
}

// HuggingFaceSkill represents a skill extracted using Hugging Face models
type HuggingFaceSkill struct {
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	Confidence     float64  `json:"confidence"`
	StartPosition  int      `json:"start_position"`
	EndPosition    int      `json:"end_position"`
	Context        string   `json:"context,omitempty"`
	NormalizedName string   `json:"normalized_name"`
	Synonyms       []string `json:"synonyms,omitempty"`
}

// HuggingFaceNERResult represents the raw NER result from Hugging Face
type HuggingFaceNERResult struct {
	Entity string  `json:"entity"`
	Score  float64 `json:"score"`
	Index  int     `json:"index"`
	Word   string  `json:"word"`
	Start  int     `json:"start"`
	End    int     `json:"end"`
	Label  string  `json:"label"`
}

// HuggingFaceNERResponse represents the actual NER response format from Hugging Face API
type HuggingFaceNERResponse struct {
	EntityGroup string  `json:"entity_group"`
	Score       float64 `json:"score"`
	Word        string  `json:"word"`
	Start       int     `json:"start"`
	End         int     `json:"end"`
}

// HuggingFaceModelConfig represents configuration for Hugging Face models
type HuggingFaceModelConfig struct {
	ModelName           string  `json:"model_name"`
	APIEndpoint         string  `json:"api_endpoint"`
	ConfidenceThreshold float64 `json:"confidence_threshold"`
	MaxTokens           int     `json:"max_tokens"`
	Temperature         float64 `json:"temperature"`
	TopP                float64 `json:"top_p"`
	MaxLength           int     `json:"max_length"`
}

// SkillExtractionConfig represents configuration for skill extraction
type SkillExtractionConfig struct {
	DefaultModel        string                            `json:"default_model"`
	FallbackModel       string                            `json:"fallback_model"`
	Models              map[string]HuggingFaceModelConfig `json:"models"`
	SkillCategories     []string                          `json:"skill_categories"`
	ConfidenceThreshold float64                           `json:"confidence_threshold"`
	MaxSkillsPerText    int                               `json:"max_skills_per_text"`
	EnableCaching       bool                              `json:"enable_caching"`
	CacheExpiry         time.Duration                     `json:"cache_expiry"`
}

// SkillExtractionStats represents statistics for skill extraction
type SkillExtractionStats struct {
	TotalRequests         int64            `json:"total_requests"`
	SuccessfulRequests    int64            `json:"successful_requests"`
	FailedRequests        int64            `json:"failed_requests"`
	AverageProcessingTime time.Duration    `json:"average_processing_time"`
	MostUsedModel         string           `json:"most_used_model"`
	SkillsExtracted       int64            `json:"skills_extracted"`
	CategoriesFound       map[string]int64 `json:"categories_found"`
	LastUpdated           time.Time        `json:"last_updated"`
}

// HuggingFaceSkillExtractionService interface
type HuggingFaceSkillExtractionService interface {
	ExtractSkills(request *HuggingFaceSkillExtractionRequest) (*HuggingFaceSkillExtractionResponse, error)
	ExtractSkillsFromText(text string) (*HuggingFaceSkillExtractionResponse, error)
	GetAvailableModels() ([]string, error)
	GetModelConfig(modelName string) (*HuggingFaceModelConfig, error)
	GetStats() (*SkillExtractionStats, error)
	HealthCheck() error
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"stafind-backend/internal/models"
)

// HuggingFaceSkillService implements skill extraction using Hugging Face models
type HuggingFaceSkillService struct {
	apiKey          string
	baseURL         string
	config          *models.SkillExtractionConfig
	httpClient      *http.Client
	stats           *models.SkillExtractionStats
	statsMutex      sync.RWMutex
	cache           map[string]*models.HuggingFaceSkillExtractionResponse
	cacheMutex      sync.RWMutex
	skillCategories map[string][]string
	categoryMutex   sync.RWMutex
}

// NewHuggingFaceSkillService creates a new Hugging Face skill extraction service
func NewHuggingFaceSkillService(apiKey string) *HuggingFaceSkillService {
	config := &models.SkillExtractionConfig{
		DefaultModel:        "dbmdz/bert-large-cased-finetuned-conll03-english",
		FallbackModel:       "dslim/bert-base-NER",
		ConfidenceThreshold: 0.5,
		MaxSkillsPerText:    50,
		EnableCaching:       true,
		CacheExpiry:         30 * time.Minute,
		Models: map[string]models.HuggingFaceModelConfig{
			"dbmdz/bert-large-cased-finetuned-conll03-english": {
				ModelName:           "dbmdz/bert-large-cased-finetuned-conll03-english",
				APIEndpoint:         "https://api-inference.huggingface.co/models/dbmdz/bert-large-cased-finetuned-conll03-english",
				ConfidenceThreshold: 0.5,
				MaxTokens:           512,
			},
			"dslim/bert-base-NER": {
				ModelName:           "dslim/bert-base-NER",
				APIEndpoint:         "https://api-inference.huggingface.co/models/dslim/bert-base-NER",
				ConfidenceThreshold: 0.4,
				MaxTokens:           512,
			},
			"microsoft/DialoGPT-medium": {
				ModelName:           "microsoft/DialoGPT-medium",
				APIEndpoint:         "https://api-inference.huggingface.co/models/microsoft/DialoGPT-medium",
				ConfidenceThreshold: 0.3,
				MaxTokens:           256,
			},
		},
		SkillCategories: []string{
			"Programming Languages",
			"Frameworks",
			"Databases",
			"Cloud Platforms",
			"DevOps Tools",
			"Operating Systems",
			"Methodologies",
			"Soft Skills",
		},
	}

	service := &HuggingFaceSkillService{
		apiKey:          apiKey,
		baseURL:         "https://api-inference.huggingface.co",
		config:          config,
		httpClient:      &http.Client{Timeout: 30 * time.Second},
		stats:           &models.SkillExtractionStats{},
		cache:           make(map[string]*models.HuggingFaceSkillExtractionResponse),
		skillCategories: make(map[string][]string),
	}

	// Initialize skill categories
	service.initializeSkillCategories()

	return service
}

// ExtractSkills extracts skills from text using Hugging Face models
func (h *HuggingFaceSkillService) ExtractSkills(request *models.HuggingFaceSkillExtractionRequest) (*models.HuggingFaceSkillExtractionResponse, error) {
	startTime := time.Now()

	// Update stats
	h.updateStats(true, startTime)

	// Check cache if enabled
	if h.config.EnableCaching {
		if cached, found := h.getFromCache(request.Text); found {
			return cached, nil
		}
	}

	// Validate request
	if err := h.validateRequest(request); err != nil {
		h.updateStats(false, startTime)
		return &models.HuggingFaceSkillExtractionResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Determine which model to use
	modelName := request.ModelName
	if modelName == "" {
		modelName = h.config.DefaultModel
	}

	// Extract skills using the specified model
	response, err := h.extractSkillsWithModel(request.Text, modelName, request.ConfidenceThreshold)
	if err != nil {
		// Try fallback model if primary fails
		if modelName != h.config.FallbackModel {
			response, err = h.extractSkillsWithModel(request.Text, h.config.FallbackModel, request.ConfidenceThreshold)
		}
		if err != nil {
			h.updateStats(false, startTime)
			return &models.HuggingFaceSkillExtractionResponse{
				Success: false,
				Error:   err.Error(),
			}, err
		}
	}

	// Process and categorize skills
	processedResponse := h.processSkills(response, request, time.Since(startTime))

	// Cache the result if enabled
	if h.config.EnableCaching {
		h.setCache(request.Text, processedResponse)
	}

	return processedResponse, nil
}

// ExtractSkillsFromText is a convenience method for simple text extraction
func (h *HuggingFaceSkillService) ExtractSkillsFromText(text string) (*models.HuggingFaceSkillExtractionResponse, error) {
	request := &models.HuggingFaceSkillExtractionRequest{
		Text: text,
	}
	return h.ExtractSkills(request)
}

// extractSkillsWithModel calls the Hugging Face API for a specific model
func (h *HuggingFaceSkillService) extractSkillsWithModel(text, modelName string, confidenceThreshold float64) (*models.HuggingFaceSkillExtractionResponse, error) {
	modelConfig, exists := h.config.Models[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found in configuration", modelName)
	}

	// Prepare the request payload
	payload := map[string]interface{}{
		"inputs": text,
		"parameters": map[string]interface{}{
			"aggregation_strategy": "simple", // For NER models, use aggregation_strategy instead
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", modelConfig.APIEndpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.apiKey)

	// Make the request
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Hugging Face API: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		var apiError map[string]interface{}
		if err := json.Unmarshal(body, &apiError); err == nil {
			if errorMsg, ok := apiError["error"].(string); ok {
				// Check for specific permission errors
				if strings.Contains(errorMsg, "sufficient permissions") || strings.Contains(errorMsg, "authentication method") {
					return nil, fmt.Errorf("hugging Face API permission error: %s. Please check your API key has Inference API permissions enabled", errorMsg)
				}
				return nil, fmt.Errorf("hugging Face API error: %s", errorMsg)
			}
		}
		return nil, fmt.Errorf("hugging Face API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the NER results - NER models return a different format
	var nerResponses []models.HuggingFaceNERResponse
	if err := json.Unmarshal(body, &nerResponses); err != nil {
		// Log the raw response for debugging
		fmt.Printf("DEBUG: Failed to parse NER response as array, raw response: %s\n", string(body))
		return nil, fmt.Errorf("failed to parse NER results: %w", err)
	}

	// Convert to our internal format
	var nerResults []models.HuggingFaceNERResult
	for _, resp := range nerResponses {
		nerResults = append(nerResults, models.HuggingFaceNERResult{
			Entity: resp.EntityGroup,
			Score:  resp.Score,
			Word:   resp.Word,
			Start:  resp.Start,
			End:    resp.End,
			Label:  resp.EntityGroup,
		})
	}

	// Debug logging
	fmt.Printf("DEBUG: Found %d NER results\n", len(nerResults))
	for i, result := range nerResults {
		fmt.Printf("DEBUG: Result %d: Word='%s', Label='%s', Score=%.3f\n", i+1, result.Word, result.Label, result.Score)
	}

	// Process the results
	response := &models.HuggingFaceSkillExtractionResponse{
		Success:     true,
		ModelUsed:   modelName,
		RawResponse: nerResults,
	}

	// Convert NER results to skills
	skills := h.convertNERResultsToSkills(nerResults, text, confidenceThreshold)
	fmt.Printf("DEBUG: Extracted %d skills after filtering\n", len(skills))
	for i, skill := range skills {
		fmt.Printf("DEBUG: Skill %d: %s (%s) - Confidence: %.3f\n", i+1, skill.Name, skill.Category, skill.Confidence)
	}

	response.Skills = skills
	response.TotalSkills = len(skills)
	response.Categories = h.categorizeSkills(skills)

	return response, nil
}

// convertNERResultsToSkills converts Hugging Face NER results to skill objects
func (h *HuggingFaceSkillService) convertNERResultsToSkills(nerResults []models.HuggingFaceNERResult, text string, confidenceThreshold float64) []models.HuggingFaceSkill {
	var skills []models.HuggingFaceSkill
	seenSkills := make(map[string]bool)

	for _, result := range nerResults {
		// Filter by confidence threshold
		if result.Score < confidenceThreshold {
			continue
		}

		// Filter for skill-related entities (MISC, ORG, PER, etc.)
		if !h.isSkillRelatedEntity(result.Label) {
			continue
		}

		// Normalize skill name
		normalizedName := h.normalizeSkillName(result.Word)
		if normalizedName == "" || seenSkills[normalizedName] {
			continue
		}

		// Extract context around the skill
		context := h.extractContext(text, result.Start, result.End)

		skill := models.HuggingFaceSkill{
			Name:           result.Word,
			Category:       h.determineSkillCategory(normalizedName),
			Confidence:     result.Score,
			StartPosition:  result.Start,
			EndPosition:    result.End,
			Context:        context,
			NormalizedName: normalizedName,
			Synonyms:       h.generateSynonyms(normalizedName),
		}

		skills = append(skills, skill)
		seenSkills[normalizedName] = true
	}

	// Sort by confidence score (highest first)
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Confidence > skills[j].Confidence
	})

	return skills
}

// processSkills processes and finalizes the skill extraction results
func (h *HuggingFaceSkillService) processSkills(response *models.HuggingFaceSkillExtractionResponse, request *models.HuggingFaceSkillExtractionRequest, processingTime time.Duration) *models.HuggingFaceSkillExtractionResponse {
	// Apply max skills limit if specified
	if request.MaxSkills > 0 && len(response.Skills) > request.MaxSkills {
		response.Skills = response.Skills[:request.MaxSkills]
		response.TotalSkills = request.MaxSkills
	}

	// Filter by categories if specified
	if len(request.Categories) > 0 {
		filteredSkills := []models.HuggingFaceSkill{}
		for _, skill := range response.Skills {
			for _, category := range request.Categories {
				if strings.EqualFold(skill.Category, category) {
					filteredSkills = append(filteredSkills, skill)
					break
				}
			}
		}
		response.Skills = filteredSkills
		response.TotalSkills = len(filteredSkills)
	}

	// Recalculate categories after filtering
	response.Categories = h.categorizeSkills(response.Skills)

	// Calculate overall confidence score
	response.ConfidenceScore = h.calculateOverallConfidence(response.Skills)

	// Set processing time
	response.ProcessingTime = processingTime

	return response
}

// categorizeSkills groups skills by their categories
func (h *HuggingFaceSkillService) categorizeSkills(skills []models.HuggingFaceSkill) map[string][]string {
	categories := make(map[string][]string)

	for _, skill := range skills {
		if categories[skill.Category] == nil {
			categories[skill.Category] = []string{}
		}
		categories[skill.Category] = append(categories[skill.Category], skill.Name)
	}

	return categories
}

// Helper methods

func (h *HuggingFaceSkillService) validateRequest(request *models.HuggingFaceSkillExtractionRequest) error {
	if request.Text == "" {
		return fmt.Errorf("text is required")
	}
	if len(request.Text) > 10000 { // Reasonable limit
		return fmt.Errorf("text is too long (max 10000 characters)")
	}
	return nil
}

func (h *HuggingFaceSkillService) isSkillRelatedEntity(label string) bool {
	// Common NER labels that might contain skills
	skillLabels := []string{"MISC", "ORG", "PER", "LOC", "TECH", "SKILL", "B-TECH", "I-TECH", "B-MISC", "I-MISC"}
	label = strings.ToUpper(label)
	for _, skillLabel := range skillLabels {
		if strings.Contains(label, skillLabel) {
			return true
		}
	}
	// Also accept any label that's not empty and not a common non-skill label
	nonSkillLabels := []string{"O", "B-PER", "I-PER", "B-LOC", "I-LOC", "B-ORG", "I-ORG"}
	for _, nonSkill := range nonSkillLabels {
		if strings.Contains(label, nonSkill) {
			return false
		}
	}
	// If it's not a known non-skill label and has some content, consider it
	return label != "" && len(label) > 1
}

func (h *HuggingFaceSkillService) normalizeSkillName(name string) string {
	// Clean and normalize skill names
	normalized := strings.TrimSpace(name)
	normalized = strings.ToLower(normalized)

	// Remove common prefixes/suffixes
	normalized = strings.TrimPrefix(normalized, "the ")
	normalized = strings.TrimSuffix(normalized, ".")
	normalized = strings.TrimSuffix(normalized, ",")

	// Skip very short or very long names
	if len(normalized) < 2 || len(normalized) > 50 {
		return ""
	}

	return normalized
}

func (h *HuggingFaceSkillService) determineSkillCategory(skillName string) string {
	h.categoryMutex.RLock()
	defer h.categoryMutex.RUnlock()

	skillLower := strings.ToLower(skillName)

	// Check against predefined skill categories
	for category, skills := range h.skillCategories {
		for _, skill := range skills {
			if strings.Contains(skillLower, strings.ToLower(skill)) {
				return category
			}
		}
	}

	// Default categorization based on common patterns
	if h.isProgrammingLanguage(skillName) {
		return "Programming Languages"
	} else if h.isFramework(skillName) {
		return "Frameworks"
	} else if h.isDatabase(skillName) {
		return "Databases"
	} else if h.isCloudPlatform(skillName) {
		return "Cloud Platforms"
	} else if h.isDevOpsTool(skillName) {
		return "DevOps Tools"
	} else if h.isOperatingSystem(skillName) {
		return "Operating Systems"
	} else if h.isMethodology(skillName) {
		return "Methodologies"
	}

	return "Other"
}

func (h *HuggingFaceSkillService) extractContext(text string, start, end int) string {
	// Extract 50 characters before and after the skill
	contextStart := start - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := end + 50
	if contextEnd > len(text) {
		contextEnd = len(text)
	}

	context := text[contextStart:contextEnd]
	return strings.TrimSpace(context)
}

func (h *HuggingFaceSkillService) generateSynonyms(skillName string) []string {
	synonyms := []string{skillName}

	// Add common variations
	variations := []string{
		strings.ToUpper(skillName),
		strings.Title(skillName),
		strings.ReplaceAll(skillName, ".", ""),
		strings.ReplaceAll(skillName, "-", ""),
		strings.ReplaceAll(skillName, "_", ""),
	}

	for _, variation := range variations {
		if variation != skillName && !containsString(synonyms, variation) {
			synonyms = append(synonyms, variation)
		}
	}

	return synonyms
}

func (h *HuggingFaceSkillService) calculateOverallConfidence(skills []models.HuggingFaceSkill) float64 {
	if len(skills) == 0 {
		return 0.0
	}

	totalConfidence := 0.0
	for _, skill := range skills {
		totalConfidence += skill.Confidence
	}

	return totalConfidence / float64(len(skills))
}

// Category detection methods
func (h *HuggingFaceSkillService) isProgrammingLanguage(skillName string) bool {
	programmingLanguages := []string{
		"javascript", "python", "java", "go", "rust", "c++", "c#", "php", "ruby", "swift",
		"kotlin", "typescript", "scala", "r", "matlab", "perl", "haskell", "clojure",
		"erlang", "elixir", "dart", "lua", "assembly", "cobol", "fortran", "pascal",
	}
	return containsString(programmingLanguages, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) isFramework(skillName string) bool {
	frameworks := []string{
		"react", "angular", "vue", "node", "express", "django", "flask", "spring", "laravel",
		"rails", "asp.net", "jquery", "bootstrap", "tailwind", "sass", "less", "webpack",
		"gulp", "grunt", "babel", "typescript", "next.js", "nuxt.js", "svelte",
	}
	return containsString(frameworks, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) isDatabase(skillName string) bool {
	databases := []string{
		"mysql", "postgresql", "mongodb", "redis", "sqlite", "oracle", "sql server",
		"mariadb", "cassandra", "elasticsearch", "dynamodb", "couchdb", "neo4j",
		"influxdb", "timescaledb", "clickhouse", "snowflake", "bigquery",
	}
	return containsString(databases, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) isCloudPlatform(skillName string) bool {
	cloudPlatforms := []string{
		"aws", "azure", "gcp", "google cloud", "amazon web services", "microsoft azure",
		"digital ocean", "linode", "vultr", "heroku", "netlify", "vercel", "cloudflare",
	}
	return containsString(cloudPlatforms, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) isDevOpsTool(skillName string) bool {
	devopsTools := []string{
		"docker", "kubernetes", "jenkins", "gitlab", "github", "bitbucket", "ansible",
		"terraform", "vagrant", "prometheus", "grafana", "elk", "splunk", "newrelic",
		"datadog", "sentry", "circleci", "travis", "github actions", "azure devops",
	}
	return containsString(devopsTools, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) isOperatingSystem(skillName string) bool {
	operatingSystems := []string{
		"linux", "ubuntu", "centos", "redhat", "debian", "windows", "macos", "unix",
		"freebsd", "openbsd", "solaris", "aix", "hp-ux", "android", "ios",
	}
	return containsString(operatingSystems, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) isMethodology(skillName string) bool {
	methodologies := []string{
		"agile", "scrum", "kanban", "waterfall", "devops", "ci/cd", "tdd", "bdd",
		"pair programming", "code review", "refactoring", "microservices", "api",
		"rest", "graphql", "soap", "mvc", "mvp", "mvvm", "clean architecture",
	}
	return containsString(methodologies, strings.ToLower(skillName))
}

func (h *HuggingFaceSkillService) initializeSkillCategories() {
	h.categoryMutex.Lock()
	defer h.categoryMutex.Unlock()

	h.skillCategories = map[string][]string{
		"Programming Languages": {
			"javascript", "python", "java", "go", "rust", "c++", "c#", "php", "ruby",
			"swift", "kotlin", "typescript", "scala", "r", "matlab", "perl", "haskell",
		},
		"Frameworks": {
			"react", "angular", "vue", "node", "express", "django", "flask", "spring",
			"laravel", "rails", "asp.net", "jquery", "bootstrap", "tailwind",
		},
		"Databases": {
			"mysql", "postgresql", "mongodb", "redis", "sqlite", "oracle", "sql server",
			"mariadb", "cassandra", "elasticsearch", "dynamodb", "couchdb", "neo4j",
		},
		"Cloud Platforms": {
			"aws", "azure", "gcp", "google cloud", "amazon web services", "microsoft azure",
			"digital ocean", "linode", "vultr", "heroku", "netlify", "vercel",
		},
		"DevOps Tools": {
			"docker", "kubernetes", "jenkins", "gitlab", "github", "bitbucket", "ansible",
			"terraform", "vagrant", "prometheus", "grafana", "elk", "splunk",
		},
		"Operating Systems": {
			"linux", "ubuntu", "centos", "redhat", "debian", "windows", "macos", "unix",
			"freebsd", "openbsd", "solaris", "android", "ios",
		},
		"Methodologies": {
			"agile", "scrum", "kanban", "waterfall", "devops", "ci/cd", "tdd", "bdd",
			"pair programming", "code review", "refactoring", "microservices",
		},
	}
}

// Cache methods
func (h *HuggingFaceSkillService) getFromCache(text string) (*models.HuggingFaceSkillExtractionResponse, bool) {
	h.cacheMutex.RLock()
	defer h.cacheMutex.RUnlock()

	response, exists := h.cache[text]
	return response, exists
}

func (h *HuggingFaceSkillService) setCache(text string, response *models.HuggingFaceSkillExtractionResponse) {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()

	h.cache[text] = response
}

// Stats methods
func (h *HuggingFaceSkillService) updateStats(success bool, startTime time.Time) {
	h.statsMutex.Lock()
	defer h.statsMutex.Unlock()

	h.stats.TotalRequests++
	if success {
		h.stats.SuccessfulRequests++
	} else {
		h.stats.FailedRequests++
	}

	processingTime := time.Since(startTime)
	h.stats.AverageProcessingTime = (h.stats.AverageProcessingTime + processingTime) / 2
	h.stats.LastUpdated = time.Now()
}

// Interface methods
func (h *HuggingFaceSkillService) GetAvailableModels() ([]string, error) {
	var models []string
	for modelName := range h.config.Models {
		models = append(models, modelName)
	}
	return models, nil
}

func (h *HuggingFaceSkillService) GetModelConfig(modelName string) (*models.HuggingFaceModelConfig, error) {
	config, exists := h.config.Models[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelName)
	}
	return &config, nil
}

func (h *HuggingFaceSkillService) GetStats() (*models.SkillExtractionStats, error) {
	h.statsMutex.RLock()
	defer h.statsMutex.RUnlock()

	// Create a copy to avoid race conditions
	stats := *h.stats
	return &stats, nil
}

func (h *HuggingFaceSkillService) HealthCheck() error {
	// Test with a simple request
	_, err := h.ExtractSkillsFromText("test")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}

// Utility function
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

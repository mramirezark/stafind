package services

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"stafind-backend/internal/repositories"

	"github.com/jdkato/prose/v2"
)

// DatabaseSkillExtractor provides skill extraction using database categories and skills
type DatabaseSkillExtractor struct {
	skillRepo       repositories.SkillRepository
	categoryRepo    repositories.CategoryRepository
	skillsCache     map[string]SkillInfo
	categoriesCache map[string]CategoryInfo
	cacheMutex      sync.RWMutex
	lastCacheUpdate time.Time
	cacheExpiry     time.Duration
}

// SkillInfo contains metadata about a skill from the database
type SkillInfo struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
	Synonyms   []string `json:"synonyms"`
}

// CategoryInfo contains metadata about a category from the database
type CategoryInfo struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Skills []string `json:"skills"`
}

// ExtractedSkills represents the structured skills extracted from text using dynamic categories
type ExtractedSkills struct {
	Categories        map[string][]string `json:"categories"`          // Dynamic categories from database
	EducationLevel    []string            `json:"education_level"`     // Still needed for NER
	LanguagesDetected []string            `json:"languages_detected"`  // Still needed for NER
	YearsOfExperience []string            `json:"years_of_experience"` // Still needed for NER
	Summary           string              `json:"summary"`
	ConfidenceScore   float64             `json:"confidence_score"`
}

// SkillExtractionResult represents the complete result of skill extraction
type SkillExtractionResult struct {
	Skills           ExtractedSkills `json:"skills"`
	TotalSkillsFound int             `json:"total_skills_found"`
	ExtractionMethod string          `json:"extraction_method"`
	AIConfidence     string          `json:"ai_confidence"`
	RawText          string          `json:"raw_text"`
	ProcessingTime   string          `json:"processing_time"`
}

// NERService provides Named Entity Recognition capabilities for skill extraction
type NERService struct {
	databaseExtractor *DatabaseSkillExtractor
}

// NewNERService creates a new instance of NERService with database integration
func NewNERService(skillRepo repositories.SkillRepository, categoryRepo repositories.CategoryRepository) *NERService {
	return &NERService{
		databaseExtractor: NewDatabaseSkillExtractor(skillRepo, categoryRepo),
	}
}

// ExtractSkillsFromText extracts skills from text using database categories and skills
func (n *NERService) ExtractSkillsFromText(text string) (*SkillExtractionResult, error) {
	return n.databaseExtractor.ExtractSkillsFromText(text)
}

// NewDatabaseSkillExtractor creates a new database-backed skill extractor
func NewDatabaseSkillExtractor(skillRepo repositories.SkillRepository, categoryRepo repositories.CategoryRepository) *DatabaseSkillExtractor {
	return &DatabaseSkillExtractor{
		skillRepo:       skillRepo,
		categoryRepo:    categoryRepo,
		skillsCache:     make(map[string]SkillInfo),
		categoriesCache: make(map[string]CategoryInfo),
		cacheExpiry:     30 * time.Minute, // Cache for 30 minutes
	}
}

// LoadSkillsFromDB loads all skills from the database with their categories
func (d *DatabaseSkillExtractor) LoadSkillsFromDB() error {
	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	// Check if cache is still valid
	if time.Since(d.lastCacheUpdate) < d.cacheExpiry && len(d.skillsCache) > 0 {
		return nil
	}

	// Load skills with categories
	skills, err := d.skillRepo.GetSkillsWithCategories()
	if err != nil {
		return fmt.Errorf("failed to load skills from database: %w", err)
	}

	// Clear existing cache
	d.skillsCache = make(map[string]SkillInfo)

	// Build skill cache with normalized names and categories
	for _, skill := range skills {
		skillInfo := SkillInfo{
			ID:         skill.ID,
			Name:       skill.Name,
			Categories: make([]string, len(skill.Categories)),
			Synonyms:   generateSynonyms(skill.Name),
		}

		// Add categories
		for i, category := range skill.Categories {
			skillInfo.Categories[i] = category.Name
		}

		// Add to cache with normalized name
		normalizedName := normalizeSkillName(skill.Name)
		d.skillsCache[normalizedName] = skillInfo

		// Also add original name if different
		if normalizedName != strings.ToLower(skill.Name) {
			d.skillsCache[strings.ToLower(skill.Name)] = skillInfo
		}
	}

	// Load categories
	categories, err := d.categoryRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to load categories from database: %w", err)
	}

	d.categoriesCache = make(map[string]CategoryInfo)

	for _, category := range categories {
		categoryInfo := CategoryInfo{
			ID:     category.ID,
			Name:   category.Name,
			Skills: []string{}, // Will be populated when needed
		}
		d.categoriesCache[strings.ToLower(category.Name)] = categoryInfo
	}

	d.lastCacheUpdate = time.Now()
	return nil
}

// ExtractSkillsFromText extracts skills from text using database data
func (d *DatabaseSkillExtractor) ExtractSkillsFromText(text string) (*SkillExtractionResult, error) {
	// Ensure skills are loaded from database
	if err := d.LoadSkillsFromDB(); err != nil {
		return nil, fmt.Errorf("failed to load skills from database: %w", err)
	}

	// Use Prose for initial NER
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil, fmt.Errorf("failed to create prose document: %w", err)
	}

	// Extract entities using Prose
	entities := doc.Entities()
	tokens := doc.Tokens()

	// Process text for skill extraction
	extractedSkills := &ExtractedSkills{
		Categories: make(map[string][]string),
	}

	// Combine entities and tokens for comprehensive analysis
	allTextParts := make([]string, 0, len(entities)+len(tokens))

	// Add entities
	for _, entity := range entities {
		allTextParts = append(allTextParts, entity.Text)
	}

	// Add tokens that might be skills
	for _, token := range tokens {
		if d.isLikelySkill(token.Text) {
			allTextParts = append(allTextParts, token.Text)
		}
	}

	// Extract skills using database lookup
	d.extractSkillsFromTextParts(allTextParts, extractedSkills)

	// Also search for skills in the full text using regex patterns
	d.extractSkillsWithRegex(text, extractedSkills)

	// Remove duplicates and return result
	d.deduplicateSkills(extractedSkills)

	totalSkills := d.countTotalSkills(extractedSkills)
	extractedSkills.Summary = fmt.Sprintf("Extracted %d skills using database categories", totalSkills)
	extractedSkills.ConfidenceScore = d.calculateConfidenceScore(extractedSkills)

	return &SkillExtractionResult{
		Skills:           *extractedSkills,
		TotalSkillsFound: totalSkills,
		ExtractionMethod: "database_ner",
		AIConfidence:     fmt.Sprintf("%.2f", extractedSkills.ConfidenceScore),
		RawText:          text,
		ProcessingTime:   time.Now().Format(time.RFC3339),
	}, nil
}

// extractSkillsFromTextParts extracts skills from individual text parts
func (d *DatabaseSkillExtractor) extractSkillsFromTextParts(textParts []string, skills *ExtractedSkills) {
	d.cacheMutex.RLock()
	defer d.cacheMutex.RUnlock()

	for _, part := range textParts {
		normalized := normalizeSkillName(part)

		if skillInfo, exists := d.skillsCache[normalized]; exists {
			d.categorizeSkill(skillInfo, skills)
		}
	}
}

// extractSkillsWithRegex searches for skills using regex patterns
func (d *DatabaseSkillExtractor) extractSkillsWithRegex(text string, skills *ExtractedSkills) {
	d.cacheMutex.RLock()
	defer d.cacheMutex.RUnlock()

	// Create regex patterns for common skill formats
	patterns := []string{
		`\b[A-Z][a-z]+(?:\.[A-Z][a-z]+)*\b`, // Capitalized words (React.js, Node.js)
		`\b[a-z]+(?:\.[a-z]+)*\b`,           // Lowercase words (mysql, redis)
		`\b[A-Z]{2,}\b`,                     // Acronyms (AWS, API, SQL)
	}

	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern)
		matches := regex.FindAllString(text, -1)

		for _, match := range matches {
			normalized := normalizeSkillName(match)
			if skillInfo, exists := d.skillsCache[normalized]; exists {
				d.categorizeSkill(skillInfo, skills)
			}
		}
	}
}

// categorizeSkill categorizes a skill into the appropriate category using pure dynamic categories
func (d *DatabaseSkillExtractor) categorizeSkill(skillInfo SkillInfo, skills *ExtractedSkills) {
	for _, category := range skillInfo.Categories {
		// Initialize category slice if it doesn't exist
		if skills.Categories[category] == nil {
			skills.Categories[category] = []string{}
		}

		// Add skill to the category if not already present
		if !contains(skills.Categories[category], skillInfo.Name) {
			skills.Categories[category] = append(skills.Categories[category], skillInfo.Name)
		}
	}
}

// Helper functions
func normalizeSkillName(name string) string {
	// Remove common prefixes/suffixes and normalize
	normalized := strings.ToLower(name)
	normalized = strings.TrimSpace(normalized)
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")
	return normalized
}

func generateSynonyms(name string) []string {
	synonyms := []string{name}

	// Add variations
	variations := []string{
		strings.ToLower(name),
		strings.ToUpper(name),
		strings.ReplaceAll(name, ".", ""),
		strings.ReplaceAll(name, "-", ""),
		strings.ReplaceAll(name, "_", ""),
	}

	for _, variation := range variations {
		if variation != name && !contains(synonyms, variation) {
			synonyms = append(synonyms, variation)
		}
	}

	return synonyms
}

func (d *DatabaseSkillExtractor) isLikelySkill(text string) bool {
	// Simple heuristic to identify potential skills
	if len(text) < 2 || len(text) > 50 {
		return false
	}

	// Check if it's mostly alphanumeric
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.]+$`)
	if !alphanumeric.MatchString(text) {
		return false
	}

	// Check if it contains at least one letter
	hasLetter := regexp.MustCompile(`[a-zA-Z]`)
	return hasLetter.MatchString(text)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (d *DatabaseSkillExtractor) deduplicateSkills(skills *ExtractedSkills) {
	// Deduplicate skills in each dynamic category
	for category, skillList := range skills.Categories {
		skills.Categories[category] = removeDuplicates(skillList)
	}
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

func (d *DatabaseSkillExtractor) countTotalSkills(skills *ExtractedSkills) int {
	total := 0
	for _, skillList := range skills.Categories {
		total += len(skillList)
	}
	return total
}

func (d *DatabaseSkillExtractor) calculateConfidenceScore(skills *ExtractedSkills) float64 {
	totalSkills := d.countTotalSkills(skills)
	if totalSkills == 0 {
		return 0.0
	}

	// Simple confidence calculation based on number of skills found
	// More skills found = higher confidence (up to a point)
	confidence := float64(totalSkills) * 0.1
	if confidence > 0.9 {
		confidence = 0.9
	}

	return confidence
}

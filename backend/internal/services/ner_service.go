package services

import (
	"fmt"
	"regexp"
	"stafind-backend/internal/constants"
	"strconv"
	"strings"
	"time"

	"github.com/jdkato/prose/v2"
)

// NERService handles Named Entity Recognition for skill extraction
type NERService struct{}

// NewNERService creates a new NER service
func NewNERService() *NERService {
	return &NERService{}
}

// ExtractedSkills represents the structured output of skill extraction
type ExtractedSkills struct {
	ProgrammingLanguages []string `json:"programming_languages"`
	WebTechnologies      []string `json:"web_technologies"`
	Databases            []string `json:"databases"`
	CloudDevOps          []string `json:"cloud_devops"`
	SoftSkills           []string `json:"soft_skills"`
	ToolsFrameworks      []string `json:"tools_frameworks"`
	YearsOfExperience    *int     `json:"years_of_experience"`
	EducationLevel       []string `json:"education_level"`
	LanguagesDetected    []string `json:"languages_detected"`
	Summary              string   `json:"summary"`
	ConfidenceScore      float64  `json:"confidence_score"`
}

// SkillExtractionResult represents the complete result of skill extraction
type SkillExtractionResult struct {
	Skills           ExtractedSkills `json:"extracted_skills"`
	TotalSkillsFound int             `json:"total_skills_found"`
	ExtractionMethod string          `json:"skill_extraction_method"`
	AIConfidence     string          `json:"ai_confidence"`
	RawText          string          `json:"raw_text"`
	ProcessingTime   string          `json:"processing_time"`
}

// ExtractSkillsFromText uses Go NER libraries to extract skills from text
func (n *NERService) ExtractSkillsFromText(text string) (*SkillExtractionResult, error) {
	start := time.Now()

	// Create a new document for analysis
	doc, err := prose.NewDocument(text)
	if err != nil {
		return nil, fmt.Errorf("failed to create prose document: %v", err)
	}

	// Extract entities using Prose
	entities := doc.Entities()

	// Initialize skill categories
	skills := &ExtractedSkills{
		ProgrammingLanguages: []string{},
		WebTechnologies:      []string{},
		Databases:            []string{},
		CloudDevOps:          []string{},
		SoftSkills:           []string{},
		ToolsFrameworks:      []string{},
		EducationLevel:       []string{},
		LanguagesDetected:    []string{},
		Summary:              "Skills extracted using Go NER library (Prose)",
		ConfidenceScore:      0.9,
	}

	// Process entities from Prose NER
	for _, entity := range entities {
		entityText := strings.ToLower(entity.Text)
		entityLabel := entity.Label

		// Categorize entities based on their type and content
		switch entityLabel {
		case constants.EntityPerson, constants.EntityOrg, constants.EntityGPE:
			// Check if it's a technology company or skill
			if n.isTechnologyRelated(entityText) {
				n.categorizeTechnology(entityText, skills)
			}
		case constants.EntityWorkOfArt, constants.EntityEvent:
			// Check if it's a framework, library, or tool
			n.categorizeTechnology(entityText, skills)
		}
	}

	// Additional regex-based extraction for comprehensive coverage
	n.extractSkillsWithRegex(text, skills)

	// Extract years of experience
	skills.YearsOfExperience = n.extractYearsOfExperience(text)

	// Extract education level
	skills.EducationLevel = n.extractEducationLevel(text)

	// Detect languages
	skills.LanguagesDetected = n.detectLanguages(text)

	// Calculate total skills found
	totalSkills := len(skills.ProgrammingLanguages) + len(skills.WebTechnologies) +
		len(skills.Databases) + len(skills.CloudDevOps) +
		len(skills.SoftSkills) + len(skills.ToolsFrameworks)

	processingTime := time.Since(start).String()

	return &SkillExtractionResult{
		Skills:           *skills,
		TotalSkillsFound: totalSkills,
		ExtractionMethod: "go_ner_library_prose",
		AIConfidence:     "high",
		RawText:          text,
		ProcessingTime:   processingTime,
	}, nil
}

// isTechnologyRelated checks if an entity is related to technology
func (n *NERService) isTechnologyRelated(text string) bool {
	techKeywords := []string{
		"javascript", "python", "java", "react", "angular", "vue", "node",
		"aws", "azure", "docker", "kubernetes", "mysql", "postgresql",
		"mongodb", "redis", "git", "github", "gitlab", "jenkins",
		"terraform", "ansible", "microservices", "api", "rest",
		"graphql", "typescript", "html", "css", "bootstrap", "tailwind",
		"django", "flask", "spring", "laravel", "rails", "express",
		"redis", "elasticsearch", "cassandra", "oracle", "sql server",
		"sqlite", "dynamodb", "neo4j", "chef", "puppet", "ci", "cd",
	}

	text = strings.ToLower(text)
	for _, keyword := range techKeywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// categorizeTechnology categorizes a technology into appropriate skill category
func (n *NERService) categorizeTechnology(text string, skills *ExtractedSkills) {
	text = strings.ToLower(strings.TrimSpace(text))

	// Programming languages
	programmingLanguages := []string{
		"javascript", "js", "typescript", "ts", "python", "py", "java",
		"c#", "csharp", "c++", "cpp", "php", "ruby", "go", "golang",
		"rust", "swift", "kotlin", "scala", "r", "matlab", "perl",
		"bash", "shell", "powershell", "sql", "pl/sql", "t-sql",
	}

	// Web technologies
	webTechnologies := []string{
		"react", "angular", "vue", "vue.js", "node.js", "nodejs", "node",
		"express", "express.js", "django", "flask", "spring", "spring boot",
		"laravel", "rails", "ruby on rails", "asp.net", "html", "html5",
		"css", "css3", "sass", "scss", "less", "bootstrap", "tailwind",
		"tailwindcss", "jquery", "jquery.js", "next.js", "nuxt.js",
		"gatsby", "svelte", "ember", "backbone", "knockout", "knockout.js",
	}

	// Databases
	databases := []string{
		"mysql", "postgresql", "postgres", "mongodb", "mongo", "redis",
		"elasticsearch", "cassandra", "oracle", "sql server", "sqlite",
		"dynamodb", "neo4j", "couchdb", "couchbase", "influxdb", "timescaledb",
		"mariadb", "percona", "clickhouse", "snowflake", "bigquery",
	}

	// Cloud & DevOps
	cloudDevOps := []string{
		"aws", "amazon web services", "azure", "gcp", "google cloud",
		"docker", "kubernetes", "k8s", "jenkins", "gitlab", "github",
		"terraform", "ansible", "chef", "puppet", "ci", "cd", "cicd",
		"continuous integration", "continuous deployment", "microservices",
		"serverless", "lambda", "ec2", "s3", "rds", "vpc", "iam",
		"cloudformation", "cloudwatch", "route53", "elb", "alb", "nlb",
	}

	// Soft skills
	softSkills := []string{
		"leadership", "communication", "teamwork", "team work",
		"problem solving", "problem-solving", "analytical thinking",
		"creativity", "adaptability", "time management", "project management",
		"mentoring", "collaboration", "liderazgo", "comunicación",
		"trabajo en equipo", "resolución de problemas", "pensamiento analítico",
		"creatividad", "adaptabilidad", "gestión del tiempo", "gestión de proyectos",
		"mentoring", "colaboración", "agile", "scrum", "kanban", "lean",
	}

	// Check programming languages
	for _, lang := range programmingLanguages {
		if strings.Contains(text, lang) {
			if !contains(skills.ProgrammingLanguages, lang) {
				skills.ProgrammingLanguages = append(skills.ProgrammingLanguages, lang)
			}
			return
		}
	}

	// Check web technologies
	for _, tech := range webTechnologies {
		if strings.Contains(text, tech) {
			if !contains(skills.WebTechnologies, tech) {
				skills.WebTechnologies = append(skills.WebTechnologies, tech)
			}
			return
		}
	}

	// Check databases
	for _, db := range databases {
		if strings.Contains(text, db) {
			if !contains(skills.Databases, db) {
				skills.Databases = append(skills.Databases, db)
			}
			return
		}
	}

	// Check cloud & devops
	for _, cloud := range cloudDevOps {
		if strings.Contains(text, cloud) {
			if !contains(skills.CloudDevOps, cloud) {
				skills.CloudDevOps = append(skills.CloudDevOps, cloud)
			}
			return
		}
	}

	// Check soft skills
	for _, skill := range softSkills {
		if strings.Contains(text, skill) {
			if !contains(skills.SoftSkills, skill) {
				skills.SoftSkills = append(skills.SoftSkills, skill)
			}
			return
		}
	}

	// If not categorized, add to tools and frameworks
	if !contains(skills.ToolsFrameworks, text) {
		skills.ToolsFrameworks = append(skills.ToolsFrameworks, text)
	}
}

// extractSkillsWithRegex performs additional regex-based extraction
func (n *NERService) extractSkillsWithRegex(text string, skills *ExtractedSkills) {
	lowerText := strings.ToLower(text)

	// Comprehensive skill patterns
	skillPatterns := map[string][]string{
		"programming_languages": {
			`\b(javascript|js|typescript|ts|python|py|java|c#|csharp|c\+\+|cpp|php|ruby|go|golang|rust|swift|kotlin|scala|r|matlab|perl|bash|shell)\b`,
		},
		"web_technologies": {
			`\b(react|angular|vue|node\.?js|nodejs|express|django|flask|spring|laravel|rails|asp\.net|html|css|sass|scss|less|bootstrap|tailwind|jquery)\b`,
		},
		"databases": {
			`\b(mysql|postgresql|postgres|mongodb|mongo|redis|elasticsearch|cassandra|oracle|sql\s+server|sqlite|dynamodb|neo4j)\b`,
		},
		"cloud_devops": {
			`\b(aws|azure|gcp|docker|kubernetes|k8s|jenkins|gitlab|github|terraform|ansible|chef|puppet|ci/cd|microservices)\b`,
		},
	}

	for category, patterns := range skillPatterns {
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			matches := re.FindAllString(lowerText, -1)

			for _, match := range matches {
				match = strings.TrimSpace(match)
				switch category {
				case "programming_languages":
					if !contains(skills.ProgrammingLanguages, match) {
						skills.ProgrammingLanguages = append(skills.ProgrammingLanguages, match)
					}
				case "web_technologies":
					if !contains(skills.WebTechnologies, match) {
						skills.WebTechnologies = append(skills.WebTechnologies, match)
					}
				case "databases":
					if !contains(skills.Databases, match) {
						skills.Databases = append(skills.Databases, match)
					}
				case "cloud_devops":
					if !contains(skills.CloudDevOps, match) {
						skills.CloudDevOps = append(skills.CloudDevOps, match)
					}
				}
			}
		}
	}
}

// extractYearsOfExperience extracts years of experience from text
func (n *NERService) extractYearsOfExperience(text string) *int {
	experiencePatterns := []string{
		`(\d+)\s*(?:years?|años?)\s*(?:of\s*)?(?:experience|experiencia)`,
		`(?:experience|experiencia)\s*(?:of\s*)?(\d+)\s*(?:years?|años?)`,
		`(\d+)\s*(?:years?|años?)\s*(?:in|en|with)`,
		`(\d+)\+?\s*(?:years?|años?)`,
	}

	for _, pattern := range experiencePatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			if years, err := strconv.Atoi(matches[1]); err == nil {
				return &years
			}
		}
	}

	return nil
}

// extractEducationLevel extracts education level from text
func (n *NERService) extractEducationLevel(text string) []string {
	educationLevels := []string{}
	lowerText := strings.ToLower(text)

	educationPatterns := map[string][]string{
		"en": {"bachelor", "master", "phd", "doctorate", "degree", "diploma", "certification", "certificate", "university", "college"},
		"es": {"licenciatura", "maestría", "doctorado", "grado", "diploma", "certificación", "certificado", "universidad", "colegio", "ingeniería", "ingeniero"},
	}

	for _, patterns := range educationPatterns {
		for _, pattern := range patterns {
			if strings.Contains(lowerText, pattern) && !contains(educationLevels, pattern) {
				educationLevels = append(educationLevels, pattern)
			}
		}
	}

	return educationLevels
}

// detectLanguages detects the languages used in the text
func (n *NERService) detectLanguages(text string) []string {
	languages := []string{}
	lowerText := strings.ToLower(text)

	// English indicators
	englishWords := []string{"experience", "skills", "developer", "engineer", "years", "required", "looking", "need", "hiring"}
	hasEnglish := false
	for _, word := range englishWords {
		if strings.Contains(lowerText, word) {
			hasEnglish = true
			break
		}
	}

	// Spanish indicators
	spanishWords := []string{"experiencia", "habilidades", "desarrollador", "ingeniero", "años", "requerido", "buscamos", "necesitamos", "contratando"}
	hasSpanish := false
	for _, word := range spanishWords {
		if strings.Contains(lowerText, word) {
			hasSpanish = true
			break
		}
	}

	if hasEnglish {
		languages = append(languages, "english")
	}
	if hasSpanish {
		languages = append(languages, "spanish")
	}

	return languages
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

package services

import (
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"stafind-backend/internal/models"
	"stafind-backend/internal/repositories"
	"strings"
)

type skillService struct {
	skillRepo    repositories.SkillRepository
	employeeRepo repositories.EmployeeRepository
}

// NewSkillService creates a new skill service
func NewSkillService(skillRepo repositories.SkillRepository, employeeRepo repositories.EmployeeRepository) SkillService {
	return &skillService{
		skillRepo:    skillRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *skillService) GetAllSkills() ([]models.Skill, error) {
	return s.skillRepo.GetAll()
}

func (s *skillService) GetSkillByID(id int) (*models.Skill, error) {
	return s.skillRepo.GetByID(id)
}

func (s *skillService) GetSkillsByCategoryID(categoryID int) ([]models.Skill, error) {
	if categoryID <= 0 {
		return nil, &ValidationError{Field: "category_id", Message: "Valid category ID is required"}
	}
	return s.skillRepo.GetByCategoryID(categoryID)
}

func (s *skillService) GetSkillsWithCategories() ([]models.Skill, error) {
	return s.skillRepo.GetSkillsWithCategories()
}

func (s *skillService) SearchSkills(query string) ([]models.Skill, error) {
	if query == "" {
		return nil, &ValidationError{Field: "query", Message: "Search query is required"}
	}

	// Sanitize query to prevent injection
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return nil, &ValidationError{Field: "query", Message: "Search query must be at least 2 characters"}
	}

	return s.skillRepo.SearchSkills(query)
}

func (s *skillService) GetPopularSkills(limit int) ([]models.Skill, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	return s.skillRepo.GetPopularSkills(limit)
}

func (s *skillService) GetSkillsByIDs(ids []int) ([]models.Skill, error) {
	if len(ids) == 0 {
		return []models.Skill{}, nil
	}

	// Remove duplicates and validate IDs
	uniqueIDs := make(map[int]bool)
	validIDs := []int{}
	for _, id := range ids {
		if id > 0 && !uniqueIDs[id] {
			uniqueIDs[id] = true
			validIDs = append(validIDs, id)
		}
	}

	return s.skillRepo.GetSkillsByIDs(validIDs)
}

func (s *skillService) GetSkillsWithEmployeeCount() ([]models.SkillWithCount, error) {
	return s.skillRepo.GetSkillsWithEmployeeCount()
}

func (s *skillService) CreateSkill(skill *models.Skill) (*models.Skill, error) {
	// Validate skill
	if err := s.ValidateSkill(skill); err != nil {
		return nil, err
	}

	// Check if skill already exists
	existingSkill, err := s.skillRepo.GetByName(skill.Name)
	if err == nil && existingSkill != nil {
		return nil, &ConflictError{Resource: "skill", Message: "Skill already exists"}
	}

	return s.skillRepo.Create(skill)
}

func (s *skillService) CreateSkillsBatch(skills []models.Skill) ([]models.Skill, error) {
	if len(skills) == 0 {
		return []models.Skill{}, nil
	}

	if len(skills) > 1000 {
		return nil, &ValidationError{Field: "skills", Message: "Cannot create more than 1000 skills at once"}
	}

	// Validate all skills
	for i, skill := range skills {
		if err := s.ValidateSkill(&skill); err != nil {
			return nil, fmt.Errorf("skill at index %d: %w", i, err)
		}
	}

	return s.skillRepo.CreateBatch(skills)
}

func (s *skillService) UpdateSkill(id int, skill *models.Skill) (*models.Skill, error) {
	// Check if skill exists
	_, err := s.skillRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Resource: "skill", ID: id}
		}
		return nil, err
	}

	// Validate skill
	if err := s.ValidateSkill(skill); err != nil {
		return nil, err
	}

	// Check if another skill with the same name exists
	existingSkill, err := s.skillRepo.GetByName(skill.Name)
	if err == nil && existingSkill != nil && existingSkill.ID != id {
		return nil, &ConflictError{Resource: "skill", Message: "Another skill with this name already exists"}
	}

	return s.skillRepo.Update(id, skill)
}

func (s *skillService) UpdateSkillsBatch(updates []models.SkillUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	if len(updates) > 1000 {
		return &ValidationError{Field: "updates", Message: "Cannot update more than 1000 skills at once"}
	}

	// Validate all updates
	for i, update := range updates {
		if update.ID <= 0 {
			return fmt.Errorf("update at index %d: invalid ID", i)
		}

		if update.Name != nil && *update.Name == "" {
			return fmt.Errorf("update at index %d: name cannot be empty", i)
		}
	}

	return s.skillRepo.UpdateBatch(updates)
}

func (s *skillService) DeleteSkill(id int) error {
	// Check if skill exists
	_, err := s.skillRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{Resource: "skill", ID: id}
		}
		return err
	}

	return s.skillRepo.Delete(id)
}

func (s *skillService) DeleteSkillsBatch(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	if len(ids) > 1000 {
		return &ValidationError{Field: "ids", Message: "Cannot delete more than 1000 skills at once"}
	}

	// Remove duplicates and validate IDs
	uniqueIDs := make(map[int]bool)
	validIDs := []int{}
	for _, id := range ids {
		if id > 0 && !uniqueIDs[id] {
			uniqueIDs[id] = true
			validIDs = append(validIDs, id)
		}
	}

	return s.skillRepo.DeleteBatch(validIDs)
}

func (s *skillService) GetSkillStats() (*models.SkillStats, error) {
	return s.skillRepo.GetSkillStats()
}

func (s *skillService) ValidateSkill(skill *models.Skill) error {
	if skill == nil {
		return &ValidationError{Field: "skill", Message: "Skill cannot be nil"}
	}

	if skill.Name == "" {
		return &ValidationError{Field: "name", Message: "Skill name is required"}
	}

	if len(skill.Name) > 100 {
		return &ValidationError{Field: "name", Message: "Skill name cannot exceed 100 characters"}
	}

	// Validate name format (alphanumeric, spaces, hyphens, underscores, dots)
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_\.]+$`)
	if !nameRegex.MatchString(skill.Name) {
		return &ValidationError{Field: "name", Message: "Skill name contains invalid characters"}
	}

	return nil
}

func (s *skillService) SuggestSkillCategories(name string) ([]string, error) {
	if name == "" {
		return []string{}, nil
	}

	// Common category mappings based on skill names
	categoryMappings := map[string][]string{
		// Programming Languages
		"javascript": {"Programming Language"},
		"python":     {"Programming Language"},
		"java":       {"Programming Language"},
		"go":         {"Programming Language"},
		"rust":       {"Programming Language"},
		"swift":      {"Programming Language"},
		"kotlin":     {"Programming Language"},
		"php":        {"Programming Language"},
		"ruby":       {"Programming Language"},
		"c#":         {"Programming Language"},
		"c++":        {"Programming Language"},
		"c":          {"Programming Language"},
		"typescript": {"Programming Language"},
		"dart":       {"Programming Language"},
		"scala":      {"Programming Language"},
		"r":          {"Programming Language"},
		"matlab":     {"Programming Language"},
		"lua":        {"Programming Language"},
		"haskell":    {"Programming Language"},
		"clojure":    {"Programming Language"},
		"erlang":     {"Programming Language"},
		"elixir":     {"Programming Language"},
		"f#":         {"Programming Language"},
		"assembly":   {"Programming Language"},
		"shell":      {"Programming Language"},
		"bash":       {"Programming Language"},
		"powershell": {"Programming Language"},

		// Frontend Frameworks
		"react":   {"Frontend Framework"},
		"vue":     {"Frontend Framework"},
		"angular": {"Frontend Framework"},
		"next":    {"Frontend Framework"},
		"nuxt":    {"Frontend Framework"},
		"svelte":  {"Frontend Framework"},
		"gatsby":  {"Frontend Framework"},
		"astro":   {"Frontend Framework"},
		"remix":   {"Frontend Framework"},
		"solidjs": {"Frontend Framework"},
		"alpine":  {"Frontend Framework"},
		"jquery":  {"Frontend Framework"},
		"lit":     {"Frontend Framework"},
		"preact":  {"Frontend Framework"},

		// Backend Frameworks
		"express":  {"Backend Framework"},
		"fastapi":  {"Backend Framework"},
		"django":   {"Backend Framework"},
		"flask":    {"Backend Framework"},
		"spring":   {"Backend Framework"},
		"laravel":  {"Backend Framework"},
		"symfony":  {"Backend Framework"},
		"rails":    {"Backend Framework"},
		"asp.net":  {"Backend Framework"},
		"deno":     {"Backend Framework"},
		"bun":      {"Backend Framework"},
		"nestjs":   {"Backend Framework"},
		"koa":      {"Backend Framework"},
		"hapi":     {"Backend Framework"},
		"sails":    {"Backend Framework"},
		"meteor":   {"Backend Framework"},
		"feathers": {"Backend Framework"},
		"adonis":   {"Backend Framework"},
		"fastify":  {"Backend Framework"},
		"gin":      {"Backend Framework"},
		"echo":     {"Backend Framework"},
		"fiber":    {"Backend Framework"},

		// Databases
		"postgresql":    {"Database"},
		"mongodb":       {"Database"},
		"redis":         {"Database"},
		"mysql":         {"Database"},
		"sqlite":        {"Database"},
		"oracle":        {"Database"},
		"sql server":    {"Database"},
		"mariadb":       {"Database"},
		"neo4j":         {"Database"},
		"couchdb":       {"Database"},
		"rethinkdb":     {"Database"},
		"influxdb":      {"Database"},
		"timescale":     {"Database"},
		"clickhouse":    {"Database"},
		"snowflake":     {"Database"},
		"bigquery":      {"Database"},
		"redshift":      {"Database"},
		"firebase":      {"Database"},
		"supabase":      {"Database"},
		"planetscale":   {"Database"},
		"cockroach":     {"Database"},
		"fauna":         {"Database"},
		"arangodb":      {"Database"},
		"orientdb":      {"Database"},
		"amazon rds":    {"Database"},
		"azure sql":     {"Database"},
		"dynamodb":      {"Database"},
		"cassandra":     {"Database"},
		"elasticsearch": {"Database"},

		// Cloud Platforms
		"aws":          {"Cloud Platform"},
		"azure":        {"Cloud Platform"},
		"google cloud": {"Cloud Platform"},
		"vercel":       {"Cloud Platform"},
		"netlify":      {"Cloud Platform"},
		"digitalocean": {"Cloud Platform"},
		"linode":       {"Cloud Platform"},
		"heroku":       {"Cloud Platform"},
		"railway":      {"Cloud Platform"},
		"render":       {"Cloud Platform"},
		"fly":          {"Cloud Platform"},
		"cloudflare":   {"Cloud Platform"},
		"alibaba":      {"Cloud Platform"},
		"ibm cloud":    {"Cloud Platform"},
		"oracle cloud": {"Cloud Platform"},
		"tencent":      {"Cloud Platform"},

		// DevOps
		"docker":     {"DevOps"},
		"kubernetes": {"DevOps"},
		"terraform":  {"DevOps"},
		"ansible":    {"DevOps"},
		"jenkins":    {"DevOps"},
		"gitlab":     {"DevOps"},
		"github":     {"DevOps"},
		"circleci":   {"DevOps"},
		"travis":     {"DevOps"},
		"bamboo":     {"DevOps"},
		"teamcity":   {"DevOps"},
		"helm":       {"DevOps"},
		"prometheus": {"DevOps"},
		"grafana":    {"DevOps"},
		"elk":        {"DevOps"},
		"splunk":     {"DevOps"},
		"datadog":    {"DevOps"},
		"new relic":  {"DevOps"},
		"pagerduty":  {"DevOps"},
		"vault":      {"DevOps"},
		"consul":     {"DevOps"},
		"istio":      {"DevOps"},
		"linkerd":    {"DevOps"},
		"nginx":      {"DevOps"},
		"apache":     {"DevOps"},
		"haproxy":    {"DevOps"},
		"traefik":    {"DevOps"},
		"lambda":     {"DevOps"},
		"functions":  {"DevOps"},
		"ci/cd":      {"DevOps"},

		// API Technologies
		"graphql":   {"API"},
		"rest":      {"API"},
		"grpc":      {"API"},
		"websocket": {"API"},
		"webrtc":    {"API"},
		"openapi":   {"API"},
		"swagger":   {"API"},
		"postman":   {"API"},
		"insomnia":  {"API"},
		"thunder":   {"API"},
		"apollo":    {"API"},
		"hasura":    {"API"},
		"prisma":    {"API"},

		// Version Control
		"git":       {"Version Control"},
		"mercurial": {"Version Control"},
		"svn":       {"Version Control"},
		"perforce":  {"Version Control"},
		"plastic":   {"Version Control"},

		// Architecture
		"microservices":      {"Architecture"},
		"serverless":         {"Architecture"},
		"event-driven":       {"Architecture"},
		"cqrs":               {"Architecture"},
		"event sourcing":     {"Architecture"},
		"domain-driven":      {"Architecture"},
		"clean architecture": {"Architecture"},
		"hexagonal":          {"Architecture"},
		"layered":            {"Architecture"},
		"monolithic":         {"Architecture"},
		"service mesh":       {"Architecture"},
		"api gateway":        {"Architecture"},
		"load balancing":     {"Architecture"},
		"circuit breaker":    {"Architecture"},
		"bulkhead":           {"Architecture"},
		"saga":               {"Architecture"},

		// Frontend Technologies
		"html":       {"Frontend Technology"},
		"css":        {"Frontend Technology"},
		"sass":       {"Frontend Technology"},
		"less":       {"Frontend Technology"},
		"stylus":     {"Frontend Technology"},
		"postcss":    {"Frontend Technology"},
		"tailwind":   {"Frontend Technology"},
		"material":   {"Frontend Technology"},
		"ant design": {"Frontend Technology"},
		"chakra":     {"Frontend Technology"},
		"bulma":      {"Frontend Technology"},
		"foundation": {"Frontend Technology"},
		"semantic":   {"Frontend Technology"},
		"vuetify":    {"Frontend Technology"},
		"quasar":     {"Frontend Technology"},

		// Build Tools
		"webpack": {"Build Tool"},
		"vite":    {"Build Tool"},
		"parcel":  {"Build Tool"},
		"rollup":  {"Build Tool"},
		"esbuild": {"Build Tool"},
		"swc":     {"Build Tool"},
		"turbo":   {"Build Tool"},
		"nx":      {"Build Tool"},
		"lerna":   {"Build Tool"},
		"rush":    {"Build Tool"},
		"yarn":    {"Build Tool"},
		"pnpm":    {"Build Tool"},
		"npm":     {"Build Tool"},

		// Testing
		"jest":            {"Testing Framework"},
		"cypress":         {"Testing Framework"},
		"playwright":      {"Testing Framework"},
		"selenium":        {"Testing Framework"},
		"puppeteer":       {"Testing Framework"},
		"vitest":          {"Testing Framework"},
		"testing library": {"Testing Framework"},
		"mocha":           {"Testing Framework"},
		"chai":            {"Testing Framework"},
		"sinon":           {"Testing Framework"},
		"enzyme":          {"Testing Framework"},
		"karma":           {"Testing Framework"},
		"jasmine":         {"Testing Framework"},
		"ava":             {"Testing Framework"},
		"tap":             {"Testing Framework"},

		// Mobile Development
		"react native": {"Mobile Development"},
		"flutter":      {"Mobile Development"},
		"ionic":        {"Mobile Development"},
		"xamarin":      {"Mobile Development"},
		"cordova":      {"Mobile Development"},
		"phonegap":     {"Mobile Development"},
		"expo":         {"Mobile Development"},
		"nativescript": {"Mobile Development"},
		"framework7":   {"Mobile Development"},

		// Desktop Development
		"electron":   {"Desktop Development"},
		"tauri":      {"Desktop Development"},
		"proton":     {"Desktop Development"},
		"neutralino": {"Desktop Development"},
		"wails":      {"Desktop Development"},
		"qt":         {"Desktop Development"},
		"gtk":        {"Desktop Development"},
		"wxwidgets":  {"Desktop Development"},

		// Data Science & AI
		"tensorflow": {"Data Science & AI"},
		"pytorch":    {"Data Science & AI"},
		"scikit":     {"Data Science & AI"},
		"pandas":     {"Data Science & AI"},
		"numpy":      {"Data Science & AI"},
		"matplotlib": {"Data Science & AI"},
		"seaborn":    {"Data Science & AI"},
		"plotly":     {"Data Science & AI"},
		"d3":         {"Data Science & AI"},
		"spark":      {"Data Science & AI"},
		"hadoop":     {"Data Science & AI"},
		"kafka":      {"Data Science & AI"},
		"airflow":    {"Data Science & AI"},
		"jupyter":    {"Data Science & AI"},
		"mlflow":     {"Data Science & AI"},
		"weights":    {"Data Science & AI"},
		"comet":      {"Data Science & AI"},
		"neptune":    {"Data Science & AI"},

		// Security
		"owasp":         {"Security"},
		"jwt":           {"Security"},
		"oauth":         {"Security"},
		"saml":          {"Security"},
		"ldap":          {"Security"},
		"ssl":           {"Security"},
		"tls":           {"Security"},
		"https":         {"Security"},
		"cors":          {"Security"},
		"csrf":          {"Security"},
		"xss":           {"Security"},
		"sql injection": {"Security"},
		"penetration":   {"Security"},
		"vulnerability": {"Security"},
		"auditing":      {"Security"},
		"cryptography":  {"Security"},

		// Soft Skills
		"leadership":        {"Soft Skills"},
		"management":        {"Soft Skills"},
		"agile":             {"Soft Skills"},
		"scrum":             {"Soft Skills"},
		"kanban":            {"Soft Skills"},
		"devops":            {"Soft Skills"},
		"code review":       {"Soft Skills"},
		"writing":           {"Soft Skills"},
		"speaking":          {"Soft Skills"},
		"mentoring":         {"Soft Skills"},
		"collaboration":     {"Soft Skills"},
		"problem solving":   {"Soft Skills"},
		"critical thinking": {"Soft Skills"},
		"communication":     {"Soft Skills"},
		"time management":   {"Soft Skills"},
	}

	// Normalize the skill name for lookup
	normalizedName := strings.ToLower(strings.TrimSpace(name))

	// Direct lookup
	if categories, exists := categoryMappings[normalizedName]; exists {
		return categories, nil
	}

	// Partial matching for compound names
	var suggestions []string
	for key, categories := range categoryMappings {
		if strings.Contains(normalizedName, key) || strings.Contains(key, normalizedName) {
			suggestions = append(suggestions, categories...)
		}
	}

	// Remove duplicates and sort
	uniqueSuggestions := make(map[string]bool)
	var result []string
	for _, suggestion := range suggestions {
		if !uniqueSuggestions[suggestion] {
			uniqueSuggestions[suggestion] = true
			result = append(result, suggestion)
		}
	}

	sort.Strings(result)
	return result, nil
}

func (s *skillService) GetSkillsByEmployeeID(employeeID int) ([]models.Skill, error) {
	return s.employeeRepo.GetSkills(employeeID)
}

func (s *skillService) GetSkillsByEmployeeIDs(employeeIDs []int) (map[int][]models.Skill, error) {
	result := make(map[int][]models.Skill)

	for _, employeeID := range employeeIDs {
		skills, err := s.employeeRepo.GetSkills(employeeID)
		if err != nil {
			return nil, err
		}
		result[employeeID] = skills
	}

	return result, nil
}

func (s *skillService) AddSkillToCategory(skillID, categoryID int) error {
	if skillID <= 0 {
		return &ValidationError{Field: "skill_id", Message: "Valid skill ID is required"}
	}
	if categoryID <= 0 {
		return &ValidationError{Field: "category_id", Message: "Valid category ID is required"}
	}

	// Check if skill exists
	_, err := s.skillRepo.GetByID(skillID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{Resource: "skill", ID: skillID}
		}
		return err
	}

	return s.skillRepo.AddSkillToCategory(skillID, categoryID)
}

func (s *skillService) RemoveSkillFromCategory(skillID, categoryID int) error {
	if skillID <= 0 {
		return &ValidationError{Field: "skill_id", Message: "Valid skill ID is required"}
	}
	if categoryID <= 0 {
		return &ValidationError{Field: "category_id", Message: "Valid category ID is required"}
	}

	// Check if skill exists
	_, err := s.skillRepo.GetByID(skillID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &NotFoundError{Resource: "skill", ID: skillID}
		}
		return err
	}

	return s.skillRepo.RemoveSkillFromCategory(skillID, categoryID)
}

func (s *skillService) GetSkillCategories(skillID int) ([]models.Category, error) {
	if skillID <= 0 {
		return nil, &ValidationError{Field: "skill_id", Message: "Valid skill ID is required"}
	}

	// Check if skill exists
	_, err := s.skillRepo.GetByID(skillID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &NotFoundError{Resource: "skill", ID: skillID}
		}
		return nil, err
	}

	return s.skillRepo.GetSkillCategories(skillID)
}

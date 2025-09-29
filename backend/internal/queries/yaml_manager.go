package queries

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed config/*.yaml *.sql
var yamlQueryFiles embed.FS

// QueryConfig represents the structure of the YAML configuration
type QueryConfig struct {
	Queries    map[string]map[string]QueryDefinition `yaml:"queries"`
	Categories map[string]CategoryDefinition         `yaml:"categories"`
	Settings   SettingsDefinition                    `yaml:"settings"`
	Domains    []DomainDefinition                    `yaml:"domains,omitempty"`
}

// DomainDefinition represents a domain configuration file
type DomainDefinition struct {
	File        string `yaml:"file"`
	Description string `yaml:"description"`
}

// QueryDefinition represents a single query configuration
type QueryDefinition struct {
	Description string                `yaml:"description"`
	Category    string                `yaml:"category"`
	Operation   string                `yaml:"operation"`
	Parameters  []ParameterDefinition `yaml:"parameters"`
	Tags        []string              `yaml:"tags"`
	SQLFile     string                `yaml:"sql_file"`
	SQL         string                `yaml:"sql,omitempty"` // For inline SQL
}

// ParameterDefinition represents a query parameter
type ParameterDefinition struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
}

// CategoryDefinition represents a query category
type CategoryDefinition struct {
	Description string `yaml:"description"`
	Color       string `yaml:"color"`
}

// SettingsDefinition represents global settings
type SettingsDefinition struct {
	DefaultSQLDirectory       string `yaml:"default_sql_directory"`
	EnableQueryValidation     bool   `yaml:"enable_query_validation"`
	EnableParameterValidation bool   `yaml:"enable_parameter_validation"`
	CacheQueries              bool   `yaml:"cache_queries"`
}

// YAMLQueryManager manages queries using YAML configuration
type YAMLQueryManager struct {
	config     *QueryConfig
	queries    map[string]string
	categories map[string]CategoryDefinition
	settings   SettingsDefinition
}

// NewYAMLQueryManager creates a new YAML-based query manager
func NewYAMLQueryManager() (*YAMLQueryManager, error) {
	// Load main configuration
	configData, err := yamlQueryFiles.ReadFile("config/main.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config/main.yaml: %w", err)
	}

	var mainConfig QueryConfig
	if err := yaml.Unmarshal(configData, &mainConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config/main.yaml: %w", err)
	}

	// Initialize the query manager with main config
	qm := &YAMLQueryManager{
		config:     &mainConfig,
		queries:    make(map[string]string),
		categories: mainConfig.Categories,
		settings:   mainConfig.Settings,
	}

	// Load domain-specific configuration files
	if err := qm.loadDomainConfigs(); err != nil {
		return nil, fmt.Errorf("failed to load domain configs: %w", err)
	}

	// Load all SQL files and map them to query names
	if err := qm.loadSQLQueries(); err != nil {
		return nil, fmt.Errorf("failed to load SQL queries: %w", err)
	}

	// Validate configuration
	if qm.settings.EnableQueryValidation {
		if err := qm.validateQueries(); err != nil {
			return nil, fmt.Errorf("query validation failed: %w", err)
		}
	}

	return qm, nil
}

// loadSQLQueries loads SQL queries from files and maps them to query names
func (qm *YAMLQueryManager) loadSQLQueries() error {
	// Get all SQL files
	entries, err := yamlQueryFiles.ReadDir(".")
	if err != nil {
		return err
	}

	// Create a map of SQL file contents
	sqlFiles := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			content, err := yamlQueryFiles.ReadFile(entry.Name())
			if err != nil {
				return fmt.Errorf("failed to read SQL file %s: %w", entry.Name(), err)
			}
			sqlFiles[entry.Name()] = string(content)
		}
	}

	// Parse SQL files and extract queries
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			if err := qm.parseSQLFile(entry.Name(), sqlFiles[entry.Name()]); err != nil {
				return fmt.Errorf("failed to parse SQL file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// loadDomainConfigs loads domain-specific configuration files
func (qm *YAMLQueryManager) loadDomainConfigs() error {
	// Get all YAML files in the config directory
	entries, err := yamlQueryFiles.ReadDir("config")
	if err != nil {
		return err
	}

	// Load each domain configuration file
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") && entry.Name() != "main.yaml" {
			configPath := "config/" + entry.Name()
			configData, err := yamlQueryFiles.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", configPath, err)
			}

			var domainConfig QueryConfig
			if err := yaml.Unmarshal(configData, &domainConfig); err != nil {
				return fmt.Errorf("failed to parse %s: %w", configPath, err)
			}

			// Merge queries from domain config into main config
			for category, queries := range domainConfig.Queries {
				if qm.config.Queries == nil {
					qm.config.Queries = make(map[string]map[string]QueryDefinition)
				}
				if qm.config.Queries[category] == nil {
					qm.config.Queries[category] = make(map[string]QueryDefinition)
				}
				for queryName, queryDef := range queries {
					qm.config.Queries[category][queryName] = queryDef
				}
			}
		}
	}

	return nil
}

// parseSQLFile parses a SQL file and extracts individual queries
func (qm *YAMLQueryManager) parseSQLFile(filename, content string) error {
	lines := strings.Split(content, "\n")
	var currentQuery strings.Builder
	var currentQueryName string
	var inQuery bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for query name comment
		if strings.HasPrefix(line, "-- Query name:") {
			// Save previous query if exists
			if inQuery && currentQueryName != "" {
				query := strings.TrimSpace(currentQuery.String())
				qm.queries[currentQueryName] = query
			}

			// Start new query
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				currentQueryName = strings.TrimSpace(parts[1])
				currentQuery.Reset()
				inQuery = true
			}
			continue
		}

		// Skip comments and empty lines
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}

		// Add line to current query
		if inQuery {
			currentQuery.WriteString(line)
			currentQuery.WriteString("\n")
		}
	}

	// Save the last query
	if inQuery && currentQueryName != "" {
		query := strings.TrimSpace(currentQuery.String())
		qm.queries[currentQueryName] = query
	}

	return nil
}

// GetQuery returns a query by name
func (qm *YAMLQueryManager) GetQuery(name string) (string, error) {
	query, exists := qm.queries[name]
	if !exists {
		return "", fmt.Errorf("query '%s' not found", name)
	}
	return query, nil
}

// MustGetQuery returns a query by name, panics if not found
func (qm *YAMLQueryManager) MustGetQuery(name string) string {
	query, err := qm.GetQuery(name)
	if err != nil {
		panic(fmt.Sprintf("Required query '%s' not found: %v", name, err))
	}
	return query
}

// GetQueryDefinition returns the full query definition including metadata
func (qm *YAMLQueryManager) GetQueryDefinition(name string) (*QueryDefinition, error) {
	for category, queries := range qm.config.Queries {
		if queryDef, exists := queries[name]; exists {
			queryDef.Category = category
			return &queryDef, nil
		}
	}
	return nil, fmt.Errorf("query definition '%s' not found", name)
}

// GetQueriesByCategory returns all queries in a specific category
func (qm *YAMLQueryManager) GetQueriesByCategory(category string) ([]string, error) {
	queries, exists := qm.config.Queries[category]
	if !exists {
		return nil, fmt.Errorf("category '%s' not found", category)
	}

	var queryNames []string
	for name := range queries {
		queryNames = append(queryNames, name)
	}
	return queryNames, nil
}

// GetQueriesByTag returns all queries with a specific tag
func (qm *YAMLQueryManager) GetQueriesByTag(tag string) []string {
	var matchingQueries []string

	for category, queries := range qm.config.Queries {
		for name, queryDef := range queries {
			for _, queryTag := range queryDef.Tags {
				if queryTag == tag {
					matchingQueries = append(matchingQueries, fmt.Sprintf("%s.%s", category, name))
				}
			}
		}
	}

	return matchingQueries
}

// GetCategories returns all available categories
func (qm *YAMLQueryManager) GetCategories() map[string]CategoryDefinition {
	return qm.categories
}

// ListQueries returns all available query names
func (qm *YAMLQueryManager) ListQueries() []string {
	var names []string
	for name := range qm.queries {
		names = append(names, name)
	}
	return names
}

// ListQueriesWithMetadata returns all queries with their metadata
func (qm *YAMLQueryManager) ListQueriesWithMetadata() map[string]*QueryDefinition {
	result := make(map[string]*QueryDefinition)

	for category, queries := range qm.config.Queries {
		for name, queryDef := range queries {
			fullName := fmt.Sprintf("%s.%s", category, name)
			queryDef.Category = category
			result[fullName] = &queryDef
		}
	}

	return result
}

// ValidateParameters validates query parameters
func (qm *YAMLQueryManager) ValidateParameters(queryName string, args []interface{}) error {
	if !qm.settings.EnableParameterValidation {
		return nil
	}

	queryDef, err := qm.GetQueryDefinition(queryName)
	if err != nil {
		return err
	}

	requiredCount := 0
	for _, param := range queryDef.Parameters {
		if param.Required {
			requiredCount++
		}
	}

	if len(args) < requiredCount {
		return fmt.Errorf("query '%s' requires at least %d parameters, got %d", queryName, requiredCount, len(args))
	}

	return nil
}

// validateQueries validates the query configuration
func (qm *YAMLQueryManager) validateQueries() error {
	// Check that all queries defined in YAML have corresponding SQL
	for category, queries := range qm.config.Queries {
		for name, queryDef := range queries {
			if queryDef.SQLFile != "" {
				// Check if the query exists in the loaded SQL
				if _, exists := qm.queries[name]; !exists {
					return fmt.Errorf("query '%s' in category '%s' references SQL file '%s' but query not found", name, category, queryDef.SQLFile)
				}
			}
		}
	}

	return nil
}

// GetSettings returns the global settings
func (qm *YAMLQueryManager) GetSettings() SettingsDefinition {
	return qm.settings
}

// ExportQueries exports all queries to a map for external use
func (qm *YAMLQueryManager) ExportQueries() map[string]string {
	result := make(map[string]string)
	for name, query := range qm.queries {
		result[name] = query
	}
	return result
}

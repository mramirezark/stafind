package queries

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed *.sql
var queryFiles embed.FS

// QueryManager manages SQL queries loaded from files
type QueryManager struct {
	queries map[string]string
}

// NewQueryManager creates a new query manager and loads all SQL files
func NewQueryManager() (*QueryManager, error) {
	qm := &QueryManager{
		queries: make(map[string]string),
	}

	// Load embedded SQL files
	entries, err := queryFiles.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read query directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			content, err := queryFiles.ReadFile(entry.Name())
			if err != nil {
				return nil, fmt.Errorf("failed to read query file %s: %w", entry.Name(), err)
			}

			queryName := strings.TrimSuffix(entry.Name(), ".sql")
			qm.queries[queryName] = string(content)
		}
	}

	// Also add the predefined queries as fallback
	qm.loadPredefinedQueries()

	return qm, nil
}

// GetQuery returns a query by name
func (qm *QueryManager) GetQuery(name string) (string, error) {
	query, exists := qm.queries[name]
	if !exists {
		return "", fmt.Errorf("query '%s' not found", name)
	}
	return query, nil
}

// MustGetQuery returns a query by name, panics if not found
func (qm *QueryManager) MustGetQuery(name string) string {
	query, err := qm.GetQuery(name)
	if err != nil {
		panic(fmt.Sprintf("Required query '%s' not found: %v", name, err))
	}
	return query
}

// loadPredefinedQueries loads the predefined queries as fallback
func (qm *QueryManager) loadPredefinedQueries() {
	predefinedQueries := map[string]string{
		"get_all_employees":             GetAllEmployees,
		"get_employee_by_id":            GetEmployeeByID,
		"create_employee":               CreateEmployee,
		"update_employee":               UpdateEmployee,
		"delete_employee":               DeleteEmployee,
		"get_employee_skills":           GetEmployeeSkills,
		"get_skill_by_id":               GetSkillByID,
		"get_skill_by_name":             GetSkillByName,
		"create_skill":                  CreateSkill,
		"add_employee_skill":            AddEmployeeSkill,
		"remove_employee_skills":        RemoveEmployeeSkills,
		"get_all_skills":                GetAllSkills,
		"update_skill":                  UpdateSkill,
		"delete_skill":                  DeleteSkill,
		"create_match":                  CreateMatch,
		"delete_match":                  DeleteMatch,
		"get_employee_skills_for_match": GetEmployeeSkillsForMatch,
	}

	for name, query := range predefinedQueries {
		if _, exists := qm.queries[name]; !exists {
			qm.queries[name] = query
		}
	}
}

// ListQueries returns all available query names
func (qm *QueryManager) ListQueries() []string {
	var names []string
	for name := range qm.queries {
		names = append(names, name)
	}
	return names
}

// ValidateQueries checks if all required queries are available
func (qm *QueryManager) ValidateQueries() error {
	requiredQueries := []string{
		"get_all_employees",
		"get_employee_by_id",
		"create_employee",
		"update_employee",
		"delete_employee",
		"get_employee_skills",
		"get_all_skills",
		"get_skill_by_id",
		"get_skill_by_name",
		"create_skill",
		"update_skill",
		"delete_skill",
	}

	var missing []string
	for _, required := range requiredQueries {
		if _, exists := qm.queries[required]; !exists {
			missing = append(missing, required)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required queries: %s", strings.Join(missing, ", "))
	}

	return nil
}

// BuildDynamicQuery builds a dynamic query based on conditions
func (qm *QueryManager) BuildDynamicQuery(baseQuery string, conditions map[string]interface{}) (string, []interface{}) {
	// This is a simplified version - in a real application, you might want
	// a more sophisticated approach that handles WHERE clauses, JOINs, etc.
	return baseQuery, []interface{}{}
}

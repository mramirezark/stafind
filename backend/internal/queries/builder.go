package queries

import (
	"fmt"
	"strings"
)

// QueryBuilder provides a fluent interface for building SQL queries
type QueryBuilder struct {
	query    strings.Builder
	args     []interface{}
	argCount int
}

// NewQueryBuilder creates a new query builder instance
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		args: make([]interface{}, 0),
	}
}

// Select starts a SELECT query
func (qb *QueryBuilder) Select(columns string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf("SELECT %s", columns))
	return qb
}

// From adds a FROM clause
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" FROM %s", table))
	return qb
}

// Join adds a JOIN clause
func (qb *QueryBuilder) Join(table, condition string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" JOIN %s ON %s", table, condition))
	return qb
}

// LeftJoin adds a LEFT JOIN clause
func (qb *QueryBuilder) LeftJoin(table, condition string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" LEFT JOIN %s ON %s", table, condition))
	return qb
}

// Where adds a WHERE clause
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	if qb.query.Len() > 0 {
		qb.query.WriteString(fmt.Sprintf(" WHERE %s", qb.replacePlaceholders(condition, args...)))
	}
	return qb
}

// And adds an AND condition
func (qb *QueryBuilder) And(condition string, args ...interface{}) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" AND %s", qb.replacePlaceholders(condition, args...)))
	return qb
}

// Or adds an OR condition
func (qb *QueryBuilder) Or(condition string, args ...interface{}) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" OR %s", qb.replacePlaceholders(condition, args...)))
	return qb
}

// OrderBy adds an ORDER BY clause
func (qb *QueryBuilder) OrderBy(columns string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" ORDER BY %s", columns))
	return qb
}

// GroupBy adds a GROUP BY clause
func (qb *QueryBuilder) GroupBy(columns string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" GROUP BY %s", columns))
	return qb
}

// Having adds a HAVING clause
func (qb *QueryBuilder) Having(condition string, args ...interface{}) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" HAVING %s", qb.replacePlaceholders(condition, args...)))
	return qb
}

// Limit adds a LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	return qb
}

// Offset adds an OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" OFFSET %d", offset))
	return qb
}

// Insert starts an INSERT query
func (qb *QueryBuilder) Insert(table string, columns []string, values []interface{}) *QueryBuilder {
	columnsStr := strings.Join(columns, ", ")
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		qb.argCount++
		placeholders[i] = fmt.Sprintf("$%d", qb.argCount)
	}

	qb.query.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, columnsStr, strings.Join(placeholders, ", ")))
	qb.args = append(qb.args, values...)
	return qb
}

// Returning adds a RETURNING clause
func (qb *QueryBuilder) Returning(columns string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf(" RETURNING %s", columns))
	return qb
}

// Update starts an UPDATE query
func (qb *QueryBuilder) Update(table string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf("UPDATE %s", table))
	return qb
}

// Set adds SET clauses for UPDATE
func (qb *QueryBuilder) Set(setClauses map[string]interface{}) *QueryBuilder {
	var sets []string
	for column, value := range setClauses {
		qb.argCount++
		sets = append(sets, fmt.Sprintf("%s = $%d", column, qb.argCount))
		qb.args = append(qb.args, value)
	}
	qb.query.WriteString(fmt.Sprintf(" SET %s", strings.Join(sets, ", ")))
	return qb
}

// Delete starts a DELETE query
func (qb *QueryBuilder) Delete(table string) *QueryBuilder {
	qb.query.WriteString(fmt.Sprintf("DELETE FROM %s", table))
	return qb
}

// Build returns the final query string and arguments
func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.query.String(), qb.args
}

// String returns just the query string (for debugging)
func (qb *QueryBuilder) String() string {
	return qb.query.String()
}

// replacePlaceholders replaces $1, $2, etc. with actual argument numbers
func (qb *QueryBuilder) replacePlaceholders(condition string, args ...interface{}) string {
	result := condition
	for _, arg := range args {
		qb.argCount++
		qb.args = append(qb.args, arg)
		result = strings.Replace(result, "?", fmt.Sprintf("$%d", qb.argCount), 1)
	}
	return result
}

// Predefined query builders for common operations

// BuildEmployeeSelect builds a query to select employees
func BuildEmployeeSelect() *QueryBuilder {
	return NewQueryBuilder().
		Select("e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.created_at, e.updated_at").
		From("employees e")
}

// BuildEmployeeWithSkills builds a query to select employees with their skills
func BuildEmployeeWithSkills() *QueryBuilder {
	return NewQueryBuilder().
		Select("e.id, e.name, e.email, e.department, e.level, e.location, e.bio, e.created_at, e.updated_at, s.id as skill_id, s.name as skill_name, s.category as skill_category, es.proficiency_level, es.years_experience").
		From("employees e").
		LeftJoin("employee_skills es", "e.id = es.employee_id").
		LeftJoin("skills s", "es.skill_id = s.id")
}

// BuildJobRequestSelect builds a query to select job requests
func BuildJobRequestSelect() *QueryBuilder {
	return NewQueryBuilder().
		Select("id, title, description, department, required_skills, preferred_skills, experience_level, location, priority, status, created_by, created_at, updated_at").
		From("job_requests")
}

// BuildSkillSelect builds a query to select skills
func BuildSkillSelect() *QueryBuilder {
	return NewQueryBuilder().
		Select("id, name, category").
		From("skills")
}

// BuildMatchSelect builds a query to select matches with employee details
func BuildMatchSelect() *QueryBuilder {
	return NewQueryBuilder().
		Select("m.id, m.job_request_id, m.employee_id, m.match_score, m.matching_skills, m.notes, m.created_at, e.id as employee_id, e.name as employee_name, e.email as employee_email, e.department as employee_department, e.level as employee_level, e.location as employee_location, e.bio as employee_bio, e.created_at as employee_created_at, e.updated_at as employee_updated_at").
		From("matches m").
		Join("employees e", "m.employee_id = e.id")
}

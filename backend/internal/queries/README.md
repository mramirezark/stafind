# Query Management System

This package provides multiple approaches for managing SQL queries in the application, eliminating hardcoded queries from repository code.

## Approaches Implemented

### 1. YAML Configuration (`config.yaml` + `yaml_manager.go`) ⭐ **RECOMMENDED**
SQL queries are organized using YAML configuration files with rich metadata, categories, and parameter definitions.

**Pros:**
- Rich metadata (descriptions, parameters, categories, tags)
- Excellent organization by domain and operation type
- Parameter validation and documentation
- CLI tools for query exploration
- Query statistics and validation
- Easy to maintain and extend

**Cons:**
- More complex setup
- Requires YAML parsing

### 2. SQL Files with Embed (`*.sql` files + `manager.go`)
SQL queries are stored in separate `.sql` files and embedded into the binary using Go's `embed` package.

**Pros:**
- SQL syntax highlighting in IDEs
- Better organization by domain
- Can be version controlled separately
- Embedded into binary (no external files needed)
- Query validation at startup

**Cons:**
- More complex setup
- Requires Go 1.16+ for embed


## Usage Examples

### Using YAML Query Manager
```go
// In repository
func (r *engineerRepository) GetAll() ([]models.Engineer, error) {
    query := r.MustGetQuery("get_all_engineers")
    rows, err := r.db.Query(query)
    // ... handle results
}

// Get query with metadata
queryDef, err := r.GetQueryManager().GetQueryDefinition("get_all_engineers")
if err == nil {
    fmt.Printf("Query: %s\n", queryDef.Description)
    fmt.Printf("Parameters: %d\n", len(queryDef.Parameters))
}
```

### Using CLI Tools
```bash
# List all queries
go run cmd/query-cli/main.go list

# Show detailed information about a query
go run cmd/query-cli/main.go show get_all_engineers

# List queries by category
go run cmd/query-cli/main.go category

# List queries by tag
go run cmd/query-cli/main.go tag engineers

# Validate all queries
go run cmd/query-cli/main.go validate

# Show statistics
go run cmd/query-cli/main.go stats

# Interactive mode
go run cmd/query-cli/main.go interactive
```

### Using Query Builder
```go
// Dynamic query building
qb := queries.BuildEngineerSelect().
    Where("department = ?", "Engineering").
    OrderBy("name ASC").
    Limit(10)

query, args := qb.Build()
rows, err := db.Query(query, args...)
```

## YAML Configuration Structure

The YAML configuration file (`config.yaml`) provides a comprehensive way to organize queries with rich metadata:

```yaml
queries:
  engineers:
    get_all_engineers:
      description: "Retrieve all engineers with their basic information"
      category: "engineers"
      operation: "select"
      parameters: []
      tags: ["engineers", "list", "basic"]
      sql_file: "engineers.sql"
      
    create_engineer:
      description: "Create a new engineer record"
      category: "engineers"
      operation: "insert"
      parameters:
        - name: "name"
          type: "string"
          required: true
          description: "Engineer full name"
        - name: "email"
          type: "string"
          required: true
          description: "Engineer email address"
      tags: ["engineers", "create", "insert"]
      sql_file: "engineers.sql"

categories:
  engineers:
    description: "Engineer-related queries"
    color: "#3498db"

settings:
  enable_query_validation: true
  enable_parameter_validation: true
  cache_queries: true
```

### Key Features:
- **Rich Metadata**: Descriptions, categories, operations, tags
- **Parameter Documentation**: Type, required status, descriptions
- **Organization**: Grouped by domain (engineers, job_requests, skills, matches)
- **Validation**: Parameter and query validation
- **CLI Tools**: Built-in tools for exploration and management

## Best Practices

1. **Choose the right approach** for your needs:
   - **YAML Configuration**: Best for production applications with rich metadata needs
   - Use SQL files for better organization and syntax highlighting

2. **Naming conventions**:
   - Use descriptive names: `GetAllEngineers`, `CreateJobRequest`
   - Include operation type: `Get`, `Create`, `Update`, `Delete`
   - Include entity name: `Engineers`, `JobRequests`, `Skills`

3. **Query organization**:
   - Group related queries together
   - Use consistent parameter ordering
   - Include comments for complex queries

4. **Error handling**:
   - Validate queries at startup
   - Provide meaningful error messages
   - Handle missing queries gracefully

## Migration from Hardcoded Queries

1. **Extract queries** from repository methods
2. **Define constants** or create SQL files
3. **Update repositories** to use the new system
4. **Test thoroughly** to ensure no regressions
5. **Update documentation** and team guidelines

## Future Enhancements

- Query caching for frequently used queries
- Query performance monitoring
- Automatic query optimization suggestions
- Support for different database dialects
- Query templates with parameter substitution

# Query Configuration Structure

This directory contains the split query configuration files for better organization and maintainability.

## File Structure

- `main.yaml` - Main configuration file with global settings and categories
- `employees.yaml` - Employee and skill management queries
- `api_keys.yaml` - API key management queries
- `ai_agents.yaml` - AI agent processing queries
- `users.yaml` - User and role management queries
- `matches.yaml` - Employee matching queries

## Benefits of Split Configuration

1. **Better Organization**: Each domain has its own configuration file
2. **Easier Maintenance**: Smaller files are easier to navigate and modify
3. **Team Collaboration**: Different teams can work on different domain files
4. **Reduced Conflicts**: Less chance of merge conflicts in version control
5. **Clear Separation**: Each file focuses on a specific business domain

## Adding New Queries

To add new queries:

1. Choose the appropriate domain file (e.g., `employees.yaml` for employee-related queries)
2. Add the query definition under the appropriate category
3. Create or update the corresponding SQL file (e.g., `employees.sql`)
4. Ensure the query name matches the SQL comment in the file

## Query Definition Structure

Each query definition includes:
- `description`: Human-readable description
- `category`: Query category for organization
- `operation`: SQL operation type (select, insert, update, delete)
- `parameters`: List of parameters with types and validation
- `tags`: Searchable tags for query discovery
- `sql_file`: Reference to the SQL file containing the query

## Categories

- `employees`: Employee and skill management
- `api_keys`: API key management
- `ai_agents`: AI agent processing
- `users`: User and role management
- `matches`: Employee matching

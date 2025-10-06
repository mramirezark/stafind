package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"stafind-backend/internal/models"
	"time"
)

type employeeRepository struct {
	*BaseRepository
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db *sql.DB) (EmployeeRepository, error) {
	baseRepo, err := NewBaseRepository(db)
	if err != nil {
		return nil, err
	}

	return &employeeRepository{BaseRepository: baseRepo}, nil
}

func (r *employeeRepository) GetAll() ([]models.Employee, error) {
	query := r.MustGetQuery("get_all_employees_with_skills")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employeeMap := make(map[int]*models.Employee)
	for rows.Next() {
		var employeeID int
		var employeeName, employeeEmail, employeeDepartment, employeeLevel, employeeLocation, employeeBio string
		var currentProject sql.NullString
		var resumeUrl sql.NullString
		var createdAt, updatedAt time.Time
		var skillID sql.NullInt64
		var skillName sql.NullString
		var proficiencyLevel sql.NullInt64
		var yearsExperience sql.NullFloat64

		err := rows.Scan(
			&employeeID, &employeeName, &employeeEmail, &employeeDepartment,
			&employeeLevel, &employeeLocation, &employeeBio, &currentProject, &resumeUrl,
			&createdAt, &updatedAt, &skillID, &skillName, &proficiencyLevel, &yearsExperience,
		)
		if err != nil {
			return nil, err
		}

		// Get or create employee
		employee, exists := employeeMap[employeeID]
		if !exists {
			employee = &models.Employee{
				ID:         employeeID,
				Name:       employeeName,
				Email:      employeeEmail,
				Department: employeeDepartment,
				Level:      employeeLevel,
				Location:   employeeLocation,
				Bio:        employeeBio,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				Skills:     []models.Skill{},
			}
			if currentProject.Valid {
				employee.CurrentProject = &currentProject.String
			}
			if resumeUrl.Valid {
				employee.ResumeUrl = &resumeUrl.String
			}
			employeeMap[employeeID] = employee
		}

		// Add the skill if it exists
		if skillID.Valid && skillName.Valid {
			employee.Skills = append(employee.Skills, models.Skill{
				ID:   int(skillID.Int64),
				Name: skillName.String,
			})
		}
	}

	// Convert map to slice
	var employees []models.Employee
	for _, employee := range employeeMap {
		employees = append(employees, *employee)
	}

	return employees, nil
}

func (r *employeeRepository) GetByID(id int) (*models.Employee, error) {
	query := r.MustGetQuery("get_employee_by_id")
	var employee models.Employee
	err := r.db.QueryRow(query, id).Scan(
		&employee.ID, &employee.Name, &employee.Email, &employee.Department,
		&employee.Level, &employee.Location, &employee.Bio, &employee.CurrentProject, &employee.ResumeUrl,
		&employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get skills for this employee
	skills, err := r.GetSkills(employee.ID)
	if err != nil {
		return nil, err
	}
	employee.Skills = skills

	return &employee, nil
}

func (r *employeeRepository) GetByEmail(email string) (*models.Employee, error) {
	query := r.MustGetQuery("get_employee_by_email")
	var employee models.Employee
	var originalText, extractionSource, extractionStatus sql.NullString
	var extractionTimestamp sql.NullTime
	var extractedDataJSON sql.NullString

	err := r.db.QueryRow(query, email).Scan(
		&employee.ID, &employee.Name, &employee.Email, &employee.Department,
		&employee.Level, &employee.Location, &employee.Bio, &employee.CurrentProject, &employee.ResumeUrl,
		&originalText, &extractedDataJSON, &extractionTimestamp, &extractionSource, &extractionStatus,
		&employee.CreatedAt, &employee.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	// Handle nullable fields
	if originalText.Valid {
		employee.OriginalText = &originalText.String
	}
	if extractionSource.Valid {
		employee.ExtractionSource = &extractionSource.String
	}
	if extractionStatus.Valid {
		employee.ExtractionStatus = &extractionStatus.String
	}
	if extractionTimestamp.Valid {
		employee.ExtractionTimestamp = &extractionTimestamp.Time
	}
	if extractedDataJSON.Valid {
		var extractedData map[string]interface{}
		if err := json.Unmarshal([]byte(extractedDataJSON.String), &extractedData); err == nil {
			employee.ExtractedData = extractedData
		}
	}

	// Get skills for this employee
	skills, err := r.GetSkills(employee.ID)
	if err != nil {
		return nil, err
	}
	employee.Skills = skills

	return &employee, nil
}

func (r *employeeRepository) Create(req *models.CreateEmployeeRequest) (*models.Employee, error) {
	fmt.Printf("DEBUG: Starting employee creation for %s (%s)\n", req.Name, req.Email)
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("DEBUG: Failed to begin transaction: %v\n", err)
		return nil, err
	}
	defer tx.Rollback()

	query := r.MustGetQuery("create_employee")
	var employee models.Employee
	// Convert string to pointer string for current_project
	var currentProject *string
	if req.CurrentProject != "" {
		currentProject = &req.CurrentProject
	}

	// Convert string to pointer string for resume_url
	var resumeUrl *string
	if req.ResumeUrl != "" {
		resumeUrl = &req.ResumeUrl
	}

	fmt.Printf("DEBUG: Executing INSERT query with params: %s, %s, %s, %s, %s, %s, %v\n",
		req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject)

	err = tx.QueryRow(query, req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject, resumeUrl).
		Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		fmt.Printf("DEBUG: Failed to execute INSERT query: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Successfully created employee with ID: %d\n", employee.ID)

	employee.Name = req.Name
	employee.Email = req.Email
	employee.Department = req.Department
	employee.Level = req.Level
	employee.Location = req.Location
	employee.Bio = req.Bio
	employee.CurrentProject = currentProject

	// Add skills
	fmt.Printf("DEBUG: Adding %d skills for employee ID: %d\n", len(req.Skills), employee.ID)
	for i, skillReq := range req.Skills {
		fmt.Printf("DEBUG: Adding skill %d: %s\n", i+1, skillReq.SkillName)
		err = r.addEmployeeSkill(tx, employee.ID, &skillReq)
		if err != nil {
			fmt.Printf("DEBUG: Failed to add skill %s: %v\n", skillReq.SkillName, err)
			return nil, err
		}
	}

	// Get the created employee with skills
	fmt.Printf("DEBUG: Getting created employee with skills\n")
	createdEmployee, err := r.GetByID(employee.ID)
	if err != nil {
		fmt.Printf("DEBUG: Failed to get created employee: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Committing transaction\n")
	err = tx.Commit()
	if err != nil {
		fmt.Printf("DEBUG: Failed to commit transaction: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Successfully created employee: %s (ID: %d)\n", createdEmployee.Name, createdEmployee.ID)
	return createdEmployee, nil
}

func (r *employeeRepository) CreateWithExtraction(req *models.CreateEmployeeRequest, originalText string, extractedData map[string]interface{}, extractionSource, extractionStatus, resumeURL string) (*models.Employee, error) {
	fmt.Printf("DEBUG: Creating employee with extraction data for %s (%s)\n", req.Name, req.Email)
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("DEBUG: Failed to begin transaction: %v\n", err)
		return nil, err
	}
	defer tx.Rollback()

	query := r.MustGetQuery("create_employee_with_extraction")
	var employee models.Employee
	// Convert string to pointer string for current_project
	var currentProject *string
	if req.CurrentProject != "" {
		currentProject = &req.CurrentProject
	}

	// Convert extractedData to JSON
	extractedDataJSON, err := json.Marshal(extractedData)
	if err != nil {
		fmt.Printf("DEBUG: Failed to marshal extracted data: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Executing INSERT with extraction query with params: %s, %s, %s, %s, %s, %s, %v\n",
		req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject)

	err = tx.QueryRow(query, req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject, resumeURL,
		originalText, string(extractedDataJSON), time.Now(), extractionSource, extractionStatus).
		Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		fmt.Printf("DEBUG: Failed to execute INSERT with extraction query: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Successfully created employee with extraction data, ID: %d\n", employee.ID)

	employee.Name = req.Name
	employee.Email = req.Email
	employee.Department = req.Department
	employee.Level = req.Level
	employee.Location = req.Location
	employee.Bio = req.Bio
	employee.CurrentProject = currentProject

	// Add skills
	fmt.Printf("DEBUG: Adding %d skills for employee ID: %d\n", len(req.Skills), employee.ID)
	for i, skillReq := range req.Skills {
		fmt.Printf("DEBUG: Adding skill %d: %s\n", i+1, skillReq.SkillName)
		err = r.addEmployeeSkill(tx, employee.ID, &skillReq)
		if err != nil {
			fmt.Printf("DEBUG: Failed to add skill %s: %v\n", skillReq.SkillName, err)
			return nil, err
		}
	}

	fmt.Printf("DEBUG: Committing transaction\n")
	err = tx.Commit()
	if err != nil {
		fmt.Printf("DEBUG: Failed to commit transaction: %v\n", err)
		return nil, err
	}

	// Get the created employee with skills after transaction is committed
	fmt.Printf("DEBUG: Getting created employee with skills\n")
	createdEmployee, err := r.GetByID(employee.ID)
	if err != nil {
		fmt.Printf("DEBUG: Failed to get created employee: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Successfully created employee with extraction data: %s (ID: %d)\n", createdEmployee.Name, createdEmployee.ID)
	return createdEmployee, nil
}

func (r *employeeRepository) Update(id int, req *models.CreateEmployeeRequest) (*models.Employee, error) {
	fmt.Printf("DEBUG: Updating employee %d\n", id)
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("DEBUG: Failed to begin transaction: %v\n", err)
		return nil, err
	}
	defer tx.Rollback()

	query := r.MustGetQuery("update_employee")
	// Convert string to pointer string for current_project
	var currentProject *string
	if req.CurrentProject != "" {
		currentProject = &req.CurrentProject
	}

	fmt.Printf("DEBUG: Updating employee basic info: %s, %s, %s, %s, %s, %s\n",
		req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio)

	_, err = tx.Exec(query, req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject, id)
	if err != nil {
		fmt.Printf("DEBUG: Failed to update employee: %v\n", err)
		return nil, err
	}

	// Remove existing skills within transaction
	fmt.Printf("DEBUG: Removing existing skills for employee %d\n", id)
	removeSkillsQuery := r.MustGetQuery("remove_employee_skills")
	_, err = tx.Exec(removeSkillsQuery, id)
	if err != nil {
		fmt.Printf("DEBUG: Failed to remove existing skills: %v\n", err)
		return nil, err
	}

	// Add new skills
	fmt.Printf("DEBUG: Adding %d new skills for employee ID: %d\n", len(req.Skills), id)
	for i, skillReq := range req.Skills {
		fmt.Printf("DEBUG: Adding skill %d: %s\n", i+1, skillReq.SkillName)
		err = r.addEmployeeSkill(tx, id, &skillReq)
		if err != nil {
			fmt.Printf("DEBUG: Failed to add skill %s: %v\n", skillReq.SkillName, err)
			return nil, err
		}
	}

	fmt.Printf("DEBUG: Committing transaction\n")
	err = tx.Commit()
	if err != nil {
		fmt.Printf("DEBUG: Failed to commit transaction: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Successfully updated employee %d\n", id)
	return r.GetByID(id)
}

func (r *employeeRepository) UpdateWithExtraction(id int, req *models.CreateEmployeeRequest, originalText string, extractedData map[string]interface{}, extractionSource, extractionStatus, resumeURL string) (*models.Employee, error) {
	fmt.Printf("DEBUG: Updating employee %d with extraction data\n", id)
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("DEBUG: Failed to begin transaction: %v\n", err)
		return nil, err
	}
	defer tx.Rollback()

	query := r.MustGetQuery("update_employee_extraction")
	// Convert string to pointer string for current_project
	var currentProject *string
	if req.CurrentProject != "" {
		currentProject = &req.CurrentProject
	}

	// Convert extractedData to JSON
	extractedDataJSON, err := json.Marshal(extractedData)
	if err != nil {
		fmt.Printf("DEBUG: Failed to marshal extracted data: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Executing UPDATE with extraction query for employee ID: %d\n", id)
	_, err = tx.Exec(query, req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject, resumeURL,
		originalText, string(extractedDataJSON), time.Now(), extractionSource, extractionStatus, id)
	if err != nil {
		fmt.Printf("DEBUG: Failed to execute UPDATE with extraction query: %v\n", err)
		return nil, err
	}

	// Remove existing skills within transaction
	fmt.Printf("DEBUG: Removing existing skills for employee %d\n", id)
	removeSkillsQuery := r.MustGetQuery("remove_employee_skills")
	_, err = tx.Exec(removeSkillsQuery, id)
	if err != nil {
		fmt.Printf("DEBUG: Failed to remove existing skills: %v\n", err)
		return nil, err
	}

	// Add new skills
	fmt.Printf("DEBUG: Adding %d skills for employee ID: %d\n", len(req.Skills), id)
	for i, skillReq := range req.Skills {
		fmt.Printf("DEBUG: Adding skill %d: %s\n", i+1, skillReq.SkillName)
		err = r.addEmployeeSkill(tx, id, &skillReq)
		if err != nil {
			fmt.Printf("DEBUG: Failed to add skill %s: %v\n", skillReq.SkillName, err)
			return nil, err
		}
	}

	fmt.Printf("DEBUG: Committing extraction update transaction\n")
	err = tx.Commit()
	if err != nil {
		fmt.Printf("DEBUG: Failed to commit extraction update transaction: %v\n", err)
		return nil, err
	}

	fmt.Printf("DEBUG: Successfully updated employee with extraction data\n")
	return r.GetByID(id)
}

func (r *employeeRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_employee")
	_, err := r.db.Exec(query, id)
	return err
}

func (r *employeeRepository) GetSkills(employeeID int) ([]models.Skill, error) {
	query := r.MustGetQuery("get_employee_skills")
	fmt.Printf("DEBUG: Getting skills for employee ID: %d\n", employeeID)
	rows, err := r.db.Query(query, employeeID)
	if err != nil {
		fmt.Printf("DEBUG: GetSkills query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var skill models.Skill
		var proficiencyLevel int
		var yearsExperience float64

		err := rows.Scan(&skill.ID, &skill.Name, &proficiencyLevel, &yearsExperience)
		if err != nil {
			fmt.Printf("DEBUG: GetSkills scan error: %v\n", err)
			return nil, err
		}
		skills = append(skills, skill)
	}

	fmt.Printf("DEBUG: Found %d skills for employee ID: %d\n", len(skills), employeeID)
	return skills, nil
}

func (r *employeeRepository) AddSkill(employeeID int, skillReq *models.EmployeeSkillReq) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.addEmployeeSkill(tx, employeeID, skillReq)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *employeeRepository) RemoveSkills(employeeID int) error {
	query := r.MustGetQuery("remove_employee_skills")
	_, err := r.db.Exec(query, employeeID)
	return err
}

func (r *employeeRepository) addEmployeeSkill(tx *sql.Tx, employeeID int, skillReq *models.EmployeeSkillReq) error {
	// First, get or create the skill
	var skillID int
	var skillName string
	query := r.MustGetQuery("get_skill_by_name")
	err := tx.QueryRow(query, skillReq.SkillName).Scan(&skillID, &skillName)
	if err == sql.ErrNoRows {
		// Create the skill
		query = r.MustGetQuery("create_skill")
		err = tx.QueryRow(query, skillReq.SkillName).Scan(&skillID, &skillName)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Add the employee skill relationship
	query = r.MustGetQuery("add_employee_skill")
	_, err = tx.Exec(query, employeeID, skillID, skillReq.ProficiencyLevel, skillReq.YearsExperience)
	return err
}

// GetEmployeesWithSkills gets employees who have any of the specified skills (optimized, no N+1)
func (r *employeeRepository) GetEmployeesWithSkills(skillNames []string) ([]models.Employee, error) {
	query := r.MustGetQuery("get_employees_with_skills_optimized")
	fmt.Printf("DEBUG: Getting employees with skills: %v\n", skillNames)

	rows, err := r.db.Query(query, skillNames)
	if err != nil {
		fmt.Printf("DEBUG: GetEmployeesWithSkills query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	employeeMap := make(map[int]*models.Employee)
	for rows.Next() {
		var employeeID int
		var employeeName, employeeEmail, employeeDepartment, employeeLevel, employeeLocation, employeeBio string
		var currentProject sql.NullString
		var createdAt, updatedAt time.Time
		var skillID int
		var skillName string
		var proficiencyLevel int
		var yearsExperience float64

		err := rows.Scan(
			&employeeID, &employeeName, &employeeEmail, &employeeDepartment,
			&employeeLevel, &employeeLocation, &employeeBio, &currentProject,
			&createdAt, &updatedAt, &skillID, &skillName, &proficiencyLevel, &yearsExperience,
		)
		if err != nil {
			fmt.Printf("DEBUG: GetEmployeesWithSkills scan error: %v\n", err)
			return nil, err
		}

		// Get or create employee
		employee, exists := employeeMap[employeeID]
		if !exists {
			employee = &models.Employee{
				ID:         employeeID,
				Name:       employeeName,
				Email:      employeeEmail,
				Department: employeeDepartment,
				Level:      employeeLevel,
				Location:   employeeLocation,
				Bio:        employeeBio,
				CreatedAt:  createdAt,
				UpdatedAt:  updatedAt,
				Skills:     []models.Skill{},
			}
			if currentProject.Valid {
				employee.CurrentProject = &currentProject.String
			}
			employeeMap[employeeID] = employee
		}

		// Add the skill
		employee.Skills = append(employee.Skills, models.Skill{
			ID:   skillID,
			Name: skillName,
		})
	}

	// Convert map to slice
	var employees []models.Employee
	for _, employee := range employeeMap {
		employees = append(employees, *employee)
	}

	fmt.Printf("DEBUG: Found %d employees with matching skills (optimized query)\n", len(employees))
	return employees, nil
}

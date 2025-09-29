package repositories

import (
	"database/sql"
	"stafind-backend/internal/models"
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
	query := r.MustGetQuery("get_all_employees")
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(
			&employee.ID, &employee.Name, &employee.Email, &employee.Department,
			&employee.Level, &employee.Location, &employee.Bio, &employee.CurrentProject,
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

		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *employeeRepository) GetByID(id int) (*models.Employee, error) {
	query := r.MustGetQuery("get_employee_by_id")
	var employee models.Employee
	err := r.db.QueryRow(query, id).Scan(
		&employee.ID, &employee.Name, &employee.Email, &employee.Department,
		&employee.Level, &employee.Location, &employee.Bio, &employee.CurrentProject,
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

func (r *employeeRepository) Create(req *models.CreateEmployeeRequest) (*models.Employee, error) {
	tx, err := r.db.Begin()
	if err != nil {
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

	err = tx.QueryRow(query, req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject).
		Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return nil, err
	}

	employee.Name = req.Name
	employee.Email = req.Email
	employee.Department = req.Department
	employee.Level = req.Level
	employee.Location = req.Location
	employee.Bio = req.Bio
	employee.CurrentProject = currentProject

	// Add skills
	for _, skillReq := range req.Skills {
		err = r.addEmployeeSkill(tx, employee.ID, &skillReq)
		if err != nil {
			return nil, err
		}
	}

	// Get the created employee with skills
	createdEmployee, err := r.GetByID(employee.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return createdEmployee, nil
}

func (r *employeeRepository) Update(id int, req *models.CreateEmployeeRequest) (*models.Employee, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := r.MustGetQuery("update_employee")
	// Convert string to pointer string for current_project
	var currentProject *string
	if req.CurrentProject != "" {
		currentProject = &req.CurrentProject
	}

	_, err = tx.Exec(query, req.Name, req.Email, req.Department, req.Level, req.Location, req.Bio, currentProject, id)
	if err != nil {
		return nil, err
	}

	// Remove existing skills
	err = r.RemoveSkills(id)
	if err != nil {
		return nil, err
	}

	// Add new skills
	for _, skillReq := range req.Skills {
		err = r.addEmployeeSkill(tx, id, &skillReq)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *employeeRepository) Delete(id int) error {
	query := r.MustGetQuery("delete_employee")
	_, err := r.db.Exec(query, id)
	return err
}

func (r *employeeRepository) GetSkills(employeeID int) ([]models.Skill, error) {
	query := r.MustGetQuery("get_employee_skills")
	rows, err := r.db.Query(query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var skill models.Skill
		var proficiencyLevel int
		var yearsExperience float64

		err := rows.Scan(&skill.ID, &skill.Name, &skill.Category, &proficiencyLevel, &yearsExperience)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

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
	query := r.MustGetQuery("get_skill_by_name")
	err := tx.QueryRow(query, skillReq.SkillName).Scan(&skillID)
	if err == sql.ErrNoRows {
		// Create the skill
		query = r.MustGetQuery("create_skill")
		err = tx.QueryRow(query, skillReq.SkillName).Scan(&skillID)
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

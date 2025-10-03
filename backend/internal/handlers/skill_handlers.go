package handlers

import (
	"encoding/json"
	"net/http"
	"stafind-backend/internal/models"
	"stafind-backend/internal/services"
	"strconv"
	"strings"
)

type SkillHandlers struct {
	skillService services.SkillService
}

// NewSkillHandlers creates a new skill handlers instance
func NewSkillHandlers(skillService services.SkillService) *SkillHandlers {
	return &SkillHandlers{
		skillService: skillService,
	}
}

// GetAllSkills handles GET /api/skills
func (h *SkillHandlers) GetAllSkills(w http.ResponseWriter, r *http.Request) {
	skills, err := h.skillService.GetAllSkills()
	if err != nil {
		http.Error(w, "Failed to fetch skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    skills,
		"count":   len(skills),
	})
}

// GetSkillByID handles GET /api/skills/{id}
func (h *SkillHandlers) GetSkillByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	skill, err := h.skillService.GetSkillByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Skill not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    skill,
	})
}

// GetSkillsByCategory handles GET /api/skills/category/{category}
func (h *SkillHandlers) GetSkillsByCategory(w http.ResponseWriter, r *http.Request) {
	category := strings.TrimPrefix(r.URL.Path, "/api/skills/category/")
	if category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	// For now, get all skills and filter by category name
	// TODO: Implement GetSkillsByCategoryName method in service
	allSkills, err := h.skillService.GetSkillsWithCategories()
	if err != nil {
		http.Error(w, "Failed to fetch skills by category", http.StatusInternalServerError)
		return
	}

	// Filter skills by category name
	var filteredSkills []models.Skill
	for _, skill := range allSkills {
		for _, skillCategory := range skill.Categories {
			if strings.EqualFold(skillCategory.Name, category) {
				filteredSkills = append(filteredSkills, skill)
				break
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"data":     filteredSkills,
		"category": category,
		"count":    len(filteredSkills),
	})
}

// GetCategories handles GET /api/skills/categories
func (h *SkillHandlers) GetCategories(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement category service integration
	http.Error(w, "Categories endpoint not yet implemented", http.StatusNotImplemented)
}

// SearchSkills handles GET /api/skills/search?q={query}
func (h *SkillHandlers) SearchSkills(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	skills, err := h.skillService.SearchSkills(query)
	if err != nil {
		http.Error(w, "Failed to search skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    skills,
		"query":   query,
		"count":   len(skills),
	})
}

// GetPopularSkills handles GET /api/skills/popular?limit={limit}
func (h *SkillHandlers) GetPopularSkills(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default limit

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	skills, err := h.skillService.GetPopularSkills(limit)
	if err != nil {
		http.Error(w, "Failed to fetch popular skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    skills,
		"limit":   limit,
		"count":   len(skills),
	})
}

// GetSkillsWithEmployeeCount handles GET /api/skills/with-count
func (h *SkillHandlers) GetSkillsWithEmployeeCount(w http.ResponseWriter, r *http.Request) {
	skills, err := h.skillService.GetSkillsWithEmployeeCount()
	if err != nil {
		http.Error(w, "Failed to fetch skills with employee count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    skills,
		"count":   len(skills),
	})
}

// CreateSkill handles POST /api/skills
func (h *SkillHandlers) CreateSkill(w http.ResponseWriter, r *http.Request) {
	var skill models.Skill
	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	createdSkill, err := h.skillService.CreateSkill(&skill)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    createdSkill,
		"message": "Skill created successfully",
	})
}

// CreateSkillsBatch handles POST /api/skills/batch
func (h *SkillHandlers) CreateSkillsBatch(w http.ResponseWriter, r *http.Request) {
	var skills []models.Skill
	if err := json.NewDecoder(r.Body).Decode(&skills); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	createdSkills, err := h.skillService.CreateSkillsBatch(skills)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    createdSkills,
		"count":   len(createdSkills),
		"message": "Skills created successfully",
	})
}

// UpdateSkill handles PUT /api/skills/{id}
func (h *SkillHandlers) UpdateSkill(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	var skill models.Skill
	if err := json.NewDecoder(r.Body).Decode(&skill); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	updatedSkill, err := h.skillService.UpdateSkill(id, &skill)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Skill not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    updatedSkill,
		"message": "Skill updated successfully",
	})
}

// UpdateSkillsBatch handles PUT /api/skills/batch
func (h *SkillHandlers) UpdateSkillsBatch(w http.ResponseWriter, r *http.Request) {
	var updates []models.SkillUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := h.skillService.UpdateSkillsBatch(updates)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to update skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Skills updated successfully",
	})
}

// DeleteSkill handles DELETE /api/skills/{id}
func (h *SkillHandlers) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid skill ID", http.StatusBadRequest)
		return
	}

	err = h.skillService.DeleteSkill(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Skill not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete skill", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Skill deleted successfully",
	})
}

// DeleteSkillsBatch handles DELETE /api/skills/batch
func (h *SkillHandlers) DeleteSkillsBatch(w http.ResponseWriter, r *http.Request) {
	var request struct {
		IDs []int `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := h.skillService.DeleteSkillsBatch(request.IDs)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to delete skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Skills deleted successfully",
	})
}

// GetSkillStats handles GET /api/skills/stats
func (h *SkillHandlers) GetSkillStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.skillService.GetSkillStats()
	if err != nil {
		http.Error(w, "Failed to fetch skill statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// SuggestSkillCategories handles GET /api/skills/suggest-categories?name={name}
func (h *SkillHandlers) SuggestSkillCategories(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Skill name is required", http.StatusBadRequest)
		return
	}

	// TODO: Implement category suggestion logic
	http.Error(w, "Category suggestion endpoint not yet implemented", http.StatusNotImplemented)
}

// GetSkillsByEmployeeID handles GET /api/skills/employee/{id}
func (h *SkillHandlers) GetSkillsByEmployeeID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/skills/employee/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	skills, err := h.skillService.GetSkillsByEmployeeID(id)
	if err != nil {
		http.Error(w, "Failed to fetch employee skills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"data":        skills,
		"employee_id": id,
		"count":       len(skills),
	})
}

// GetSkillsByEmployeeIDs handles POST /api/skills/employees
func (h *SkillHandlers) GetSkillsByEmployeeIDs(w http.ResponseWriter, r *http.Request) {
	var request struct {
		EmployeeIDs []int `json:"employee_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	skillsMap, err := h.skillService.GetSkillsByEmployeeIDs(request.EmployeeIDs)
	if err != nil {
		http.Error(w, "Failed to fetch skills for employees", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    skillsMap,
		"count":   len(skillsMap),
	})
}

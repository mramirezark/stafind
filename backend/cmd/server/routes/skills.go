package routes

import (
	"net/http"
	"stafind-backend/internal/handlers"
	"stafind-backend/internal/services"
)

// SetupSkillRoutes configures all skill-related routes
func SetupSkillRoutes(mux *http.ServeMux, skillService services.SkillService) {
	skillHandlers := handlers.NewSkillHandlers(skillService)

	// Public routes (no authentication required)
	mux.HandleFunc("/api/skills", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			skillHandlers.GetAllSkills(w, r)
		case http.MethodPost:
			// Create skill
			skillHandlers.CreateSkill(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Skill by ID routes
	mux.HandleFunc("/api/skills/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Handle different sub-paths
		switch {
		case path == "/api/skills/categories":
			if r.Method == http.MethodGet {
				skillHandlers.GetCategories(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/popular":
			if r.Method == http.MethodGet {
				skillHandlers.GetPopularSkills(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/with-count":
			if r.Method == http.MethodGet {
				skillHandlers.GetSkillsWithEmployeeCount(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/stats":
			if r.Method == http.MethodGet {
				skillHandlers.GetSkillStats(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/suggest-categories":
			if r.Method == http.MethodGet {
				skillHandlers.SuggestSkillCategories(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/search":
			if r.Method == http.MethodGet {
				skillHandlers.SearchSkills(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/batch":
			switch r.Method {
			case http.MethodPost:
				// Batch create
				skillHandlers.CreateSkillsBatch(w, r)
			case http.MethodPut:
				// Batch update
				skillHandlers.UpdateSkillsBatch(w, r)
			case http.MethodDelete:
				// Batch delete
				skillHandlers.DeleteSkillsBatch(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/employees":
			if r.Method == http.MethodPost {
				skillHandlers.GetSkillsByEmployeeIDs(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/category/":
			if r.Method == http.MethodGet {
				skillHandlers.GetSkillsByCategory(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case path == "/api/skills/employee/":
			if r.Method == http.MethodGet {
				skillHandlers.GetSkillsByEmployeeID(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			// Handle individual skill operations by ID
			switch r.Method {
			case http.MethodGet:
				skillHandlers.GetSkillByID(w, r)
			case http.MethodPut:
				// Update
				skillHandlers.UpdateSkill(w, r)
			case http.MethodDelete:
				// Delete
				skillHandlers.DeleteSkill(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})
}

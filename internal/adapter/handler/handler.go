package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sahib002/deployflow/internal/core/domain"
	"github.com/sahib002/deployflow/internal/core/ports"
	"github.com/sahib002/deployflow/internal/core/services"
)

type Handler struct {
	repo    ports.RepositoryPort
	service *services.DeploymentService // <--- NEW: The Logic Engine
}

func NewHandler(repo ports.RepositoryPort, service *services.DeploymentService) *Handler {
	return &Handler{
		repo:    repo,
		service: service,
	}
}

// CreateProject handles POST /projects
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	project := &domain.Project{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	if err := h.repo.CreateProject(r.Context(), project); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// CreateDeployment handles POST /deployments
// This is the "Magic Button" that builds the container
func (h *Handler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID string `json:"project_id"`
		Code      string `json:"code"` // The actual Python code as a string
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Parse UUID
	id, err := uuid.Parse(req.ProjectID)
	if err != nil {
		http.Error(w, "Invalid Project ID", http.StatusBadRequest)
		return
	}

	// Call the Service (The Robot Arm)
	deployment, err := h.service.CreateDeployment(r.Context(), id, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deployment)
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/projects", h.CreateProject)
	r.Post("/deployments", h.CreateDeployment) // <--- Register the new route
	return r
}

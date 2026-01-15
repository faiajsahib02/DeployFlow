package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/sahib002/deployflow/internal/core/domain"
)

// Repository implements ports.RepositoryPort
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new postgres repository
func NewRepository(dsn string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Config connection pool for performance
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Repository{db: db}, nil
}

// CreateProject inserts a new project
func (r *Repository) CreateProject(ctx context.Context, project *domain.Project) error {
	query := `INSERT INTO projects (id, name, created_at) VALUES (:id, :name, :created_at)`

	// NamedExec is a cool sqlx feature that maps struct fields to SQL params automatically
	_, err := r.db.NamedExecContext(ctx, query, project)
	return err
}

// GetProjectByName finds a project
func (r *Repository) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	var project domain.Project
	query := `SELECT * FROM projects WHERE name = $1`

	err := r.db.GetContext(ctx, &project, query, name)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	return &project, err
}

// CreateDeployment inserts a deployment record
func (r *Repository) CreateDeployment(ctx context.Context, d *domain.Deployment) error {
	query := `
		INSERT INTO deployments (id, project_id, status, container_id, port, image_tag, created_at, updated_at)
		VALUES (:id, :project_id, :status, :container_id, :port, :image_tag, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, d)
	return err
}

// UpdateDeploymentStatus updates the status only
func (r *Repository) UpdateDeploymentStatus(ctx context.Context, id uuid.UUID, status domain.DeploymentStatus) error {
	query := `UPDATE deployments SET status = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// UpdateDeployment updates the full deployment record
func (r *Repository) UpdateDeployment(ctx context.Context, d *domain.Deployment) error {
	query := `
		UPDATE deployments 
		SET status = $1, container_id = $2, port = $3, image_tag = $4, updated_at = $5 
		WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query, d.Status, d.ContainerID, d.Port, d.ImageTag, time.Now(), d.ID)
	return err
}

// GetActiveDeployment finds the latest running container for a project
func (r *Repository) GetActiveDeployment(ctx context.Context, projectID uuid.UUID) (*domain.Deployment, error) {
	var deployment domain.Deployment
	// We want the MOST RECENT (ORDER BY created_at DESC) one that is RUNNING
	query := `
		SELECT * FROM deployments 
		WHERE project_id = $1 AND status = 'running' 
		ORDER BY created_at DESC 
		LIMIT 1`

	err := r.db.GetContext(ctx, &deployment, query, projectID)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}

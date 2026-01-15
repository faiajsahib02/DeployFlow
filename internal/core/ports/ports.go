package ports

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/sahib002/deployflow/internal/core/domain"
)

// RepositoryPort defines how we talk to the Database
type RepositoryPort interface {
	CreateProject(ctx context.Context, project *domain.Project) error
	CreateDeployment(ctx context.Context, deployment *domain.Deployment) error
	UpdateDeploymentStatus(ctx context.Context, id uuid.UUID, status domain.DeploymentStatus) error
	UpdateDeployment(ctx context.Context, deployment *domain.Deployment) error

	// --- NEW METHODS ---
	GetProjectByName(ctx context.Context, name string) (*domain.Project, error)
	GetActiveDeployment(ctx context.Context, projectID uuid.UUID) (*domain.Deployment, error)
}

// RuntimePort defines how we talk to Docker
type RuntimePort interface {
	BuildImage(ctx context.Context, tag string, buildContext io.Reader) error
	RunContainer(ctx context.Context, imageTag string) (string, int, error) // Returns containerID, port, error
	StopContainer(ctx context.Context, containerID string) error
}

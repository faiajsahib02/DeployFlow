package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sahib002/deployflow/internal/core/domain"
	"github.com/sahib002/deployflow/internal/core/ports"
	"github.com/sahib002/deployflow/internal/utils"
)

type DeploymentService struct {
	repo   ports.RepositoryPort
	docker ports.RuntimePort
}

func NewDeploymentService(repo ports.RepositoryPort, docker ports.RuntimePort) *DeploymentService {
	return &DeploymentService{repo: repo, docker: docker}
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, projectID uuid.UUID, pythonCode string) (*domain.Deployment, error) {
	// 1. Create DB Record
	id := uuid.New()
	deployment := &domain.Deployment{
		ID:        id,
		ProjectID: projectID,
		Status:    domain.StatusBuilding,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateDeployment(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to create deployment record: %w", err)
	}

	// 2. Prepare Docker Context (The Dockerfile + Code)
	dockerfile := `
FROM python:3.9-slim
WORKDIR /app
RUN pip install flask
COPY app.py /app/app.py
EXPOSE 8000
CMD ["python", "app.py"]
`
	tarStream, err := utils.CreateTarArchive(map[string]string{
		"Dockerfile": dockerfile,
		"app.py":     pythonCode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create tar: %w", err)
	}

	// 3. Build Image
	imageTag := fmt.Sprintf("deployflow-%s", id.String())
	fmt.Printf("🚧 Building Image: %s...\n", imageTag)

	if err := s.docker.BuildImage(ctx, imageTag, tarStream); err != nil {
		s.repo.UpdateDeploymentStatus(ctx, id, domain.StatusFailed)
		return nil, fmt.Errorf("docker build failed: %w", err)
	}

	// 4. Run Container
	fmt.Printf("🚀 Starting Container...\n")
	containerID, port, err := s.docker.RunContainer(ctx, imageTag)
	if err != nil {
		s.repo.UpdateDeploymentStatus(ctx, id, domain.StatusFailed)
		return nil, fmt.Errorf("docker run failed: %w", err)
	}

	// 5. Success
	fmt.Printf("✅ Success! Container %s running on port %d\n", containerID[:12], port)

	// Update the full deployment record with container info
	deployment.Status = domain.StatusRunning
	deployment.ContainerID = containerID
	deployment.Port = port
	deployment.ImageTag = imageTag

	if err := s.repo.UpdateDeployment(ctx, deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}

	return deployment, nil
}

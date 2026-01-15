package domain

import (
	"time"

	"github.com/google/uuid" // You will need: go get github.com/google/uuid
)

// DeploymentStatus Enum for type safety
type DeploymentStatus string

const (
	StatusQueued   DeploymentStatus = "queued"
	StatusBuilding DeploymentStatus = "building"
	StatusRunning  DeploymentStatus = "running"
	StatusFailed   DeploymentStatus = "failed"
	StatusStopped  DeploymentStatus = "stopped"
)

// Project represents a user's repository/application
type Project struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"` // e.g., "my-flask-app"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Deployment represents a specific version/instance of a Project
type Deployment struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	ProjectID uuid.UUID        `json:"project_id" db:"project_id"`
	Status    DeploymentStatus `json:"status" db:"status"`

	// Infrastructure Details (The "Meat" of the project)
	ContainerID string `json:"container_id" db:"container_id"` // Docker Container ID
	Port        int    `json:"port" db:"port"`                 // The random port on Host (e.g. 49154)
	ImageTag    string `json:"image_tag" db:"image_tag"`       // docker.io/my-flask-app:v1

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

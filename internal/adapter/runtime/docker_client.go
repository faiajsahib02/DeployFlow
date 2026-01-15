package runtime

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type DockerClient struct {
	cli *client.Client
}

// NewDockerClient creates a connection to the local Docker Daemon
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to docker: %w", err)
	}
	return &DockerClient{cli: cli}, nil
}

// BuildImage builds a docker image from a tar stream
func (d *DockerClient) BuildImage(ctx context.Context, tag string, buildContext io.Reader) error {
	// Docker requires the tar stream to be passed here
	resp, err := d.cli.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Tags:       []string{tag},
		Dockerfile: "Dockerfile", // This must match the filename inside the Tar
		Remove:     true,         // Remove intermediate containers
	})
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	// Crucial: We must read the response stream to ensure the build finishes
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

// RunContainer starts a container from an image
func (d *DockerClient) RunContainer(ctx context.Context, imageTag string) (string, int, error) {
	// 1. Configure Host: Bind container port 8000 to a random Host Port
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "0", // 0 = Random Port
				},
			},
		},
	}

	// 2. Create
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: imageTag,
		ExposedPorts: nat.PortSet{
			"8000/tcp": struct{}{},
		},
	}, hostConfig, nil, nil, "")
	if err != nil {
		return "", 0, fmt.Errorf("failed to create container: %w", err)
	}

	// 3. Start
	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", 0, fmt.Errorf("failed to start container: %w", err)
	}

	// 4. Inspect to find the assigned port
	inspect, err := d.cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to inspect container: %w", err)
	}

	ports := inspect.NetworkSettings.Ports["8000/tcp"]
	if len(ports) > 0 {
		var port int
		fmt.Sscanf(ports[0].HostPort, "%d", &port)
		return resp.ID, port, nil
	}

	return "", 0, fmt.Errorf("port not assigned")
}

// StopContainer kills a running container
func (d *DockerClient) StopContainer(ctx context.Context, containerID string) error {
	timeout := 5 // seconds
	return d.cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

package e2esuite

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Utility for pulling and using an image of mulberry
func PullMulberryImage(image string) error {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Context for the Docker API call
	ctx := context.Background()

	// Pull the Docker image
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer out.Close()

	// Print the output from the pull operation
	if _, err := io.Copy(os.Stdout, out); err != nil {
		return fmt.Errorf("failed to read pull output: %w", err)
	}

	fmt.Println("\nImage pulled successfully")
	return nil
}

// RunContainer creates and starts a container from the given image.
func RunContainer(image string, containerName string) (string, error) {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Context for the Docker API calls
	ctx := context.Background()

	// Create the container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			// Command or entrypoint (optional)
			Cmd: []string{"sleep", "3600"}, // Example: long-running command
		},
		&container.HostConfig{
			AutoRemove: true, // Automatically remove container after it stops
		},
		nil, // NetworkingConfig
		nil, // Platform
		containerName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	fmt.Printf("Container started with ID: %s\n", resp.ID)
	return resp.ID, nil
}

// ExecCommandInContainer executes a command inside a running container
func ExecCommandInContainer(containerID string, command []string) error {
	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Context for the Docker API calls
	ctx := context.Background()

	// Create an exec instance in the container
	execConfig := types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	}
	execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return fmt.Errorf("failed to create exec instance: %w", err)
	}

	// Start the exec process
	resp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{})
	if err != nil {
		return fmt.Errorf("failed to attach to exec instance: %w", err)
	}
	defer resp.Close()

	// Stream the command output
	if _, err := io.Copy(os.Stdout, resp.Reader); err != nil {
		return fmt.Errorf("failed to read exec output: %w", err)
	}

	return nil

}

func StreamContainerLogs(containerID string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	ctx := context.Background()
	options := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true}

	out, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return fmt.Errorf("failed to fetch container logs: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(os.Stdout, out) // Stream logs to the local terminal
	return err
}

package e2esuite

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
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

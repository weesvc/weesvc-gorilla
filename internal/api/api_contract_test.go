package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/k6"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type serviceContainer struct {
	testcontainers.Container
	Host string
	Port string
}

// TestAPIContract runs the API compliance script against our service to ensure contracts.
func TestAPIContract(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	t.Parallel()

	ctx := context.Background()

	// Build the image and start a container with our service.
	service, err := buildServiceContainer(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	// Start a k6 container to run our API compliance test.
	k6c, err := k6.RunContainer(
		ctx,
		k6.WithCache(),
		// TODO Temporarily holding script in repository until k6 container can retrieve from a remote url.
		k6.WithTestScript("../testhelpers/api-compliance.js"),
		// k6.WithTestScript("https://raw.githubusercontent.com/weesvc/workbench/main/scripts/api-compliance.js"),
		k6.SetEnvVar("HOST", service.Host),
		k6.SetEnvVar("PORT", service.Port),
	)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if kerr := k6c.Terminate(ctx); kerr != nil {
			t.Fatalf("failed to terminate k6 container: %s", kerr)
		}
	})
}

// buildServiceContainer will build and start our service within a container based on current source.
func buildServiceContainer(ctx context.Context, t *testing.T) (*serviceContainer, error) {
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				FromDockerfile: testcontainers.FromDockerfile{
					Context:       "..",
					Dockerfile:    "Dockerfile",
					PrintBuildLog: false,
					KeepImage:     false,
				},
				ExposedPorts: []string{"9092/tcp"},
				Cmd:          []string{"/bin/sh", "-c", "/app/weesvc migrate; /app/weesvc serve"},
				WaitingFor:   wait.ForHTTP("/api/hello").WithStartupTimeout(10 * time.Second),
			},
			Started: true,
		},
	)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		if terr := container.Terminate(ctx); terr != nil {
			t.Fatalf("failed to terminate ServiceContainer: %s", terr)
		}
	})

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "9092")
	if err != nil {
		return nil, err
	}

	return &serviceContainer{
		Container: container,
		Host:      ip,
		Port:      mappedPort.Port(),
	}, nil
}

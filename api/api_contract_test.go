package api

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type serviceContainer struct {
	testcontainers.Container
	URI string
}

func TestAPIContract(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	t.Parallel()

	// Build the image and start a container with our service.
	svcContainer, err := buildServiceContainer(t)
	if err != nil {
		t.Fatal(err)
	}

	// TODO Start a k6 container to run our compliance script
	//nolint:noctx
	resp, err := http.Get(svcContainer.URI + "/api/places")
	assert.NoError(t, err)

	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// buildServiceContainer will build and start our service within a container based on current source.
func buildServiceContainer(t *testing.T) (*serviceContainer, error) {
	ctx := context.Background()

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
		URI:       "http://" + net.JoinHostPort(ip, mappedPort.Port()),
	}, nil
}

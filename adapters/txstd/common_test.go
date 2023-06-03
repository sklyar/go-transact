//go:build integration

package txstd

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // used by migrator
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     = "test_db"
	dbUser     = "postgres"
	dbPassword = "password"
)

type Container struct {
	ConnectionStr string

	pc *postgres.PostgresContainer
}

func (c *Container) Close(ctx context.Context) error {
	return c.pc.Container.Terminate(ctx)
}

func startContainer(ctx context.Context) (*Container, error) {
	waitForLogs := wait.
		ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)

	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(waitForLogs),
	)
	if err != nil {
		return nil, err
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	return &Container{
		ConnectionStr: connStr,
		pc:            container,
	}, nil
}

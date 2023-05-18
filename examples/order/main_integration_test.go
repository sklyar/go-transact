package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sklyar/go-transact"
	"github.com/sklyar/go-transact/adapters/txstd"
	"github.com/sklyar/go-transact/txsql"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     = "test_db"
	dbUser     = "postgres"
	dbPassword = "password"
)

var (
	db        txsql.DB
	txManager *transact.Manager
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

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container, err := startContainer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close(ctx)

	sqlDB, err := sql.Open("pgx", container.ConnectionStr)
	if err != nil {
		log.Fatal(err)
	}

	txManager, db, err = transact.NewManager(txstd.Wrap(sqlDB))
	if err != nil {
		log.Fatal(err)
	}

	initb, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatalf("failed to read init.sql: %v", err)
	}

	_, err = db.Exec(ctx, string(initb))
	if err != nil {
		log.Fatalf("failed to execute init.sql: %v", err)
	}

	os.Exit(m.Run())
}

func TestIntegration_App_CreateOrder(t *testing.T) {
	t.Parallel()

	orderRepo := &orderRepository{db: db}
	inventoryRepo := &inventoryRepository{db: db}
	orderService := NewOrderService(txManager, orderRepo, inventoryRepo)

	customerID := 1
	products := []int{1, 2}

	ctx := context.Background()
	orderID, err := orderService.Create(ctx, customerID, products)
	assert.NoError(t, err)
	assert.Equal(t, 1, orderID)
}

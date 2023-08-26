//go:build integration

package transactstd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sklyar/go-transact"
	"github.com/sklyar/go-transact/txsql"
	"github.com/stretchr/testify/require"
)

var (
	db        txsql.DB
	txManager *transact.Manager
)

type TestEntity struct {
	ID    uint64
	Value string
}

func newTestEntity(id uint64, value string) *TestEntity {
	return &TestEntity{
		ID:    id,
		Value: value,
	}
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

	txManager, db, err = transact.NewManager(Wrap(sqlDB))
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestDatabase_Exec(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_exec"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	assertRowsCount(t, ctx, table, 1)
}

func TestDatabase_ExecInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_exec_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)

	ctx, tx, err := txManager.Begin(ctx)
	require.NoError(t, err)

	insertTestData(ctx, t, table, entity)

	ctx, err = tx.Commit(ctx)
	require.NoError(t, err)

	assertRowsCount(t, ctx, table, 1)
}

func TestDatabase_FailedExecInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_failed_exec_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)

	someErr := fmt.Errorf("some error")
	err := txManager.BeginFunc(ctx, func(tx context.Context) error {
		insertTestData(tx, t, table, entity)

		return someErr
	})
	require.ErrorContains(t, err, someErr.Error())

	assertRowsCount(t, ctx, table, 0)
}

func TestDatabase_Query(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_query"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	rows, err := db.Query(ctx, "SELECT * FROM test_query")
	require.NoError(t, err)

	assertRowsValues(t, rows, entity)
}

func TestDatabase_QueryInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_query_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	ctx, tx, err := txManager.Begin(ctx)
	require.NoError(t, err)

	rows, err := db.Query(ctx, "DELETE FROM test_query_in_tx RETURNING *")
	require.NoError(t, err)

	assertRowsValues(t, rows, entity)

	ctx, err = tx.Commit(ctx)
	require.NoError(t, err)

	assertRowsCount(t, ctx, table, 0)
}

func TestDatabase_FailedQueryInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_failed_query_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	someErr := fmt.Errorf("some error")
	err := txManager.BeginFunc(ctx, func(tx context.Context) error {
		rows, err := db.Query(tx, "SELECT * FROM test_failed_query_in_tx")
		require.NoError(t, err)

		assertRowsValues(t, rows, entity)

		res, err := db.Exec(tx, "DELETE FROM test_failed_query_in_tx WHERE id = 1")
		require.NoError(t, err)

		assertResult(t, res)

		return someErr
	})
	require.ErrorContains(t, err, someErr.Error())

	assertRowsCount(t, ctx, table, 1)
}

func TestDatabase_QueryRow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_query_row"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	row := db.QueryRow(ctx, "SELECT * FROM test_query_row")
	assertRowValues(t, row, entity)
}

func TestDatabase_QueryRowInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_query_row_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	txctx, tx, err := txManager.Begin(ctx, txsql.WithIsolationLevel(txsql.LevelRepeatableRead))
	require.NoError(t, err)

	queryRow := db.QueryRow(txctx, "SELECT * FROM test_query_row_in_tx ORDER BY id DESC")
	assertRowValues(t, queryRow, entity)

	// insert another entity after transaction started,
	// so it will not be visible in transaction.
	entity2 := newTestEntity(2, "test2")
	insertTestData(txctx, t, table, entity2)

	queryRow = db.QueryRow(ctx, "SELECT * FROM test_query_row_in_tx ORDER BY id DESC")
	assertRowValues(t, queryRow, entity)

	_, err = tx.Commit(txctx)
	require.NoError(t, err)
}

func TestDatabase_FailedQueryRowInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_failed_query_row_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	insertTestData(ctx, t, table, entity)

	someErr := fmt.Errorf("some error")
	err := txManager.BeginFunc(ctx, func(tx context.Context) error {
		queryRow := db.QueryRow(tx, "SELECT * FROM test_failed_query_row_in_tx")
		assertRowValues(t, queryRow, entity)

		res, err := db.Exec(tx, "DELETE FROM test_failed_query_row_in_tx WHERE id = 1")
		require.NoError(t, err)

		assertResult(t, res)

		return someErr
	})
	require.ErrorContains(t, err, someErr.Error())

	assertRowsCount(t, ctx, table, 1)
}

func TestDatabase_Prepare(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_prepare"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)
	stmt, err := db.Prepare(ctx, "INSERT INTO test_prepare (name) VALUES ($1)")
	require.NoError(t, err)

	result, err := stmt.Exec(entity.Value)
	require.NoError(t, err)

	assertResult(t, result)
}

func TestDatabase_PrepareInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_prepare_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)

	ctxTx, tx, err := txManager.Begin(ctx)
	require.NoError(t, err)

	stmt, err := db.Prepare(ctxTx, "INSERT INTO test_prepare_in_tx (name) VALUES ($1)")
	require.NoError(t, err)

	result, err := stmt.Exec(entity.Value)
	require.NoError(t, err)

	assertResult(t, result)

	// check that entity is not visible in another transaction.
	assertRowsCount(t, ctx, table, 0)

	ctx, err = tx.Commit(ctxTx)
	require.NoError(t, err)

	assertRowsCount(t, ctx, table, 1)
}

func TestDatabase_FailedPrepareInTx(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const table = "test_failed_prepare_in_tx"
	entity := newTestEntity(1, "test")

	setupTable(ctx, t, table)

	someErr := fmt.Errorf("some error")
	err := txManager.BeginFunc(ctx, func(tx context.Context) error {
		stmt, err := db.Prepare(tx, "INSERT INTO test_failed_prepare_in_tx (name) VALUES ($1)")
		require.NoError(t, err)

		result, err := stmt.Exec(entity.Value)
		require.NoError(t, err)

		assertResult(t, result)

		return someErr
	})
	require.ErrorContains(t, err, someErr.Error())

	assertRowsCount(t, ctx, table, 0)
}

func TestDatabase_Ping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	err := db.Ping(ctx)
	require.NoError(t, err)
}

func setupTable(ctx context.Context, t *testing.T, tableName string) {
	t.Helper()

	query := fmt.Sprintf("CREATE TABLE %s (id SERIAL PRIMARY KEY, name TEXT)", tableName)
	_, err := db.Exec(ctx, query)
	require.NoError(t, err)
}

func insertTestData(ctx context.Context, t *testing.T, tableName string, entity *TestEntity) {
	t.Helper()

	query := fmt.Sprintf("INSERT INTO %s (id, name) VALUES ($1, $2)", tableName)
	res, err := db.Exec(ctx, query, entity.ID, entity.Value)
	require.NoError(t, err)

	assertResult(t, res)
}

func assertResult(t *testing.T, res txsql.Result) {
	t.Helper()

	n, err := res.RowsAffected()
	require.NoError(t, err)
	require.EqualValues(t, 1, n)
}

func assertRowValues(t *testing.T, row txsql.Row, expEntity *TestEntity) {
	t.Helper()

	var (
		id   int
		name string
	)

	err := row.Scan(&id, &name)
	require.NoError(t, err)

	require.EqualValues(t, expEntity.ID, id)
	require.EqualValues(t, expEntity.Value, name)
}

func assertRowsValues(t *testing.T, rows txsql.Rows, expEntities ...*TestEntity) {
	t.Helper()

	for i := 0; i < len(expEntities); i++ {
		require.True(t, rows.Next())

		var (
			id   int
			name string
		)
		err := rows.Scan(&id, &name)
		require.NoError(t, err)

		require.EqualValues(t, expEntities[i].ID, id)
		require.EqualValues(t, expEntities[i].Value, name)
	}

	require.NoError(t, rows.Close())
}

func assertRowsCount(t *testing.T, ctx context.Context, tableName string, expCount int) {
	t.Helper()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	row := db.QueryRow(ctx, query)
	require.NoError(t, row.Err())

	var count int
	require.NoError(t, row.Scan(&count))
	require.EqualValues(t, expCount, count)
}

package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/romakorinenko/task-manager/internal/config"
	dbpool2 "github.com/romakorinenko/task-manager/internal/dbpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBForTest struct {
	PostgresContainer *PostgresContainer
	DBPool            *pgxpool.Pool
}

type PostgresContainer struct {
	Container        *postgres.PostgresContainer
	ConnectionString string
}

func CreateDBForTest(t *testing.T, migrationsDir string) *DBForTest {
	t.Helper()

	_, filename, line, ok := runtime.Caller(0)
	_ = line
	if !ok {
		log.Fatalln("cannot receive root project filename")
	}
	projectDir := path.Join(path.Dir(filename), "..")

	postgresContainer := RunPostgresContainer(t)
	dbpool, err := dbpool2.NewDBPool(context.Background(), &config.DB{
		ConnectionString: postgresContainer.ConnectionString,
	})
	require.NoError(t, err)
	db := stdlib.OpenDBFromPool(dbpool)

	err = UpMigrations(db, projectDir+migrationsDir)
	require.NoError(t, err)

	return &DBForTest{
		DBPool:            dbpool,
		PostgresContainer: postgresContainer,
	}
}

func RunPostgresContainer(t *testing.T) *PostgresContainer {
	t.Helper()

	ctx := context.Background()

	dbName := "test_db"
	dbUser := "test_user"
	dbPassword := "test_db_password123321"

	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	return &PostgresContainer{
		Container:        postgresContainer,
		ConnectionString: connectionString,
	}
}

func UpMigrations(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, dir)
}

func (t *DBForTest) Close() {
	err := t.PostgresContainer.Container.Terminate(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

package tests

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var (
	ctx               context.Context
	once              sync.Once
	postgresContainer *postgres.PostgresContainer
)

const (
	user     = "root"
	password = "root"
	port     = "5432"
	database = "root"
)

func TestWithPostgres(m *testing.M) {
	once.Do(startPostgresContainer)
	exitCode := m.Run()
	stopPostgresContainer()
	os.Exit(exitCode)
}

func startPostgresContainer() {
	ctx = context.Background()

	initScriptPath, err := getInitScriptPath()
	if err != nil {
		log.Fatal(err)
	}
	postgresContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithInitScripts(initScriptPath),
		postgres.WithDatabase(database),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}
}

func getInitScriptPath() (string, error) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(projectRoot, "migrations", "users.sql"), nil
}

func getProjectRoot() (string, error) {
	var err error
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err = os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "", os.ErrNotExist
		}

		cwd = parent
	}
}

func stopPostgresContainer() {
	if err := postgresContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to stop container: %v", err)
	}
}

func getDataSourceName() (string, error) {
	if postgresContainer == nil {
		return "", errors.New("postgres container not started")
	}

	newPort, err := nat.NewPort("tcp", port)
	if err != nil {
		return "", err
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, newPort)
	if err != nil {
		return "", err
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return "", err
	}

	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", host, mappedPort.Int()),
		Path:   database,
	}

	q := u.Query()
	q.Add("sslmode", "disable")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

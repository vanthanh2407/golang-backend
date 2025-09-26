package mysql

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func mustStartMySQLContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := mysql.Run(context.Background(),
		"mysql:8.0.36",
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUser),
		mysql.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	dbname = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "3306/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartMySQLContainer()
	if err != nil {
		log.Fatalf("could not start mysql container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mysql container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}

func TestUserCRUD(t *testing.T) {
	// Create a new database instance for this test
	dbInstance = nil // Reset the singleton
	srv := New()
	ctx := context.Background()

	// Test CreateUser
	user, err := srv.CreateUser(ctx, "testuser", "test@example.com", "password123")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.Username != "testuser" {
		t.Fatalf("expected username to be 'testuser', got %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Fatalf("expected email to be 'test@example.com', got %s", user.Email)
	}

	if user.ID <= 0 {
		t.Fatalf("expected user ID to be positive, got %d", user.ID)
	}

	// Test GetUserByID
	retrievedUser, err := srv.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get user by ID: %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Fatalf("expected username to match, got %s", retrievedUser.Username)
	}

	// Test GetUserByEmail
	userByEmail, err := srv.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	if userByEmail.ID != user.ID {
		t.Fatalf("expected user ID to match, got %d", userByEmail.ID)
	}

	// Test GetUserByUsername
	userByUsername, err := srv.GetUserByUsername(ctx, "testuser")
	if err != nil {
		t.Fatalf("failed to get user by username: %v", err)
	}

	if userByUsername.ID != user.ID {
		t.Fatalf("expected user ID to match, got %d", userByUsername.ID)
	}

	// Test GetAllUsers
	users, err := srv.GetAllUsers(ctx)
	if err != nil {
		t.Fatalf("failed to get all users: %v", err)
	}

	if len(users) == 0 {
		t.Fatalf("expected at least one user")
	}

	// Test UpdateUser
	updatedUser, err := srv.UpdateUser(ctx, user.ID, "updateduser", "updated@example.com")
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	if updatedUser.Username != "updateduser" {
		t.Fatalf("expected username to be updated, got %s", updatedUser.Username)
	}

	if updatedUser.Email != "updated@example.com" {
		t.Fatalf("expected email to be updated, got %s", updatedUser.Email)
	}

	// Test UpdateUserPassword
	err = srv.UpdateUserPassword(ctx, user.ID, "newpassword123")
	if err != nil {
		t.Fatalf("failed to update password: %v", err)
	}

	// Test DeleteUser
	err = srv.DeleteUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify user is deleted
	_, err = srv.GetUserByID(ctx, user.ID)
	if err == nil {
		t.Fatalf("expected user to be deleted")
	}
}

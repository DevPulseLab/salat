package repositories

import (
	"testing"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/db/repositories/testutils"
	"github.com/DevPulseLab/salat/internal/enum"
)

func TestRegisterUser(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	err := repo.RegisterUser("testuser", "password", "admin")
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	err = repo.RegisterUser("testuser", "password", "admin")
	if err == nil {
		t.Fatalf("user not exists: %v", err)
	}
}

func TestAuthenticateUser(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	_, err := repo.AuthenticateUser("testuser", "password")
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("authentication failed: %v", err)
	}

	repo.RegisterUser("testuser", "password", "admin")

	role, err := repo.AuthenticateUser("testuser", "password")
	if err != nil {
		t.Fatalf("authentication failed: %v", err)
	}

	if role != "admin" {
		t.Fatalf("wrong user role: %s", role)
	}

	_, err = repo.AuthenticateUser("testuser", "wrong-password")
	if err == nil {
		t.Fatalf("authentication succeeded, but should have failed")
	}
}

func TestGetAllUsers(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	results := repo.GetAllUsers()
	if len(results) != 0 {
		t.Fatalf("wrong number of users: %d", len(results))
	}

	repo.RegisterUser("testuser1", "password", "admin")
	repo.RegisterUser("testuser2", "password", "admin")

	results = repo.GetAllUsers()
	if len(results) != 2 {
		t.Fatalf("wrong number of users: %d", len(results))
	}
}

func TestGetIdByUsername(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	repo.RegisterUser("testuser", "password", "admin")

	_, err := repo.GetIdByUsername("testuser")
	if err != nil {
		t.Fatalf("failed to get user id: %v", err)
	}

	_, err = repo.GetIdByUsername("wrong-testuser")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestFindByUsername(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	repo.RegisterUser("testuser", "password", "admin")

	_, err := repo.FindByUsername("testuser")
	if err != nil {
		t.Fatalf("failed to find user by name: %v", err)
	}

	_, err = repo.FindByUsername("wrong-testuser")
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestFindById(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	repo.RegisterUser("testuser", "password", "admin")

	_, err := repo.FindById(1)
	if err != nil {
		t.Fatalf("failed to find user by id: %v", err)
	}

	_, err = repo.FindById(20)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestSetPenaltyCard(t *testing.T) {
	db := testutils.GetTestDb(t, &models.User{})

	repo := NewUserRepository(db)

	repo.RegisterUser("testuser", "password", "admin")

	userModel, _ := repo.FindById(1)
	if userModel.PenaltyCard != "" {
		t.Fatalf("penalty card should be empty before set")
	}

	repo.SetPenaltyCard(userModel.ID, string(enum.Yellow))
	userModel, _ = repo.FindById(1)
	if userModel.PenaltyCard != string(enum.Yellow) {
		t.Fatalf("penalty card should be empty before set")
	}

	err := repo.SetPenaltyCard(400, string(enum.Yellow))
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

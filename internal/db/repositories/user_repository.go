package repositories

import (
	"errors"

	"github.com/DevPulseLab/salat/internal/db/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) RegisterUser(username, password, role string) error {
	var user models.User
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err == nil {
		return errors.New("user already exists")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := models.User{
		Username: username,
		Password: string(hashedPass),
		Role:     role,
	}

	return repo.DB.Create(&newUser).Error
}

func (repo *UserRepository) AuthenticateUser(username, password string) (string, error) {
	var user models.User
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	return user.Role, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (repo *UserRepository) GetAllUsers() []models.User {
	var users []models.User

	repo.DB.Find(&users)

	return users
}

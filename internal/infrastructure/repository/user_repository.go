package repository

import (
	"errors"
	"gorm.io/gorm"
	"live-coding/internal/domain"
	"time"
)

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	GetAll() ([]domain.UserResponse, error)
	FindById(userId uint) (*domain.User, error)
	UpdateLastAccess(userId uint, LastAccess time.Time) error
	Create(user domain.User) error
	Update(userId uint, name, email, password string) error
	Delete(userId uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.Table("users").
		Select("id", "email", "password", "role_id").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAll() ([]domain.UserResponse, error) {

	var users []domain.UserResponse
	err := r.db.Table("users").
		Select("users.id,users.email,users.password,users.role_id,users.last_access,roles.name as role_name").
		Joins("JOIN roles on users.role_id = roles.id").
		Where("users.deleted_at IS NULL").
		Find(&users).Error

	return users, err
}

func (r *userRepository) FindById(userId uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Table("users").
		Select("users.id,users.email,users.password,users.role_id, users.last_access").
		Where("users.deleted_at IS NULL").
		Where("id = ?", userId).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *userRepository) UpdateLastAccess(userId uint, LastAccess time.Time) error {
	timeFormat := LastAccess.Format("15:04:05")

	return r.db.Model(
		&domain.User{}).
		Where("id = ?", userId).
		Update("last_access", timeFormat).
		Error
}

func (r *userRepository) Create(user domain.User) error {
	userToSave := struct {
		ID        uint `gorm:"primarykey"`
		RoleId    uint
		Name      string
		Email     string
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}{
		ID:        uint(user.ID),
		RoleId:    uint(user.RoleId),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.db.Table("users").Create(&userToSave).Error
}

func (r *userRepository) Update(userId uint, name, email, password string) error {

	updatesColumn := map[string]interface{}{}

	if name != "" {
		updatesColumn["name"] = name
	}
	if email != "" {
		updatesColumn["email"] = email
	}
	if password != "" {
		updatesColumn["password"] = password
	}

	if len(updatesColumn) == 0 {
		return errors.New("no update field")
	}

	return r.db.Model(&domain.User{}).
		Where("id = ?", userId).
		Updates(updatesColumn).
		UpdateColumn("updated_at", time.Now()).
		Error

}

func (r *userRepository) Delete(userId uint) error {

	return r.db.Model(&domain.User{}).
		Where("id = ?", userId).
		UpdateColumn("deleted_at", time.Now()).
		Error

}

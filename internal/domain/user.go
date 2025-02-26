package domain

import "time"

type User struct {
	ID         int64     `gorm:"primary_key;column:id"`
	RoleId     int64     `gorm:"column:role_id"`
	Name       string    `gorm:"column:name"`
	Email      string    `gorm:"column:email"`
	Password   string    `gorm:"column:password"`
	LastAccess string    `gorm:"column:last_access"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	DeletedAt  time.Time `gorm:"column:deleted_at"`
}

type UserResponse struct {
	ID         int64  `gorm:"primary_key;column:id"`
	RoleId     int64  `gorm:"column:role_id"`
	RoleName   string `gorm:"column:role_name"`
	Name       string `gorm:"column:name"`
	Email      string `gorm:"column:email"`
	Password   string `gorm:"column:password"`
	LastAccess string `gorm:"column:last_access"`
}

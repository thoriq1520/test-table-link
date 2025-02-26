package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type RoleRightsRepository interface {
	CheckPermission(roleId int64, route string, permission string) (bool, error)
}

type roleRightsRepository struct {
	db *gorm.DB
}

func NewRoleRightsRepository(db *gorm.DB) RoleRightsRepository { return &roleRightsRepository{db: db} }

func (r *roleRightsRepository) CheckPermission(roleId int64, route string, permission string) (bool, error) {

	var value bool
	query := fmt.Sprintf("SELECT %s FROM role_rights WHERE role_id = ? AND route = ?", permission)
	err := r.db.Raw(query, roleId, route).Scan(&value).Error
	if err != nil {
		return false, err
	}
	return value, nil
}

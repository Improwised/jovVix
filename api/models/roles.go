package models

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

type Role string

type AllowedRoles struct {
	roles map[Role]any
}

// QuizModel implements quiz related database operations
type RoleModel struct {
	db          *goqu.Database
	systemRoles []Role
}

func InitRoleModel(db *goqu.Database) *RoleModel {
	return &RoleModel{
		db:          db,
		systemRoles: []Role{"admin", "user"},
	}
}

func (rm *RoleModel) NewAllowedRoles(roles ...string) (AllowedRoles, error) {
	allowedRoles := AllowedRoles{roles: make(map[Role]any)}
	validRoles := []Role{}

	for _, r := range roles {
		role := Role(r)
		matched := false
		for _, ra := range rm.systemRoles {
			if ra == role {
				validRoles = append(validRoles, role)
				matched = true
			}
		}
		if !matched {
			return AllowedRoles{}, fmt.Errorf("Role not found: %s", r)
		}
	}

	for _, role := range validRoles {
		allowedRoles.roles[role] = any(nil)
	}

	return allowedRoles, nil
}

func (ar *AllowedRoles) IsAllowed(role Role) bool {
	_, found := ar.roles[role]
	return found
}

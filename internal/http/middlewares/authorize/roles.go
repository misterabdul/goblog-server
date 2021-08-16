package authorize

import "github.com/misterabdul/goblog-server/internal/models"

type UserRole struct {
	Level int
	Name  string
}

func GetAllRoles() []*UserRole {
	roles := []*UserRole{}

	roles = append(roles, &UserRole{
		Level: 0,
		Name:  "SuperAdmin",
	})
	roles = append(roles, &UserRole{
		Level: 1,
		Name:  "Admin",
	})
	roles = append(roles, &UserRole{
		Level: 2,
		Name:  "Editor",
	})
	roles = append(roles, &UserRole{
		Level: 3,
		Name:  "Writer",
	})

	return roles
}

func GetRole(name string) *UserRole {
	switch {
	default:
		fallthrough
	case name == "Writer":
		return &UserRole{
			Level: 3,
			Name:  "Writer",
		}
	case name == "Editor":
		return &UserRole{
			Level: 2,
			Name:  "Editor",
		}
	case name == "Admin":
		return &UserRole{
			Level: 1,
			Name:  "Admin",
		}
	case name == "SuperAdmin":
		return &UserRole{
			Level: 0,
			Name:  "SuperAdmin",
		}
	}
}

func CheckRoles(me *models.UserModel, role *UserRole) bool {
	for _, userRole := range me.Roles {
		if userRole.Level == role.Level {
			return true
		}
	}
	return false
}

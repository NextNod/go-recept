package models

import "gorm.io/gorm"

type UserDB struct {
	gorm.Model
	Name  string
	Image string
	GID   string
}

type UserOut struct {
	Name  string
	Image string
}

func NewUserOut(user UserDB) UserOut {
	return UserOut{
		Name:  user.Name,
		Image: user.Image,
	}
}

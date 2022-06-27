package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Recept struct {
	gorm.Model
	Author      uint
	Name        string
	Description string
}

type ReceptProduct struct {
	gorm.Model
	ReceptId uint
	ProducId int
}

type ImageReceptProduct struct {
	gorm.Model
	Image     string
	ProductId uint
	ReceptId  uint
}

type ProductIn struct {
	Name   string
	Images []string
}

type ReceptIn struct {
	Name        string
	Description string
	Author      string
	Products    []int
	Images      []string
}

type ImageResponse struct {
	ImageID string
}

type ReceptResponse struct {
	ID          uint
	Name        string
	Description string
	Author      User
	Images      []ImageReceptProduct
	Products    []ProductResponse
}

type ProductResponse struct {
	ID     uint
	Name   string
	Images []ImageReceptProduct
}

type BaseResponse struct {
	Result any
	Error  string
}

type User struct {
	gorm.Model
	Name  string
	Image string
	GID   string
}

type UserIn struct {
	Name  string
	Image string
	GID   string `gorm:"unique"`
}

type UserResponse struct {
	User any
	New  bool
}

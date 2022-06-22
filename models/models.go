package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Recept struct {
	gorm.Model
	Name string
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
	Name     string
	Products []int
	Images   []string
}

type ImageResponse struct {
	ImageID string
}

type ReceptResponse struct {
	ID       uint
	Name     string
	Images   []ImageReceptProduct
	Products []ProductResponse
}

type ProductResponse struct {
	ID     uint
	Name   string
	Images []ImageReceptProduct
}

type BaseResponse struct {
	Result any
}

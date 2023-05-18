package models

import "gorm.io/gorm"

type ImageReceptProductDB struct {
	gorm.Model
	Image     string
	ProductId uint
	ReceptId  uint
}

type ImageReceptProductOut struct {
	Image     string
	ProductId uint
	ReceptId  uint
}

func NewImageReceptProductOut(item ImageReceptProductDB) ImageReceptProductOut {
	return ImageReceptProductOut{
		Image:     item.Image,
		ProductId: item.ProductId,
		ReceptId:  item.ReceptId,
	}
}

func NewAImageReceptProductOut(item []ImageReceptProductDB) []ImageReceptProductOut {
	var result []ImageReceptProductOut
	for i := range item {
		result = append(result, NewImageReceptProductOut(item[i]))
	}
	return result
}

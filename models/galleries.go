package models

import "github.com/jinzhu/gorm"

// Gallery holds per-gallery data
type Gallery struct {
	gorm.Model
	UserID uint `gorm:"not_null;index"`
	Title string `gorm:"not_null"`

}

// GalleryService interface ... [TODO: add documentation]
type GalleryService interface {
	GalleryDB
}

// GalleryDB interface ... [TODO: add documentation]
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	// TODO: implement this
	return nil
}
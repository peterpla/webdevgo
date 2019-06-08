package models

import "github.com/jinzhu/gorm"

// Services holds service details fro each of our services
type Services struct {
	Gallery GalleryService
	User UserService
}

// NewServices opens the database connection and initializes each service
func NewServices(connectionInfo string) (*Services, error) {
	// open the database connection
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	// initialize the User and Gallery services
	s := &Services{
		User: NewUserService(db),
		Gallery: &galleryGorm{},
	}
	return s, nil
}
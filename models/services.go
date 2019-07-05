package models

import "github.com/jinzhu/gorm"

// Services holds service details fro each of our services
type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
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
		User:    NewUserService(db),
		Gallery: &galleryGorm{},
		db:      db,
	}
	return s, nil
}

// Close closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

// AutoMigrate will attempt to automagically migrate all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

// DestructiveReset drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

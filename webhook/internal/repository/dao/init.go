package dao

import "gorm.io/gorm"

func InitUserTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func InitUserProfileTables(db *gorm.DB) error {
	return db.AutoMigrate(&UserProfile{})
}

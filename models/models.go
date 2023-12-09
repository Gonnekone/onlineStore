package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	Email    string `gorm:"unique"`
	Password string
	Role     string `gorm:"default:USER"`
	Basket   Basket
	Ratings  []Rating
}

type Basket struct {
	ID            uint `gorm:"primarykey"`
	UserID        uint
	BasketDevices []BasketDevice
}

type BasketDevice struct {
	ID       uint `gorm:"primarykey"`
	BasketID uint
	Device   Device
	DeviceID uint
}

type Device struct {
	ID      uint    `gorm:"primarykey"`
	Name    string  `gorm:"unique;not null"`
	Price   float64 `gorm:"not null"`
	Rating  float32 `gorm:"default:0"`
	Img     string  `gorm:"not null"`
	TypeID  uint
	Type    Type
	BrandID uint
	Brand   Brand
	Ratings []Rating
	Info    []DeviceInfo
}

type Type struct {
	ID      uint   `gorm:"primarykey"`
	Name    string `gorm:"unique;not null"`
	Devices []Device
}

type Brand struct {
	ID      uint   `gorm:"primarykey"`
	Name    string `gorm:"unique;not null"`
	Devices []Device
}

type Rating struct {
	ID       uint    `gorm:"primarykey"`
	Rate     float32 `gorm:"not null"`
	UserID   uint
	DeviceID uint
}

type DeviceInfo struct {
	ID          uint   `gorm:"primarykey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	DeviceID    uint
}

type TypeBrand struct {
	ID      uint `gorm:"primarykey"`
	TypeID  uint
	BrandID uint
}

// AutoMigrate creates or updates the database tables
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Basket{}, &BasketDevice{}, &Device{}, &Type{}, &Brand{}, &Rating{}, &DeviceInfo{}, &TypeBrand{})
	if err != nil {
		panic("failed to migrate database")
	}
}

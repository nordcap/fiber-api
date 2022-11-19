package models

import "time"

type Order struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	ProductRefer uint    `json:"product_id"`
	Product      Product `gorm:"foreignKey:ProductRefer"`
	UserRefer    uint    `json:"user_id"`
	User         User    `gorm:"foreignKey:UserRefer"`
}

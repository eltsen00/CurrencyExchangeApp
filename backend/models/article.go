package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title" binding:"required"`
	Content string `gorm:"not null" json:"content" binding:"required"`
	Preview string `gorm:"not null" json:"preview" binding:"required"`
	Like    int    `gorm:"default:0" json:"like"`
}

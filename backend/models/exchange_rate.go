package models

import "time"

type ExchangeRate struct {
	ID           int       `gorm:"primaryKey" json:"_id"`
	FromCurrency string    `gorm:"not null" json:"fromCurrency" binding:"required"`
	ToCurrency   string    `gorm:"not null" json:"toCurrency" binding:"required"`
	Rate         float64   `gorm:"not null" json:"rate" binding:"required"`
	Date         time.Time `json:"date"`
}

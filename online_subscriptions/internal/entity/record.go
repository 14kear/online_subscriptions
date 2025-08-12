package entity

import "time"

type Record struct {
	ID          uint      `gorm:"primaryKey"`
	ServiceName string    `gorm:"not null"`
	Price       int       `gorm:"not null;check:price >= 0"`
	UserID      string    `gorm:"not null;index"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP'"`
	ExpiresAt   time.Time `gorm:"not null"`
}

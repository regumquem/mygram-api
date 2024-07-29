package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Password     string         `gorm:"not null" json:"password"`
	Age          int            `gorm:"not null" json:"age"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Photos       []Photo        `gorm:"foreignKey:UserID"`
	Comments     []Comment      `gorm:"foreignKey:UserID"`
	SocialMedias []SocialMedia  `gorm:"foreignKey:UserID"`
}

type Photo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `gorm:"not null" json:"photo_url"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `gorm:"foreignKey:PhotoID"`
}

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	PhotoID   uint      `json:"photo_id"`
	Message   string    `gorm:"not null" json:"message"`
	User      User      `gorm:"foreignKey:UserID"`
	Photo     Photo     `gorm:"foreignKey:PhotoID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SocialMedia struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	SocialMediaURL string    `gorm:"not null" json:"social_media_url"`
	UserID         uint      `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

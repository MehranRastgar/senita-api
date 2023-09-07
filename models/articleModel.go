package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	ID         int64      `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
	Title      string     `json:"title" gorm:"not null"`
	Content    string     `json:"content" gorm:"not null"`
	AuthorID   int64      `json:"author_id" gorm:"not null"`
	Author     User       `json:"author" gorm:"foreignKey:AuthorID"`
	CategoryID *int64     `json:"category_id"`
	Category   Category   `json:"category" gorm:"foreignKey:CategoryID"`
}

func (article *Article) BeforeCreate(tx *gorm.DB) error {
	article.ID = int64(uuid.New().ID())
	return nil
}

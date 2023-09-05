package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        int64      `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Name      string     `json:"name" gorm:"not null"`
	ParentID  *int64     `json:"parent_id"`
	Parent    *Category  `json:"parent" gorm:"foreignKey:ParentID"`
	Children  []Category `json:"children" gorm:"foreignKey:ParentID"`
}

func (category *Category) BeforeCreate(tx *gorm.DB) error {
	category.ID = int64(uuid.New().ID())
	return nil
}

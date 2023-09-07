package models

import (
	"senita-api/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRoleClient UserRole = "client"
	UserRoleVendor UserRole = "vendor"
)

type User struct {
	ID        int64      `json:"id" gorm:"unique;not null;index;primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Email     string     `json:"email" gorm:"unique;not null;index:idx_user_email"`
	Password  string     `json:"password"`
	UserName  string     `json:"user_name" gorm:"index:idx_name,unique"`
	Role      UserRole   `json:"role"`
	VendorId  int64      `json:"vendor_id"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = int64(uuid.New().ID())
	user.Password = utils.HashAndSalt([]byte(user.Password))
	if user.Role == "" {
		user.Role = UserRoleClient
	}
	return nil
}

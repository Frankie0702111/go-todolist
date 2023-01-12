package model

import (
	"time"
)

// Create User struct representing the user table in the database
// type User struct {
// 	ID        uint64    `gorm:"primary_key:auto_increment" json:"id"`
// 	Username  string    `gorm:"NOT NULL;type:varchar(255)" json:"username"`
// 	Email     string    `gorm:"NOT NULL;uniqueIndex;type:varchar(255)" json:"email"`
// 	Password  string    `gorm:"->;<-;NOT NULL;type:varchar(255)" json:"-"`
// 	Token     string    `gorm:"-" json:"token,omitempty"`
// 	CreatedAt time.Time `gorm:"NOT NULL;DEFAULT CURRENT_TIMESTAMP;type:timestamp" json:"created_at"`
// 	UpdatedAt time.Time `gorm:"NOT NULL;DEFAULT CURRENT_TIMESTAMP;type:timestamp" json:"updated_at"`
// }

// Create User struct representing the user table in the database
type User struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Token     string    `gorm:"-" json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Token struct {
	Token string `json:"token"`
}

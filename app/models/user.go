package models

// User represents a User schema
type User struct {
	Base
	Email    string `json:"email" gorm:"unique;index;size:255"`
	Username string `json:"username" gorm:"unique;index;size:255"`
	Password string `json:"-" gorm:"size:255"`
}

package models

// User represents a User schema
type User struct {
	Base
	Email    string `json:"email" gorm:"unique;index"`
	Username string `json:"username" gorm:"unique;index"`
	Password string `json:"-"`
}

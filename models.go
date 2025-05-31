package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
}

type Message struct {
	gorm.Model
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Content    string    `gorm:"not null" json:"content"`
	Status     string    `gorm:"default:'sent'" json:"status"` // sent, delivered, read
	Sender     User      `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver   User      `gorm:"foreignKey:ReceiverID" json:"receiver"`
	Timestamp  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
}

// HashPassword hashes a password using bcrypt
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword checks if the provided password matches the hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

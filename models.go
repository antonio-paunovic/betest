package main

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Dob         string `json:"dob" gorm:"type:datetime"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

func (u *User) createUser(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) getUser(db *gorm.DB) error {
	return db.First(&u, u.ID).Error
}

func dbGetUsers(users *[]User, db *gorm.DB) error {
	return db.Find(users).Error
}

func (u *User) updateUser(db *gorm.DB) error {
	return db.Save(&u).Error
}

func (u *User) deleteUser(db *gorm.DB) error {
	return db.Delete(u).Error
}

package main

type User struct {
  UID int `gorm:"primaryKey"`
  Username string `gorm:"unique"`
  Password string `gorm:"not null"`
}

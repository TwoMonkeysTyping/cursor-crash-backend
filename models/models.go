package models

type Document struct {
	ID		string `gorm:"primaryKe"`
	content	string
}

type User struct {
    ID       uint   `gorm:"primarykey"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}
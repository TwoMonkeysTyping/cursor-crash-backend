package models

type Document struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	IsSaved   bool		`gorm:"default:false" json:"isSaved"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type User struct {
    ID       uint   `gorm:"primarykey"`
    Email    string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
}
package model

import "time"

type PostM struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;not null"`
	PostID    string    `gorm:"column:postID;not null"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (u *PostM) TableName() string {
	return "post"
}

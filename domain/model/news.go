package model

import (
	"github.com/jinzhu/gorm"
)

// News represent entity of the news
type News struct {
	gorm.Model
	Title   string  `json:"title"`
	Slug    string  `json:"slug"`
	Content string  `json:"content" gorm:"text"`
	Status  string  `json:"status"`
	Topic   []Topic `gorm:"many2many:news_topics;"`
}
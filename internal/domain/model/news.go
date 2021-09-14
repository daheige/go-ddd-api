package model

// News represent entity of the news
type News struct {
	BaseModel
	Title   string  `json:"title"`
	Slug    string  `json:"slug"`
	Content string  `json:"content" gorm:"text"`
	Status  string  `json:"status"`
	Topic   []Topic `gorm:"many2many:news_topics;"`
}

// TableName table name
func (News) TableName() string {
	return "news"
}

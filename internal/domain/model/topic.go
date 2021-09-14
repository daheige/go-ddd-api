package model

// Topic represent entity of the topic
type Topic struct {
	BaseModel
	Name string `json:"name"`
	Slug string `json:"slug"`
	News []News `gorm:"many2many:news_topics;" json:"news"`
}

// TableName table name
func (Topic) TableName() string {
	return "topics"
}

package model

type AuthorizedIP struct {
    BaseModel
    IP      string `gorm:"unique;not null"`
    Disable bool   `gorm:"default:true"`
}

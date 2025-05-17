package model

type AuthorizedIP struct {
    BaseModel
    IP       string `db:"ip"       json:"ip"`
    Disabled bool   `db:"disabled" json:"disabled"`
}
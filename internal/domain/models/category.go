package models

type Category struct {
	Id            int    `db:"category_id" json:"id"`
	Name          string `db:"name" json:"name"`
	Description   string `db:"description" json:"description"`
	PhotoUrl      string `db:"photo_url" json:"photoUrl"`
	BackgroundUrl string `db:"background_url" json:"backgroundUrl"`
}

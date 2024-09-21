package models

type Category struct {
	Id          int    `db:"category_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	PhotoUrl    string `db:"photo_url"`
}

package models

type Course struct {
	Id            int    `db:"course_id"`
	Name          string `db:"name"`
	Description   string `db:"description"`
	MinRating     int    `db:"min_rating"`
	OptimalRating int    `db:"optimal_rating"`
	CategoryId    int    `db:"category_id"`
	PhotoUrl      string `db:"photo_url"`
}

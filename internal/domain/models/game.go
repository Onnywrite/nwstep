package models

import "time"

type Game struct {
	Id       int        `db:"game_id"`
	CourseId int        `db:"course_id"`
	CreateAt time.Time  `db:"create_at"`
	EndAt    *time.Time `db:"end_at"`
}

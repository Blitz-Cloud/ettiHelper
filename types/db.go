package types

import (
	"time"

	"gorm.io/gorm"
)

type Lab struct {
	gorm.Model
	Title              string
	Description        string
	Date               *time.Time
	Tags               string
	UniYearAndSemester uint
	Content            string
}
type Blog Lab

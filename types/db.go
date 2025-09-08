package types

import (
	"time"

	"gorm.io/gorm"
)

type Lab struct {
	gorm.Model
	// frontmatter
	Title              string
	Description        string
	Date               *time.Time
	Tags               string
	Subject            string
	UniYearAndSemester uint
	// fisierul de md
	Content string
}
type Blog Lab

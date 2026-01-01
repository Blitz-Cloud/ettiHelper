package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type Lab struct {
// 	gorm.Model
// 	// frontmatter
// 	Title              string
// 	Description        string
// 	Date               *time.Time
// 	Tags               string
// 	Subject            string
// 	UniYearAndSemester uint
// 	// fisierul de md
// 	Content     string
// 	CodeExample string
// }
// type Blog Lab

type UUIDBase struct {
	ID        string `gorm:"primaryKey;size:36"` // UUID is 36 chars
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (base *UUIDBase) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		base.ID = uuid.New().String()
	}
	return nil
}

// Namespace
type Namespace struct {
	UUIDBase
	Name string `gorm:"uniqueIndex;not null"`

	// Relationships
	Categories []Category `gorm:"foreignKey:NamespaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Category
type Category struct {
	UUIDBase
	NamespaceID string `gorm:"size:36;index"` // Foreign Key is now a String
	Name        string `gorm:"not null"`

	Protected bool `gorm:"default:false"`
	Visible   bool `gorm:"default:true"`

	// Relationships
	Namespace Namespace
	Posts     []Post `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Post
type Post struct {
	UUIDBase
	CategoryID string `gorm:"size:36;index"` // Foreign Key is now a String

	// Slug               string `gorm:"index"`
	Title              string
	Description        string
	PublishedDate      *time.Time
	UniYearAndSemester int
	Content            string `gorm:"type:text"`
	Hash               string `gorm:"uniqueIndex;not null"`
	Protected          bool
	RestrictedTo       string
	Visible            bool

	// Relationships
	Category Category
	// Tags     []Tag `gorm:"many2many:post_tags;"`
}

// Tag
// type Tag struct {
// 	UUIDBase
// 	Name  string `gorm:"uniqueIndex"`
// 	Posts []Post `gorm:"many2many:post_tags;"`
// }

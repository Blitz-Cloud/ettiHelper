package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/Blitz-Cloud/ettiHelper/types"
	"github.com/davecgh/go-spew/spew"
	"gorm.io/gorm"
)

func SeedFromInMemory(sqlDB *gorm.DB, memDB DB) error {
	return sqlDB.Transaction(func(tx *gorm.DB) error {
		namespaceName2IdMap := make(map[string]string)
		categoryName2IdMap := make(map[string]string)
		for _, dbName := range memDB.DBNames {
			nameSpace := types.Namespace{
				Name: dbName,
			}

			if err := tx.Where("name = ?", dbName).FirstOrCreate(&nameSpace).Error; err != nil {
				Log.Error(err.Error())
			}
			namespaceName2IdMap[dbName] = nameSpace.ID
		}

		for _, category := range memDB.Categories {
			Log.Dump(category)
			metaData := strings.Split(category.Name, "-")
			if len(metaData) != 2 {
				continue
			}
			categoryNamespace := metaData[0]
			dbCategory := types.Category{}

			if err := tx.Where(types.Category{
				NamespaceID: namespaceName2IdMap[categoryNamespace],
				Name:        category.Name,
				Protected:   category.Protected,
				Visible:     category.Visible,
			}).FirstOrCreate(&dbCategory).Error; err != nil {
				spew.Dump(err)
			}
			categoryName2IdMap[category.Name] = dbCategory.ID
		}

		for _, post := range memDB.Posts {
			postDB := types.Post{}
			metaData := strings.Split(post.Category, "-")
			if len(metaData) != 2 {
				continue
			}
			Log.Debug("Id to look up : %s", post.ID)
			err := tx.Where("id = ? ", post.ID).First(&postDB).Error

			if err == nil && postDB.Hash != post.Hash {
				postDB.UUIDBase.ID = post.ID
				postDB.CategoryID = categoryName2IdMap[post.Category]
				postDB.Hash = post.Hash
				postDB.Title = post.Title
				postDB.Description = post.Description
				postDB.UniYearAndSemester = post.UniYearAndSemester
				postDB.Content = post.Content
				postDB.Protected = post.Properties.Protected
				postDB.Visible = post.Properties.Visible
				if err := tx.Save(&postDB).Error; err != nil {
					Log.Error(err.Error())
				}

			} else if errors.Is(err, gorm.ErrRecordNotFound) {

				t, err := time.Parse(time.RFC3339, post.Date)
				if err != nil {
					t = time.Now().UTC().Local()
				}
				newPost := types.Post{}
				postDB.Hash = post.Hash
				newPost.UUIDBase.ID = post.ID
				newPost.CategoryID = categoryName2IdMap[post.Category]
				newPost.Title = post.Title
				newPost.Description = post.Description
				newPost.PublishedDate = &t
				newPost.UniYearAndSemester = post.UniYearAndSemester
				newPost.Content = post.Content
				newPost.Protected = post.Properties.Protected
				newPost.Visible = post.Properties.Visible
				if err := tx.Create(&newPost).Error; err != nil {
					Log.Error(err.Error())
				}
			}

		}

		return nil
	})
}

func parseTime(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	// Adjust layout based on your frontmatter format (e.g., YYYY-MM-DD)
	layout := "2006-01-02"

	// If your code uses RFC3339 for LastUpdated:
	if len(dateStr) > 10 {
		layout = time.RFC3339
	}

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return nil
	}
	return &t
}

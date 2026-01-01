package utils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gofiber/fiber/v2/log"
)

type FrontmatterMetaData struct {
	ID                 string   `yaml:"id"`
	Title              string   `yaml:"title" json:"title"`
	Date               string   `yaml:"date" json:"date"`
	Description        string   `yaml:"description" json:"description"`
	Tags               []string `yaml:"tags" json:"tags"`
	UniYearAndSemester int      `yaml:"uniYearAndSemester" json:"uniYearAndSemester"`
}

type Post struct {
	// ID string
	// frontmatter
	FrontmatterMetaData
	Category string
	Content  string `json:"content"`
	Hash     string
	Properties
}

type Category struct {
	Name      string
	Protected bool
	Visible   bool
}

type DB struct {
	DBNames    []string
	Categories []Category
	Posts      []Post
}

func ParseMdString(data string) (Post, error) {
	var post Post
	mdContent, err := frontmatter.Parse(strings.NewReader(data), &post.FrontmatterMetaData)
	if err != nil {
		return Post{}, err
	}
	post.Content = string(mdContent)
	return post, nil
}

type Properties struct {
	Protected    bool   `json:"protected"`
	RestrictedTo string `json:"restricted-to"`
	// LastUpdated  string `json:"last-updated"`
	Visible bool `json:"visible"`
}

func readPostFile(path string) (Post, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}
	hash := sha256.Sum256(fileContent)
	post, err := ParseMdString(string(fileContent))
	if err != nil {
		return Post{}, err
	}
	post.Hash = fmt.Sprintf("%x", hash)
	Log.Info("Hash here: %s", post.Hash)
	Log.Dump(post)
	return post, nil
}

func readProprieties(path string) (Properties, error) {
	proprieties, err := os.ReadFile(filepath.Join(path, ".proprieties"))
	if errors.Is(err, fs.ErrNotExist) {
		return Properties{}, fmt.Errorf("The proprieties file couldn't be found at %s This path will be ignored", path)
	} else if err != nil {
		return Properties{}, err
	}

	props := Properties{}
	err = json.Unmarshal(proprieties, &props)
	if err != nil {
		return Properties{}, fmt.Errorf("Failed to parsed proprieties:\n%s\n%s", string(proprieties), err.Error())
	}
	return props, nil
}

func getDbNames(paths ...string) []string {
	dbNames := make([]string, 0)
	for _, path := range paths {
		dbNames = append(dbNames, path[strings.LastIndex(path, "/")+1:])
	}
	return dbNames
}

func getCategories(paths ...string) []Category {
	categories := make([]Category, 0)
	for _, path := range paths {
		dirs, err := os.ReadDir(path)
		if err != nil {
			return []Category{}
		}
		dbName := path[strings.LastIndex(path, "/")+1:]
		for _, dir := range dirs {
			if dir.IsDir() {
				props, err := readProprieties(filepath.Join(path, dir.Name()))
				if err != nil {
					log.Error(err)
				}
				categories = append(categories, Category{
					Name:      fmt.Sprintf("%s-%s", dbName, dir.Name()),
					Protected: props.Protected,
					Visible:   props.Visible,
				})
			}
			// log.Debug(props)
			// categories = append(categories, fmt.Sprintf("%s-%s", dir.Name(), dbName))

		}
	}
	return categories
}

func getPosts(categories []Category, paths ...string) []Post {
	rootPath := paths[0][:strings.LastIndex(paths[0], "/")]
	Posts := []Post{}
	for _, category := range categories {
		dirPath := filepath.Join(rootPath, strings.Replace(category.Name, "-", "/", 1))
		posts, err := os.ReadDir(dirPath)
		if err != nil {
			Log.Error(err.Error())
		}
		for _, post := range posts {
			if post.Name() != ".proprieties" {

				fullPost, err := readPostFile(filepath.Join(dirPath, post.Name()))
				if err != nil {
					Log.Error(err.Error())
					return []Post{}
				}
				fullPost.Category = category.Name
				Posts = append(Posts, fullPost)
			}
		}
	}
	return Posts
}

func InMemoryDB(paths ...string) (DB, error) {
	DB := DB{}
	DB.DBNames = getDbNames(paths...)
	DB.Categories = getCategories(paths...)
	DB.Posts = getPosts(DB.Categories, paths...)
	return DB, nil
}

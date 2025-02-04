package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type File struct {
	name     string
	location string
	content  string
	parent   *FsNode
}

type FsNode struct {
	name     string
	location string
	dirs     []*FsNode
	files    []*File
	parent   *FsNode
}

type Example struct {
	Location string
	Date     string
	Name     string
	Content  string
}

func ExplorerLegacy(location string, node *FsNode, examples *[]Example) {
	dirContent, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirContent {
		if file.IsDir() {
			newDirNode := FsNode{
				name:     file.Name(),
				location: fmt.Sprintf("%s/%s", location, file.Name()),
				parent:   node,
			}
			node.dirs = append(node.dirs, &newDirNode)
			ExplorerLegacy(newDirNode.location, &newDirNode, examples)
		} else if path.Ext(file.Name()) == ".c" {

			// location needed in order to read the file contents
			fileLocation := fmt.Sprintf("%s/%s", location, file.Name())
			fileContent, err := os.ReadFile(fileLocation)

			if err != nil {
				log.Fatal(err)
			}

			newFile := File{
				name:     file.Name(),
				location: fileLocation,
				content:  string(fileContent),
				parent:   node,
			}

			date, _ := time.Parse("2_Jan_2006", (node.parent).name)

			newExample := Example{
				Location: fileLocation,
				Name:     node.name,
				Date:     fmt.Sprintf("%d-%s-%d", date.Day(), date.Month().String()[:3], date.Year()),
				Content:  strings.Replace(strings.Trim(string(fileContent), " "), "\n", "", 1),
			}
			*examples = append((*examples), newExample)

			node.files = append(node.files, &newFile)
		}
	}
}

type FactoryFunc[T any] func(*File, *[]T)

func CopyPasteParser(file *File, contentArray *[]Example) {
	if file.name == "main.c" {

		unParsedDate := ((file.parent).parent).name
		date, err := time.Parse("2_Jan_2006", unParsedDate)
		if err != nil {
			log.Fatal(err)
		}
		spew.Dump(date)
		newExample := Example{
			Location: file.location,
			Name:     (file.parent).name,
			Date:     fmt.Sprintf("%d-%s-%d", date.Day(), date.Month().String()[:3], date.Year()),
			Content:  strings.Replace(strings.Trim(string(file.content), " "), "\n", "", 1),
		}
		*contentArray = append((*contentArray), newExample)
	}
}

func Explorer[T any](location string, node *FsNode, extension string, contentArray *[]T, parser FactoryFunc[T]) {
	dirContent, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirContent {
		if file.IsDir() {
			newDirNode := FsNode{
				name:     file.Name(),
				location: fmt.Sprintf("%s/%s", location, file.Name()),
				parent:   node,
			}
			node.dirs = append(node.dirs, &newDirNode)
			Explorer(newDirNode.location, &newDirNode, extension, contentArray, parser)
		} else if path.Ext(file.Name()) == extension {
			fileLocation := fmt.Sprintf("%s/%s", location, file.Name())
			fileContent, err := os.ReadFile(fileLocation)

			if err != nil {
				log.Fatal(err)
			}

			newFile := File{
				name:     file.Name(),
				location: fileLocation,
				content:  string(fileContent),
				parent:   node,
			}
			node.files = append(node.files, &newFile)
			parser(&newFile, contentArray)
		}

	}
}

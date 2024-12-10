package main

import (
	"fmt"
	"log"
	"os"
	"time"
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

func Explorer(location string, node *FsNode, examples *[]Example) {
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
			Explorer(newDirNode.location, &newDirNode, examples)
		} else if file.Name() == "main.c" {
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
				Content:  string(fileContent),
			}
			*examples = append((*examples), newExample)
			node.files = append(node.files, &newFile)
		}
	}
}

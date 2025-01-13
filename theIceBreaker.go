package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

type Tipizat struct {
	Location     string
	Date         string
	Name         string
	Content      string
	LinkCompiler string
}

func TipizatExplorer(location string, node *FsNode, examples *[]Tipizat) {
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
			TipizatExplorer(newDirNode.location, &newDirNode, examples)
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

			rootFolder := (node.parent).name
			fmt.Println(rootFolder)
			newExample := Tipizat{
				Location:     fileLocation,
				Name:         node.name,
				Date:         rootFolder,
				Content:      strings.Replace(strings.Trim(string(fileContent), " "), "\n", "", 1),
				LinkCompiler: fmt.Sprintf("<a href='https://cpp.sh/?source=%s' class='text-blue-200' target='_blank'> Ruleaza codul cu cpp.sh </a>", url.QueryEscape(strings.Replace(strings.Trim(string(fileContent), " "), "void main", "int main", 1))),
			}
			*examples = append((*examples), newExample)

			node.files = append(node.files, &newFile)
		}
	}
}

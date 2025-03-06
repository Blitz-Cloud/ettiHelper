package utils

import (
	"fmt"
	"log"
	"os"
	"path"
)

type File struct {
	Name     string
	Location string
	Content  string
	Parent   *FsNode
}

type FsNode struct {
	Name     string
	Location string
	dirs     []*FsNode
	Files    []*File
	Parent   *FsNode
}

type FactoryFunc[T any] func(*File, *[]T)

func MdContentParser(file *File, contentArray *[]BlogPost) {
	metaData, content := ParseMdString(file.Content)
	post := BlogPost{
		FrontmatterMetaData{
			metaData.Title,
			metaData.Date,
			metaData.Description,
			metaData.Tags,
			metaData.UniYearAndSemester,
		},
		file.Location,
		content,
	}
	*contentArray = append((*contentArray), post)
}

// explorer este menit sa gaseasca in mod recursiv toate fisierele cu o anumita extensie si sa
// creeze in memorie un arbore al memorie
// acesta functie cel mai probabil va mai suferi modificari majore odata cu adaugarea metaDatelor
// pentru fiecare post insa acolo unde acestea nu pot fi create si sau adaugate aceasta functie
// isi face treaba
// pe langa ceea ce am scris mai sus, prin intermediul unei singure rulari, este posibil,
// crearea si extragerea datelor necesare din fiecare fiesier
// aceasta functie se realizeaza prin contentArray si functia parser care pot primi orice tip de
// date
func Explorer[T any](location string, node *FsNode, extension string, contentArray *[]T, parser FactoryFunc[T]) {
	dirContent, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dirContent {
		// atat timp cat gasim un folder acesta va fi adaugat in arbore
		// daca extensia este cea dorita atunci fisierul este citit si apoi parserul este invocat
		if file.IsDir() {
			newDirNode := FsNode{
				Name:     file.Name(),
				Location: fmt.Sprintf("%s/%s", location, file.Name()),
				Parent:   node,
			}
			node.dirs = append(node.dirs, &newDirNode)
			Explorer(newDirNode.Location, &newDirNode, extension, contentArray, parser)
		} else if path.Ext(file.Name()) == extension {
			fileLocation := fmt.Sprintf("%s/%s", location, file.Name())
			fileContent, err := os.ReadFile(fileLocation)

			if err != nil {
				log.Fatal(err)
			}

			newFile := File{
				Name:     file.Name(),
				Location: fileLocation,
				Content:  string(fileContent),
				Parent:   node,
			}
			node.Files = append(node.Files, &newFile)
			parser(&newFile, contentArray)
		}

	}
}

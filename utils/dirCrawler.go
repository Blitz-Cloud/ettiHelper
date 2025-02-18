package utils

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
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
	// Location string
	Name    string
	Date    string
	Content string
}

type Tipizat struct {
	// Location     string
	Name         string
	Date         string
	Tags         []string
	Content      string
	LinkCompiler string
}

type BlogPost struct {
	Title       string
	Date        string
	Tags        []string
	Description string
	Content     string
	HtmlContent string
}

type FactoryFunc[T any] func(*File, *[]T)

// acesta este un parser de continut si este folosit pentru a extrage datele necesare despre un
// laborator
// data este exrasa dupa urmatorul model /data/numeleExercitiu/main.c
func LabsContentParser(file *File, contentArray *[]BlogPost) {
	if file.name == "main.c" {
		// se efectueaza citirea si procesarea datei
		unParsedDate := ((file.parent).parent).name
		date, err := time.Parse("2_Jan_2006", unParsedDate)

		if err != nil {
			log.Fatal(err)
		}
		newExample := BlogPost{
			// Location: file.location,
			Title:   (file.parent).name,
			Date:    fmt.Sprintf("%d-%s-%d", date.Day(), date.Month().String()[:3], date.Year()),
			Content: strings.Replace(strings.Trim(string(file.content), " "), "\n", "", 1),
		}
		*contentArray = append((*contentArray), newExample)
	}
}

// acesta este un parser de continut si este folosit pentru a extrage datele despre tipizatele
// de pe primul semestru la IETTI PCLP
// bucla for de mai jos este folosita pentru a nu adauga de mai multe ori acelasi tipizat
func ClangCodeExamplesParser(file *File, contentArray *[]BlogPost) {
	if path.Ext(file.name) == ".c" {
		rootFolder := (file.parent).name
		newTipizat := BlogPost{
			// Location:     file.location,
			Title:   file.name,
			Date:    rootFolder,
			Content: strings.Replace(strings.Trim(string(file.content), " "), "\n", "", 1),
			// LinkCompiler: fmt.Sprintf("<a href='https://cpp.sh/?source=%s' class='text-ctp-mauve' target='_blank'> Ruleaza codul cu cpp.sh </a>", url.QueryEscape(strings.Replace(strings.Trim(string(file.content), " "), "void main", "int main", 1))),
		}
		ok := 1
		for i := 0; i < len(*contentArray); i++ {
			if (*contentArray)[i].Title == newTipizat.Title {
				ok = 0
			}
		}
		if ok == 1 {
			*contentArray = append((*contentArray), newTipizat)
		}
	}
}

func MdContentParser(file *File, contentArray *[]BlogPost) {
	metaData, content := ParseMdString(file.content)
	post := BlogPost{
		Title:       metaData.Title,
		Date:        metaData.Date,
		Description: metaData.Description,
		Tags:        metaData.Tags,
		Content:     content,
		HtmlContent: string(Md2Html([]byte(content))),
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

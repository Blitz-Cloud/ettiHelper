package utils

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Depado/bfchroma/v2"
	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma/v2"
	bf "github.com/russross/blackfriday/v2"
)

type FrontmatterMetaData struct {
	Title              string   "yaml:'title'"
	Date               string   "yaml:'date'"
	Description        string   "yaml:'description'"
	Tags               []string "yaml:'tags'"
	UniYearAndSemester int      "yaml:'uniYearAndSemester'"
}

type BlogPost struct {
	FrontmatterMetaData
	Location string
	Content  string
}

func ParseMdString(data string) (FrontmatterMetaData, string) {
	var frontmatterData FrontmatterMetaData
	mdContent, err := frontmatter.Parse(strings.NewReader(data), &frontmatterData)
	if err != nil {
		fmt.Println("Couldnt process post data")
	}
	return frontmatterData, string(mdContent)
}

func Md2Html(md []byte) []byte {
	data, err := os.ReadFile("./static/catppucinMocha.xml")
	if err != nil {
		log.Fatal("Couldn't read the xml file")
	}

	style := chroma.MustNewXMLStyle(strings.NewReader(string(data)))
	return bf.Run([]byte(md), bf.WithRenderer(bfchroma.NewRenderer(bfchroma.ChromaStyle(style))))
}
func SortBlogPostsInDescendingOrderByDate(examples *[]BlogPost) {
	sort.Slice(*examples, func(i, j int) bool {
		date1, _ := time.Parse("2-Jan-2006", (*examples)[i].Date)
		date2, _ := time.Parse("2-Jan-2006", (*examples)[j].Date)
		if date1.After(date2) {
			return true
		} else {
			return false
		}
	})
}

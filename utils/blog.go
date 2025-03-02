package utils

import (
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
)

type FrontmatterMetaData struct {
	Title              string   "yaml:'title'"
	Date               string   "yaml:'date'"
	Description        string   "yaml:'description'"
	Tags               []string "yaml:'tags'"
	UniYearAndSemester int      "yaml:'uniYearAndSemester'"
}

func ParseMdString(data string) (FrontmatterMetaData, string) {
	var frontmatterData FrontmatterMetaData
	mdContent, err := frontmatter.Parse(strings.NewReader(data), &frontmatterData)
	if err != nil {
		fmt.Println("Couldnt process post data")
	}
	return frontmatterData, string(mdContent)
}

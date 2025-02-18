package utils

import (
	"sort"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// nu cred ca este necesar sa folosesc generics in acest caz deoarece momentan sortez doar in functie de data asa ca nu ma deranjeaza repetiria
// plus ca ar mari complexitatea in mod inutil
// ce ar fi necesar este sa fac parsingul inante, si sa il salvez in structura

func SortExamplesInDescendingOrderByDate(examples *[]Example) {
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
func Md2Html(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

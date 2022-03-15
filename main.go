package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

var ALLOWED_DOMAINS []string = []string{}["researchworks.oclc.org", "archives.chadwyck.com", "www.newspapers.com"]
var ARCHIVE_GRID_URL_PATTERNS []string = {
	"https://researchworks.oclc.org/archivegrid/?q=%s&limit=100",
}

// https://archives.chadwyck.com/marketing/index.jsp
// https://www.newspapers.com/
// https://researchworks.oclc.org/archivegrid/
// https://en.wikipedia.org/wiki/Names_of_the_American_Civil_Warhttps://researchworks.oclc.org/archivegrid/
// 
//	"https://researchworks.oclc.org/archivegrid/?q=Jack+Hester++and+%28%22diary%22+OR+%22journal%22+OR+%22notebook%22%29&limit=100"
// Jack+Hester++and+%28%22diary%22+OR+%22journal%22+OR+%22notebook%22%29
// 


func main() {
	c := colly.NewCollector(colly.AllowedDomains(ALLOWED_DOMAINS[0]))

	c.OnHTML("div", func(h *colly.HTMLElement) {
		contents := h.ChildAttrs("a", "href")
		fmt.Println(contents)
	})

	c.Visit(ARCHIVE_GRID_URL_PATTERNS )

}

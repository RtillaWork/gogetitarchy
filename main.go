package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

var ALLOWED_DOMAINS []string = []string{"researchworks.oclc.org", "archives.chadwyck.com", "www.newspapers.com"}
var ARCHIVE_GRID_URL_PATTERNS []string = []string{
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
// Good samples
// https://researchworks.oclc.org/archivegrid/?q=%22Albert+Quincy+Porter%22
// using person.name and `AND` :
// https://researchworks.oclc.org/archivegrid/?q=person.name%3APorter+AND+person.name%3AAlbert++AND+person.name%3AQuincy&limit=100
// also George Bowen, also Christian Abraham Fleetwood
// https://researchworks.oclc.org/archivegrid/?p=1&q=event.name%3A%22american+civil+war%22
//
// https://researchworks.oclc.org/archivegrid/?q=person.name%3APorter+AND+person.name%3AAlbert+++AND+person.name%3AQuincy&limit=100

type AGOrganization struct {
	orgId               int
	name                string
	contact_information string
}

type AGDomPaths struct {
	Record                           string // AGRecord.Dom
	Record_title                     string // AGRecordTitle.Dom
	Record_author                    string // AGRecordAuthor.Dom
	Record_archive                   string // AGRecordArchive.Dom
	Record_summary                   string // AGRecordSummary.Dom
	Record_links_contact_information string // AGRecordLinksContactInformation.Dom
}

type AGRecord struct {
	Dom       string
	URL_rec_x string
}

type AGRecordTitle struct {
	Dom         string
	href        string
	description string
}

type AGRecordAuthor struct {
	Dom         string
	href        string
	description string
}

type AGRecordArchive struct {
	Dom         string
	href        string
	description string
}

type AGRecordSummary struct {
	Dom         string
	href        string
	description string
}

type AGRecordLinksContactInformation struct {
	Dom         string
	href        string
	description string
}

type ArchiveGridRecords struct {
	RecId                            int
	Record                           AGRecord
	Record_title                     AGRecordTitle
	Record_author                    AGRecordAuthor
	Record_archive                   AGRecordArchive
	Record_summary                   AGRecordSummary
	Record_links_contact_information AGRecordLinksContactInformation
}

/*
div.searchresult
	div #rec_x .record
		input type="hidden" #url_rec_x value="/archivegrid/collection/data/nnnnnnnn"
		div itemprop="name" .record_title
			h3
				a
				href="/archivegrid/collection/data/same"
					$here text collection data title
				/a

		div itemprop="author" .record_author
			span itemprop="name"
				$here text author

		div itemprop="contributor" .record_archive
			span itemprop="name"
				$here text archive name

		div .record_summary
			$here text summary

		div .record_links
			a href="/archivegrid/contact-information/nnn" title="$here text about archive org"


			a href="/archivegrid/collection/data/samennnnn" <-- ignoring this one for now


*/

func main() {
	// c := colly.NewCollector(colly.AllowedDomains(ALLOWED_DOMAINS[0]))

	// c.OnHTML("div", func(h *colly.HTMLElement) {
	// 	contents := h.ChildAttrs("a", "href")
	// 	fmt.Println(contents)
	// })

	// c.Visit(ARCHIVE_GRID_URL_PATTERNS[0])

	//

	outFileName := "outFile.csv"
	outFile, err := os.Create(outFileName)
	if err != nil {
		log.Fatalf("LOG: ERROR CREATING %s, err %q", outFileName, err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// 	//
	// 	var _AGDomPathsDefinitions_ = AGDomPaths{
	// 		Record: "",
	// 		Record_title:"",
	// 		Record_author: "",
	// 		Record_archive: "" ,
	// 		Record_summary: "",
	// 		Record_links_contact_information: "",

	// }

	type ArchiveGridRecords struct {
		RecId                            int
		Record                           AGRecord
		Record_title                     AGRecordTitle
		Record_author                    AGRecordAuthor
		Record_archive                   AGRecordArchive
		Record_summary                   AGRecordSummary
		Record_links_contact_information AGRecordLinksContactInformation
	}

	//

	c := colly.NewCollector(colly.AllowedDomains(ALLOWED_DOMAINS...))

	c.OnHTML("div.record", func(h *colly.HTMLElement) {

	})

}

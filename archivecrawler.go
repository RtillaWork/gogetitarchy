package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/url"
)

var ALLOWED_DOMAINS []string = []string{"researchworks.oclc.org", "archives.chadwyck.com", "www.newspapers.com"}

var ARCHIVE_GRID_URL_PATTERNS []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%22Albert+Quincy+Porter%22",
}

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

type ArchiveGridRecord struct {
	Id                               HashSum                         `json:"id"`
	MusicianId                       HashSum                         `json:"musician_id"`
	Query                            MusicianQuery                   `json:"musician_query"`
	Record                           AGRecord                        `json:"record"`
	Record_title                     AGRecordTitle                   `json:"record_title"`
	Record_author                    AGRecordAuthor                  `json:"record_author"`
	Record_archive                   AGRecordArchive                 `json:"record_archive"`
	Record_summary                   AGRecordSummary                 `json:"record_summary"`
	Record_links_contact_information AGRecordLinksContactInformation `json:"record_links_contact_information"`
}

func (agr ArchiveGridRecord) String() string {
	return fmt.Sprintf("{ %s | %s | %s }", agr.MusicianId, agr.Query, agr.Record)
}

func (agr ArchiveGridRecord) Hash() HashSum {
	hashfunc := md5.New()
	hashsum := hashfunc.Sum([]byte(agr.String()))
	return HashSum(hashsum)
}

type MusiciansData map[HashSum]ArchiveGridRecord

func ScanArchive(musiciansQueries MusiciansQueries) {

	//
	var AGDomPathsDefinition = AGDomPaths{
		Record:                           "div.record",                // container
		Record_title:                     "div.record_title > h3 > a", // h3>a href ANDTHEN $inner_text
		Record_author:                    "div.record_author",         // span THEN $inner_text
		Record_archive:                   "div.record_archive",        // span THEN $inner_text
		Record_summary:                   "div.record_summary",        // THEN $inner_text
		Record_links_contact_information: "div.record_links",          // a href ANDALSO title
	}

	//

	var ARCHIVE_GRID_BASE_URL = "https://researchworks.oclc.org/archivegrid"
	var AG_BASE_URL, _ = url.Parse(ARCHIVE_GRID_BASE_URL)
	log.Printf("INFO: %v", AG_BASE_URL)

	// type ArchiveGridRecord struct {
	// 	RecId                            int
	// 	Record                           AGRecord
	// 	Record_title                     AGRecordTitle
	// 	Record_author                    AGRecordAuthor
	// 	Record_archive                   AGRecordArchive
	// 	Record_summary                   AGRecordSummary
	// 	Record_links_contact_information AGRecordLinksContactInformation
	// }

	//

	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.MaxDepth(1),
	)

	log.Printf("DEBUG: c.OnHtml\n\n")
	c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {
		record_title := rec.ChildText(AGDomPathsDefinition.Record_title)
		// writer.Write({record_title})
		log.Println(record_title)

	})

	// person_url := fmt.Sprintf(ARCHIVE_GRID_URL_PATTERNS[0], "Albert Quincy Porter")
	musician_url := ARCHIVE_GRID_URL_PATTERNS[0]
	//musician_url := musiciansQueries.first()
	log.Printf("DEBUG %s\n\n", musician_url)

	c.Visit(musician_url)

}

// c := colly.NewCollector(colly.AllowedDomains(ALLOWED_DOMAINS[0]))

// c.OnHTML("div", func(h *colly.HTMLElement) {
// 	contents := h.ChildAttrs("a", "href")
// 	fmt.Println(contents)
// })

// c.Visit(ARCHIVE_GRID_URL_PATTERNS[0])

//////////////////////////////////////

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

//type ArchiveData struct {
//	Person
//	query url.URL
//	map[string]string
//}

////////////////////////////

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

//////////////////////

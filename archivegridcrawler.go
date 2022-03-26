package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var ALLOWED_DOMAINS []string = []string{"researchworks.oclc.org", "archives.chadwyck.com", "www.newspapers.com"}
var ARCHIVE_GRID_URL_PATTERNS []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%22Albert+Quincy+Porter%22",
}

type MusiciansData map[HashSum][]ArchiveGridRecord

func ScanArchiveGridAll(ms MusiciansMap, mqs MusiciansQueries) (musiciansData MusiciansData) {
	const oneSecond = 1_000_000_000 // nanoseconds
	musiciansData = MusiciansData{}

	for mhash, mq := range mqs {
		log.Printf("\nScanArchiveGridAll DEBUG: QUERY %s\n FOR MUSICIAN %s \n\n", mq, ms[mhash].ToCsv())
		musiciansData[mhash] = scanArchiveGrid(ms[mhash], mq)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("press key... ")
		presskey, _ := reader.ReadString('\n')
		fmt.Println(presskey)
		//delay := time.Duration(oneSecond * (rand.Int63n(3*oneSecond) + 1))
		//time.Sleep(delay)
		//log.Printf("DELAY %d", delay)

	}
	return musiciansData

}

func scanArchiveGrid(m Musician, mq MusicianQuery) (agRecords []ArchiveGridRecord) {
	agRecords = []ArchiveGridRecord{}

	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       3 * time.Second,
	})

	c.OnHTML(AGResultsDefinition.Results, func(results *colly.HTMLElement) {
		resultsSize, err := myAtoi(results.ChildText(AGResultsDefinition.ResultsSizeMessage))
		if resultsSize == 0 || err != nil {
			agrecord := NewArchiveGridRecord(m.Id, mq)
			agrecord.IsFound = false
			agRecords = append(agRecords, agrecord)
			return
		} else if resultsSize > 3 {
			// too many to process for now, take note and pass, set IsFound false as flag nor record as non nilfor now
			agrecord := NewArchiveGridRecord(m.Id, mq)
			agrecord.IsFound = true
			agrecord.DebugNotes = AGDEBUG(TOOMANYRECORDS)
			agRecords = append(agRecords, agrecord)
		} else {
			results.ForEach(AGDomPathsDefinition.Record, func(_ int, record *colly.HTMLElement) {
				title := record.ChildText(AGDomPathsDefinition.Record_title)
				author := record.ChildText(AGDomPathsDefinition.Record_author)
				archive := record.ChildText(AGDomPathsDefinition.Record_archive)
				summary := record.ChildText(AGDomPathsDefinition.Record_summary)
				contact := record.ChildAttr(AGDomPathsDefinition.Record_links_contact_information, "title")
				link := record.ChildAttr(AGDomPathsDefinition.Record_links_contact_information, "href")
				log.Printf("\n\n BEGINRECORD:\nTITLE: %s\nAUTHOR: %s\nARCHIVE: %s\nSUMMARY: %s\nCONTACT: %s\nLINK: %s\nENDRECORD\n\n",
					title, author, archive, summary, contact, link)
				agrecord := NewArchiveGridRecord(m.Id, mq)
				agrecord.IsFound = true
				agrecord.set(title, author, archive, summary, contact)

				agrecord.DebugNotes = AGDEBUG(TOOMANYRECORDS)
				agRecords = append(agRecords, agrecord)

			})

			//c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {
			//	record_title := rec.ChildText(AGDomPathsDefinition.Record_title)
			//	// writer.Write({record_title})
			//	log.Println(record_title)
			//
			//})
			//agrecord := NewArchiveGridRecord(m.Id, mq)
			//agrecord.IsFound = true
			//agRecords = append(agRecords, agrecord)
		}
	})

	//c.OnHTML(AGResultsDefinition.ResultsEmpty, func(rec *colly.HTMLElement) {
	//	log.Printf("NOT FOUND")
	//	return
	//})
	//
	//c.OnHTML(AGResultsDefinition.ResultsNotEmpty, func(rec *colly.HTMLElement) {
	//	log.Printf("################## FOUND")
	//	agrecord := NewArchiveGridRecord(m.Id, mq)
	//	agRecords = append(agRecords, agrecord)
	//	return
	//})

	//log.Printf("DEBUG: c.OnHtml\n\n")
	//c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {
	//	record_title := rec.ChildText(AGDomPathsDefinition.Record_title)
	//	// writer.Write({record_title})
	//	log.Println(record_title)
	//
	//})

	// person_url := fmt.Sprintf(ARCHIVE_GRID_URL_PATTERNS[0], "Albert Quincy Porter")
	log.Printf("\n\nscanArchiveGrid DEBUG QUERY %s\n", mq)

	c.Visit(mq.String())

	return agRecords
}

//c.OnHTML(AGResultsDefinition.ResultsSizeMessage, func(e *colly.HTMLElement) {
//	log.Printf("ELEMENT %#v", e.Text) // e.ChildText(AGResultsDefinition.ResultsSize))
//	resultsSize, _ := myAtoi(e.Text)
//	log.Printf("################## FOUND RESULT SIZE: %d", resultsSize)
//	switch {
//	case resultsSize > 5:
//		// too many to process for now, take note and pass, set IsFound false as flag nor record as non nilfor now
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.IsFound = true
//		agrecord.DebugNotes = AGDEBUG(TOOMANYRECORDS)
//		agRecords = append(agRecords, agrecord)
//		break
//	case resultsSize > 0:
//		// crawl each result and add
//
//		c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {
//			record_title := rec.ChildText(AGDomPathsDefinition.Record_title)
//			// writer.Write({record_title})
//			log.Println(record_title)
//
//		})
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.IsFound = true
//		agRecords = append(agRecords, agrecord)
//		break
//	case resultsSize == 0:
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.IsFound = false
//		agRecords = append(agRecords, agrecord)
//		break
//
//	}
//
//	return
//})

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

// Helpers

// like Atoi but cleanes the string out of any non digit characters like comma before the conversion
func myAtoi(s string) (n int, err error) {
	if s == "" {
		return 0, errors.New("myAtoi got empty string probably because no css selector matched OnHtml")
	} else {
		text := strings.Fields(s)
		sint := text[len(text)-1]
		text = strings.Split(sint, ",")
		sint = strings.Join(text, "")
		n, err = strconv.Atoi(sint)
		FailOn(err, "INFO myAtoi EXTRACTING RESULTS SIZE FROM SPAN")
		return n, err
	}
}

// c := colly.NewCollector(colly.AllowedDomains(ALLOWED_DOMAINS[0]))

// c.OnHTML("div", func(h *colly.HTMLElement) {
// 	contents := h.ChildAttrs("a", "href")
// 	fmt.Println(contents)
// })

// c.Visit(ARCHIVE_GRID_URL_PATTERNS[0])

//////////////////////////////////////

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

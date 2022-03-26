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

func CrawlArchiveGrid(ms MusiciansMap, mqs MusiciansQueries, size int) (musiciansData MusiciansData, ok bool) {
	const oneSecond = 1_000_000_000 // nanoseconds
	musiciansData = MusiciansData{}

	if lenms, lenmqs := len(ms), len(mqs); lenms == 0 || lenmqs == 0 || size < 1 {
		return nil, false
	} else {
		log.Printf("Processing %d queries for a MusiciansMap size of %d and a MusiciansQueries size of %d",
			size, lenms, lenmqs)

		for mhash, mq := range mqs {
			if size == 0 {
				break
			}
			log.Printf("\nCrawlArchiveGrid DEBUG: QUERY %s\n FOR MUSICIAN %s \n\n", mq, ms[mhash].ToCsv())
			musiciansData[mhash] = append(musiciansData[mhash], ScanArchiveGrid(ms[mhash], mq)...)

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("press key... ")
			presskey, _ := reader.ReadString('\n')
			fmt.Println(presskey)
			//delay := time.Duration(oneSecond * (rand.Int63n(3*oneSecond) + 1))
			//time.Sleep(delay)
			//log.Printf("DELAY %d", delay)
			size--
		}
		return musiciansData, true
	}

}

func ScanArchiveGrid(m Musician, mq MusicianQuery) (agRecords []ArchiveGridRecord) {
	//agRecord := NewArchiveGridRecord(m.Id, mq)
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

	agrecord := NewArchiveGridRecord(m.Id, mq)

	c.OnHTML(AGDomPathsDefinition.Results, func(results *colly.HTMLElement) {
		resultsSize, err := myAtoi(results.ChildText(AGDomPathsDefinition.ResultsSizeMessage))
		if resultsSize == 0 || err != nil {
			agrecord.ResultCount = 0
			agrecord.IsMatch = true
			agrecord.DebugNotes = AGDEBUG(NORESULTS)
			agRecords = append(agRecords, agrecord)
			return
		} else if resultsSize > 3 {
			// too many to process for now, take note and pass, set ResultCount false as flag nor record as non nilfor now
			agrecord.ResultCount = resultsSize
			agrecord.IsMatch = false
			agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
			agRecords = append(agRecords, agrecord)
		} else {
			agrecord.ResultCount = resultsSize
			//agrecord.set(title, author, archive, summary, contact)
			agrecord.DebugNotes = AGDEBUG(ACCEPTABLERESULTS)
		}
	})


	c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func( rec *colly.HTMLElement) {
//if sanity check
		agrecord := NewArchiveGridRecord(m.Id, mq)
		agrecord.ResultCount = true
		agrecord.set(title, author, archive, summary, contact)


		record := rec.Attr("id")
		title := rec.ChildText(AGDomPathsDefinition.Title)
		author := rec.ChildText(AGDomPathsDefinition.Author)
		archive := rec.ChildText(AGDomPathsDefinition.Archive)
		summary := rec.ChildText(AGDomPathsDefinition.Summary)
		contact := rec.ChildAttr(AGDomPathsDefinition.LinksContactInformation, "title")
		link := rec.ChildAttr(AGDomPathsDefinition.LinksContactInformation, "href")

		log.Printf("\n\n BEGINRECORD: %q\nTITLE: %s\nAUTHOR: %s\nARCHIVE: %s\nSUMMARY: %s\nCONTACT: %s\nLINK: %s\nENDRECORD\n\n",
			record, title, author, archive, summary, contact, link)


		agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
		agRecords = append(agRecords, agrecord)
		)}

	//c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
	//	record_title := rec.ChildText(AGDomPathsDefinition.Title)
	//	// writer.Write({record_title})
	//	log.Println(record_title)
	//
	//})
	//agrecord := NewArchiveGridRecord(m.Id, mq)
	//agrecord.ResultCount = true
	//agRecords = append(agRecords, agrecord)



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
	//c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
	//	record_title := rec.ChildText(AGDomPathsDefinition.Title)
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
//		// too many to process for now, take note and pass, set ResultCount false as flag nor record as non nilfor now
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.ResultCount = true
//		agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
//		agRecords = append(agRecords, agrecord)
//		break
//	case resultsSize > 0:
//		// crawl each result and add
//
//		c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
//			record_title := rec.ChildText(AGDomPathsDefinition.Title)
//			// writer.Write({record_title})
//			log.Println(record_title)
//
//		})
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.ResultCount = true
//		agRecords = append(agRecords, agrecord)
//		break
//	case resultsSize == 0:
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.ResultCount = false
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
		RecordCollectionDataPath: "div.record",                // container
		Title:                    "div.record_title > h3 > a", // h3>a href ANDTHEN $inner_text
		Author:                   "div.record_author",         // span THEN $inner_text
		Archive:                  "div.record_archive",        // span THEN $inner_text
		Summary:                  "div.record_summary",        // THEN $inner_text
		LinksContactInformation:  "div.record_links",          // a href ANDALSO title
	}

	//

	var ARCHIVE_GRID_BASE_URL = "https://researchworks.oclc.org/archivegrid"
	var AG_BASE_URL, _ = url.Parse(ARCHIVE_GRID_BASE_URL)
	log.Printf("INFO: %v", AG_BASE_URL)

	// type ArchiveGridRecord struct {
	// 	RecId                            int
	// 	RecordCollectionDataPath                           AGRecord
	// 	Title                     AGRecordTitle
	// 	Author                    AGRecordAuthor
	// 	Archive                   AGRecordArchive
	// 	Summary                   AGRecordSummary
	// 	LinksContactInformation AGRecordLinksContactInformation
	// }

	//

	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.MaxDepth(1),
	)

	log.Printf("DEBUG: c.OnHtml\n\n")
	c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
		record_title := rec.ChildText(AGDomPathsDefinition.Title)
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

//
//func scanArchiveGrid(m Musician, mq MusicianQuery) (agRecords []ArchiveGridRecord) {
//	agRecords = []ArchiveGridRecord{}
//
//	c := colly.NewCollector(
//		colly.AllowedDomains(ALLOWED_DOMAINS...),
//		colly.MaxDepth(1),
//	)
//
//	c.Limit(&colly.LimitRule{
//		DomainGlob:  "*",
//		Parallelism: 1,
//		Delay:       3 * time.Second,
//	})
//
//	c.OnHTML(AGResultsDefinition.Results, func(results *colly.HTMLElement) {
//		resultsSize, err := myAtoi(results.ChildText(AGResultsDefinition.ResultsSizeMessage))
//		if resultsSize == 0 || err != nil {
//			agrecord := NewArchiveGridRecord(m.Id, mq)
//			agrecord.ResultCount = false
//			agRecords = append(agRecords, agrecord)
//			return
//		} else if resultsSize > 3 {
//			// too many to process for now, take note and pass, set ResultCount false as flag nor record as non nilfor now
//			agrecord := NewArchiveGridRecord(m.Id, mq)
//			agrecord.ResultCount = true
//			agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
//			agRecords = append(agRecords, agrecord)
//		} else {
//			results.ForEach(AGDomPathsDefinition.RecordCollectionDataPath, func(_ int, record *colly.HTMLElement) {
//				title := record.ChildText(AGDomPathsDefinition.Title)
//				author := record.ChildText(AGDomPathsDefinition.Author)
//				archive := record.ChildText(AGDomPathsDefinition.Archive)
//				summary := record.ChildText(AGDomPathsDefinition.Summary)
//				contact := record.ChildAttr(AGDomPathsDefinition.LinksContactInformation, "title")
//				link := record.ChildAttr(AGDomPathsDefinition.LinksContactInformation, "href")
//				log.Printf("\n\n BEGINRECORD:\nTITLE: %s\nAUTHOR: %s\nARCHIVE: %s\nSUMMARY: %s\nCONTACT: %s\nLINK: %s\nENDRECORD\n\n",
//					title, author, archive, summary, contact, link)
//				agrecord := NewArchiveGridRecord(m.Id, mq)
//				agrecord.ResultCount = true
//				agrecord.set(title, author, archive, summary, contact)
//
//				agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
//				agRecords = append(agRecords, agrecord)
//
//			})
//
//			//c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
//			//	record_title := rec.ChildText(AGDomPathsDefinition.Title)
//			//	// writer.Write({record_title})
//			//	log.Println(record_title)
//			//
//			//})
//			//agrecord := NewArchiveGridRecord(m.Id, mq)
//			//agrecord.ResultCount = true
//			//agRecords = append(agRecords, agrecord)
//		}
//	})
//
//	//c.OnHTML(AGResultsDefinition.ResultsEmpty, func(rec *colly.HTMLElement) {
//	//	log.Printf("NOT FOUND")
//	//	return
//	//})
//	//
//	//c.OnHTML(AGResultsDefinition.ResultsNotEmpty, func(rec *colly.HTMLElement) {
//	//	log.Printf("################## FOUND")
//	//	agrecord := NewArchiveGridRecord(m.Id, mq)
//	//	agRecords = append(agRecords, agrecord)
//	//	return
//	//})
//
//	//log.Printf("DEBUG: c.OnHtml\n\n")
//	//c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
//	//	record_title := rec.ChildText(AGDomPathsDefinition.Title)
//	//	// writer.Write({record_title})
//	//	log.Println(record_title)
//	//
//	//})
//
//	// person_url := fmt.Sprintf(ARCHIVE_GRID_URL_PATTERNS[0], "Albert Quincy Porter")
//	log.Printf("\n\nscanArchiveGrid DEBUG QUERY %s\n", mq)
//
//	c.Visit(mq.String())
//
//	return agRecords
//}

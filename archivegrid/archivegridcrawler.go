package archivegrid

import (
	"errors"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils"

	//"github.com/RtillaWork/gogetitarchy"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const TOOMANYRESULTSVALUE int = 15000

var ALLOWED_DOMAINS []string = []string{"researchworks.oclc.org", "archives.chadwyck.com", "www.newspapers.com"}
var ARCHIVE_GRID_URL_PATTERNS []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%22Albert+Quincy+Porter%22",
}

type MusiciansData map[utils.HashSum][]*ArchiveGridRecord

func CrawlArchiveGrid(ms musician.MusiciansMap, mqs MusiciansQueries, size int, phrases []string) (musiciansData MusiciansData, ok bool) {
	const oneSecond = 1_000_000_000 // nanoseconds
	musiciansData = MusiciansData{}

	if lenms, lenmqs := len(ms), len(mqs); lenms == 0 || lenmqs == 0 || size < 1 {
		return nil, false
		log.Println("CrawlArchiveGrid Parameter(s) Error. Returned prematurely")
	} else {
		log.Printf("Processing %d queries for a MusiciansMap size of %d and a MusiciansQueries size of %d",
			size, lenms, lenmqs)

		for mhash, mq := range mqs {
			if size == 0 {
				break
			}
			log.Printf("\nCrawlArchiveGrid DEBUG: QUERY %s\n  \n\n", mq)
			musiciansData[mhash] = append(musiciansData[mhash], ScanArchiveGrid(ms[mhash], mq, phrases)...)

			utils.WaitForKeypress()
			//delay := time.Duration(oneSecond * (rand.Int63n(3*oneSecond) + 1))
			//time.Sleep(delay)
			//log.Printf("DELAY %d", delay)
			size--
		}

	}
	return musiciansData, true
}

func ScanArchiveGrid(m *musician.Musician, mq *MusicianQuery, phrases []string) (agRecords []*ArchiveGridRecord) {
	//agRecord := NewArchiveGridRecord(m.Id, mq)
	agRecords = []*ArchiveGridRecord{}

	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       7 * time.Second,
	})

	c.OnRequest(func(request *colly.Request) {
		log.Printf("\nSTARTING NEW REQUEST TIME: %v\n", time.Now())
		log.Printf("\nscanArchiveGrid YOUR QUERY %s\n", mq)
		log.Printf("\nscanArchiveGrid REQUEST: %s\n\n", request.URL)

	})

	c.OnError(func(response *colly.Response, err error) {
		log.Printf("\nERROR: Time %v \tRESPONSE STATUS: %v\nERROR: %s", time.Now(), response.StatusCode, err)

	})

	c.OnHTML(AGDomPathsDefinition.Results, func(results *colly.HTMLElement) {
		resultsSize, err := myAtoi(results.ChildText(AGDomPathsDefinition.ResultsSizeMessage))
		if resultsSize == 0 || err != nil {
			mq.SetResultCount(0)
			mq.DebugNotes = QUERYDEBUG(NORESULTS)

			log.Printf("RESULT SIZE resultsSize == 0 || err != nil %d", resultsSize)
			//agrecord.ResultCount = -1
			//agrecord.IsMatch = true
			//agrecord.DebugNotes = AGDEBUG(NORESULTS)
			//agRecords = append(agRecords, agrecord)
			return
		} else if resultsSize > TOOMANYRESULTSVALUE {
			log.Printf("RESULT SIZE resultsSize > TOOMANYRESULTSVALUE %d", resultsSize)
			// too many to process for now, take note and pass, set ResultSize false as flag nor record as non nilfor now
			mq.SetResultCount(0)
			mq.DebugNotes = QUERYDEBUG(TOOMANYRESULTS)
			//agrecord.ResultCount = resultsSize
			//agrecord.IsMatch = false
			//agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
			//agRecords = append(agRecords, agrecord)
		} else {
			log.Printf("RESULT SIZE ok supposed to process the other OnHtml for AG DOM elements %d", resultsSize)
			mq.SetResultCount(resultsSize)
			mq.DebugNotes = QUERYDEBUG(ACCEPTABLERESULTS)
			//agrecord.ResultCount = resultsSize
			//agrecord.set(title, author, archive, summary, contact)
			//agrecord.DebugNotes = AGDEBUG(ACCEPTABLERESULTS)
		}
	})

	c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {

		// exit this OnHtml if there is nothing to search for or sanity doesn't check
		if mq.ResultSize < 1 {
			return
		}
		agrecord := NewArchiveGridRecord(m.Id, *mq)
		record := rec.ChildAttr(AGDomPathsDefinition.RecordCollectionDataPath, "value")
		title := rec.ChildText(AGDomPathsDefinition.Title)
		//title := rec.DOM.Find(AGDomPathsDefinition.Title).Text()
		//.ChildText(AGDomPathsDefinition.Title)
		author := rec.ChildText(AGDomPathsDefinition.Author)
		//.DOM.Find(AGDomPathsDefinition.Author).Text()
		archive := rec.ChildText(AGDomPathsDefinition.Archive)
		//.DOM.Find(AGDomPathsDefinition.Archive).Text()
		summary := rec.ChildText(AGDomPathsDefinition.Summary)
		//.DOM.Find(AGDomPathsDefinition.Summary).Text()
		link := rec.ChildAttr(AGDomPathsDefinition.LinksContactInformation, "href")
		//.DOM.Find(AGDomPathsDefinition.LinksContactInformation).Attr("href")
		//ChildAttr(AGDomPathsDefinition.LinksContactInformation, "href")
		contact := rec.ChildAttr(AGDomPathsDefinition.ContactInformation, "title")
		//.DOM.Find(AGDomPathsDefinition.ContactInformation).Attr("title")
		//rec.ChildAttr(AGDomPathsDefinition.ContactInformation, "title")

		agrecord.Set(record, title, author, archive, summary, link, contact)
		//log.Printf("\n\n RECORDOBJECT: \nBEGINRECORD: %#v\nTITLE: %#v\nAUTHOR: %#v\nARCHIVE: %#v\nSUMMARY: %#v\nCONTACT: %#v\nLINK: %#v\nENDRECORD\n\n",
		//	record, title, author, archive, summary, contact, link)

		if matches := agrecord.ContainsAnyFolded(phrases); matches > 0 || phrases == nil {
			log.Printf("\n\n RECORDOBJECT: \nBEGINRECORD: %#v\nTITLE: %#v\nAUTHOR: %#v\nARCHIVE: %#v\nSUMMARY: %#v\nCONTACT: %#v\nLINK: %#v\nENDRECORD\n\n",
				record, title, author, archive, summary, contact, link)
			log.Printf("FILTERED IN matches = %d", matches)
			agrecord.DebugNotes = AGDEBUG(FOUNDANDVALIDATED)
			agrecord.IsMatch = true
			mq.Matches = matches
			agRecords = append(agRecords, agrecord)
		} else {
			log.Printf("FILTERED OUT matches = %d", matches)
		}

	})

	//c.OnHTML(AGDomPathsDefinition.RecordCollectionDataPath, func(rec *colly.HTMLElement) {
	//	record_title := rec.ChildText(AGDomPathsDefinition.Title)
	//	// writer.Write({record_title})
	//	log.Println(record_title)
	//
	//})
	//agrecord := NewArchiveGridRecord(m.Id, mq)
	//agrecord.ResultSize = true
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

	c.Visit(mq.String())

	return agRecords
}

//c.OnHTML(AGResultsDefinition.ResultsSizeMessage, func(e *colly.HTMLElement) {
//	log.Printf("ELEMENT %#v", e.Text) // e.ChildText(AGResultsDefinition.ResultsSize))
//	resultsSize, _ := myAtoi(e.Text)
//	log.Printf("################## FOUND RESULT SIZE: %d", resultsSize)
//	switch {
//	case resultsSize > 5:
//		// too many to process for now, take note and pass, set ResultSize false as flag nor record as non nilfor now
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.ResultSize = true
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
//		agrecord.ResultSize = true
//		agRecords = append(agRecords, agrecord)
//		break
//	case resultsSize == 0:
//		agrecord := NewArchiveGridRecord(m.Id, mq)
//		agrecord.ResultSize = false
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
		log.Printf("\nTEXT: %v\n", text)
		sint := text[len(text)-1]
		log.Printf("SINT: %v\n", sint)
		sint = sint[:len(sint)-1]
		log.Printf("SINT: %v\n", sint)
		text = strings.Split(sint, ",")
		log.Printf("\nTEXT: %v\n", text)
		n, err = strconv.Atoi(strings.Join(text, ""))
		sint, text = "", nil
		utils.FailOn(err, "INFO myAtoi EXTRACTING RESULTS SIZE FROM SPAN")
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
//			agrecord.ResultSize = false
//			agRecords = append(agRecords, agrecord)
//			return
//		} else if resultsSize > 3 {
//			// too many to process for now, take note and pass, set ResultSize false as flag nor record as non nilfor now
//			agrecord := NewArchiveGridRecord(m.Id, mq)
//			agrecord.ResultSize = true
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
//				agrecord.ResultSize = true
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
//			//agrecord.ResultSize = true
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

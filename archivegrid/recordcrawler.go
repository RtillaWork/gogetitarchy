package archivegrid

import (
	"errors"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils"
	errors2 "github.com/RtillaWork/gogetitarchy/utils/errors"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

func CrawlArchiveGrid(ms musician.MusiciansMap, mqs MusiciansQueries, size int, phrases []string) (musiciansData MusiciansData, ok bool) {
	const oneSecond = 1_000_000_000 // nanoseconds
	musiciansData = MusiciansData{}

	if lenms, lenmqs := len(ms), len(mqs); lenms == 0 || lenmqs == 0 || size < 1 {
		return nil, false
		log.Println("CrawlArchiveGrid Parameter(s) Error. Returned prematurely")
	} else {
		log.Printf("Processing %d queries for a MusiciansMap size of %d and a MusiciansQueries size of %d",
			size, lenms, lenmqs)
		utils.WaitForKeypress()

		for mhash, mq := range mqs {
			if size == 0 {
				break
			}
			mq.SetResultCountFunc(ScanQueryResultSize)

			if mq.ResultSize > 0 {
				log.Printf("\nCrawlArchiveGrid DEBUG: QUERY result count %d\n For query %s \n QUERYINg...\n", mq.ResultSize, mq.String())
				recordsresponse, err := ScanArchiveGrid(ms[mhash], mq, phrases)
				errors2.FailOn(err, "Crawling query\n"+mq.String())
				musiciansData[mhash] = append(musiciansData[mhash], recordsresponse...)
			} else {
				log.Printf("\nCrawlArchiveGrid DEBUG: QUERY result count <= 0 %d\n For query %s \n SKIPPING...\n", mq.ResultSize, mq)

			}

			utils.WaitForKeypress()
			//delay := time.Duration(oneSecond * (rand.Int63n(3*oneSecond) + 1))
			//time.Sleep(delay)
			//log.Printf("DELAY %d", delay)
			size--
		}

	}
	return musiciansData, true
}

// Get query's specific parameters particularily result size
func ScanQueryResultSize(mq MusicianQuery) (resultsize int, err error) {
	resultsize = -1
	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(request *colly.Request) {
		log.Printf("\n ScanQueryResultSize STARTING NEW REQUEST TIME: %v\n", time.Now())
		log.Printf("\n ScanQueryResultSizeQueryResultSize YOUR QUERY %s\n", mq.String())
		log.Printf("\nScanQueryResultSize scanArchiveGrid REQUEST: %s\n\n", request.URL)

	})

	c.OnError(func(response *colly.Response, e error) {
		err = e
		log.Printf("\nScanQueryResultSize ERROR: Time %v \tRESPONSE STATUS: %v\nERROR: %s", time.Now(), response.StatusCode, err)

	})

	c.OnResponse(func(response *colly.Response) {
		log.Printf("\nScanQueryResultSize Received response about request %s\n", response.Request.URL)
	})

	c.OnHTML(AGDomPathsDefinition.ResultsSize, func(elem *colly.HTMLElement) {
		//resultsizehtml, e := results.DOM.Find("#resultsize").Html()
		resultsizehtml, e := elem.DOM.Html()
		if e != nil {
			err = e
			log.Printf("ERROR ResultSize elem.DOM.Html() ERROR %s, \nresultsizeHtml: %s\n", err, resultsizehtml)
			return
		}
		log.Printf("ResultSizeelem.DOM.Html() %s", resultsizehtml)
		utils.WaitForKeypress()

		//if e != nil {
		//	err = e
		//} else {
		resultsize, err = totalPagesAtoi(resultsizehtml)
		//}
	})
	c.Visit(mq.String())

	return resultsize, err

}

func ScanArchiveGrid(m *musician.Musician, mq *MusicianQuery, phrases []string) (agRecords []*Record, err error) {
	//agRecord := NewArchiveGridRecord(m.Id, mq)
	agRecords = []*Record{}

	c := colly.NewCollector(
		colly.AllowedDomains(ALLOWED_DOMAINS...),
		colly.MaxDepth(1),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       1 * time.Second,
		// Add an additional random delay
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(request *colly.Request) {
		log.Printf("\nSTARTING NEW REQUEST TIME: %v\n", time.Now())
		log.Printf("\nscanArchiveGrid YOUR QUERY %s\n", mq)
		log.Printf("\nscanArchiveGrid REQUEST: %s\n\n", request.URL)

	})

	c.OnError(func(response *colly.Response, e error) {
		err = e
		log.Printf("\nERROR: Time %v \tRESPONSE STATUS: %v\nERROR: %s", time.Now(), response.StatusCode, err)

	})

	c.OnResponse(func(response *colly.Response) {
		log.Printf("\nReceived response about request %s\n", response.Request.URL)
	})

	c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {

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

	c.Visit(mq.String())

	return agRecords, err
}

// Helpers

// like Atoi but cleanes the string out of any non digit characters like comma before the conversion
func totalPagesAtoi(s string) (n int, err error) {
	rexTotalPages := regexp.MustCompile(`\d+`)
	sint := rexTotalPages.FindString(s)
	if s == "" || sint == "" {
		return -1, errors.New("totalPagesAtoi got empty string probably because no css selector matched OnHtml")
	} else {
		n, err = strconv.Atoi(sint)
		errors2.FailOn(err, "INFO totalPagesAtoi EXTRACTING RESULTS SIZE FROM SPAN")

		sint = ""
		return n, err
	}
}

func FilteredMusiciansDataBuilder(m *musician.Musician, mq *MusicianQuery, phrases []string) (agRecords []*Record) {
	return nil
}

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

////////////

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

	// type Record struct {
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

// OLD
//// like Atoi but cleanes the string out of any non digit characters like comma before the conversion
//func totalPagesAtoi(s string) (n int, err error) {
//	if s == "" {
//		return 0, errors.New("totalPagesAtoi got empty string probably because no css selector matched OnHtml")
//	} else {
//		text := strings.Fields(s)
//		log.Printf("\nTEXT: %v\n", text)
//		sint := text[len(text)-1]
//		log.Printf("SINT: %v\n", sint)
//		sint = sint[:len(sint)-1]
//		log.Printf("SINT: %v\n", sint)
//		text = strings.Split(sint, ",")
//		log.Printf("\nTEXT: %v\n", text)
//		n, err = strconv.Atoi(strings.Join(text, ""))
//		sint, text = "", nil
//		errors2.FailOn(err, "INFO totalPagesAtoi EXTRACTING RESULTS SIZE FROM SPAN")
//		return n, err
//	}
//}

// OLD works but result size incorrect
//func ScanArchiveGrid(m *musician.Musician, mq *MusicianQuery, phrases []string) (agRecords []*Record) {
//	//agRecord := NewArchiveGridRecord(m.Id, mq)
//	agRecords = []*Record{}
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
//	c.OnRequest(func(request *colly.Request) {
//		log.Printf("\nSTARTING NEW REQUEST TIME: %v\n", time.Now())
//		log.Printf("\nscanArchiveGrid YOUR QUERY %s\n", mq)
//		log.Printf("\nscanArchiveGrid REQUEST: %s\n\n", request.URL)
//
//	})
//
//	c.OnError(func(response *colly.Response, err error) {
//		log.Printf("\nERROR: Time %v \tRESPONSE STATUS: %v\nERROR: %s", time.Now(), response.StatusCode, err)
//
//	})
//
//	c.OnResponse(func(response *colly.Response) {
//		log.Printf("\nReceived response about request %s\n", response.Request.URL)
//	})
//
//	c.OnHTML(AGDomPathsDefinition.Results, func(results *colly.HTMLElement) {
//		resultsSize, err := totalPagesAtoi(results.ChildText(AGDomPathsDefinition.ResultsSizeMessage))
//		if resultsSize == 0 || err != nil {
//			mq.SetResultCount(0)
//			mq.DebugNotes = QUERYDEBUG(NORESULTS)
//
//			log.Printf("RESULT SIZE resultsSize == 0 || err != nil %d", resultsSize)
//			//agrecord.MatchestCount = -1
//			//agrecord.IsMatch = true
//			//agrecord.DebugNotes = AGDEBUG(NORESULTS)
//			//agRecords = append(agRecords, agrecord)
//			return
//		} else if resultsSize > TOOMANYRESULTSVALUE {
//			log.Printf("RESULT SIZE resultsSize > TOOMANYRESULTSVALUE %d", resultsSize)
//			// too many to process for now, take note and pass, set ResultSize false as flag nor record as non nilfor now
//			mq.SetResultCount(0)
//			mq.DebugNotes = QUERYDEBUG(TOOMANYRESULTS)
//			//agrecord.MatchestCount = resultsSize
//			//agrecord.IsMatch = false
//			//agrecord.DebugNotes = AGDEBUG(TOOMANYRESULTS)
//			//agRecords = append(agRecords, agrecord)
//		} else {
//			log.Printf("RESULT SIZE ok supposed to process the other OnHtml for AG DOM elements %d", resultsSize)
//			mq.SetResultCount(resultsSize)
//			mq.DebugNotes = QUERYDEBUG(ACCEPTABLERESULTS)
//			//agrecord.MatchestCount = resultsSize
//			//agrecord.set(title, author, archive, summary, contact)
//			//agrecord.DebugNotes = AGDEBUG(ACCEPTABLERESULTS)
//		}
//	})
//
//	c.OnHTML(AGDomPathsDefinition.Record, func(rec *colly.HTMLElement) {
//
//		// exit this OnHtml if there is nothing to search for or sanity doesn't check
//		switch {
//		case mq.ResultSize < 1:
//			return
//		case mq.ResultSize > QUERY_LIMIT:
//			//
//		}
//
//		agrecord := NewArchiveGridRecord(m.Id, *mq)
//		record := rec.ChildAttr(AGDomPathsDefinition.RecordCollectionDataPath, "value")
//		title := rec.ChildText(AGDomPathsDefinition.Title)
//		//title := rec.DOM.Find(AGDomPathsDefinition.Title).Text()
//		//.ChildText(AGDomPathsDefinition.Title)
//		author := rec.ChildText(AGDomPathsDefinition.Author)
//		//.DOM.Find(AGDomPathsDefinition.Author).Text()
//		archive := rec.ChildText(AGDomPathsDefinition.Archive)
//		//.DOM.Find(AGDomPathsDefinition.Archive).Text()
//		summary := rec.ChildText(AGDomPathsDefinition.Summary)
//		//.DOM.Find(AGDomPathsDefinition.Summary).Text()
//		link := rec.ChildAttr(AGDomPathsDefinition.LinksContactInformation, "href")
//		//.DOM.Find(AGDomPathsDefinition.LinksContactInformation).Attr("href")
//		//ChildAttr(AGDomPathsDefinition.LinksContactInformation, "href")
//		contact := rec.ChildAttr(AGDomPathsDefinition.ContactInformation, "title")
//		//.DOM.Find(AGDomPathsDefinition.ContactInformation).Attr("title")
//		//rec.ChildAttr(AGDomPathsDefinition.ContactInformation, "title")
//
//		agrecord.Set(record, title, author, archive, summary, link, contact)
//		//log.Printf("\n\n RECORDOBJECT: \nBEGINRECORD: %#v\nTITLE: %#v\nAUTHOR: %#v\nARCHIVE: %#v\nSUMMARY: %#v\nCONTACT: %#v\nLINK: %#v\nENDRECORD\n\n",
//		//	record, title, author, archive, summary, contact, link)
//
//		if matches := agrecord.ContainsAnyFolded(phrases); matches > 0 || phrases == nil {
//			log.Printf("\n\n RECORDOBJECT: \nBEGINRECORD: %#v\nTITLE: %#v\nAUTHOR: %#v\nARCHIVE: %#v\nSUMMARY: %#v\nCONTACT: %#v\nLINK: %#v\nENDRECORD\n\n",
//				record, title, author, archive, summary, contact, link)
//			log.Printf("FILTERED IN matches = %d", matches)
//			agrecord.DebugNotes = AGDEBUG(FOUNDANDVALIDATED)
//			agrecord.IsMatch = true
//			mq.Matches = matches
//			agRecords = append(agRecords, agrecord)
//		} else {
//			log.Printf("FILTERED OUT matches = %d", matches)
//		}
//
//	})
//
//	c.Visit(mq.String())
//
//	return agRecords
//}

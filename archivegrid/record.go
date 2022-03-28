package archivegrid

import (
	"crypto/md5"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils"

	//"github.com/RtillaWork/gogetitarchy"
	"io"
	"strings"
)

type AGDEBUG int

const (
	EMPTY AGDEBUG = iota
	INPROGRESS
	FOUNDNOTVALIDATEDYET
	FOUNDANDNOTVALIDATED
	FOUNDANDVALIDATED
	//TOOMANYRESULTS
	//NORESULTS
	//ACCEPTABLERESULTS
)

type AGOrganization struct {
	orgId               int
	name                string
	contact_information string
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

const ArchiveGridRecordSTRINGNULL = "NODATAFOUND"

//type Record struct {
//	Id                               HashSum                         `json:"id"`
//	MusicianId                       HashSum                         `json:"musician_id"`
//	Query                            MusicianQuery                   `json:"musician_query"`
//	ResultSize                            bool                            `json:"is_found"`
//	RecordCollectionDataPath                           AGRecord                        `json:"record"`
//	Title                     AGRecordTitle                   `json:"record_title"`
//	Author                    AGRecordAuthor                  `json:"record_author"`
//	Archive                   AGRecordArchive                 `json:"record_archive"`
//	Summary                   AGRecordSummary                 `json:"record_summary"`
//	LinksContactInformation AGRecordLinksContactInformation `json:"record_links_contact_information"`
//	DebugNotes                       AGDEBUG                         `json:"debug_notes"`
//}

type Record struct {
	Id                       utils.HashSum `json:"id"`
	MusicianId               utils.HashSum `json:"musician_id"`
	Query                    MusicianQuery `json:"musician_query"`
	ResultCount              int           `json:"result_count"`
	IsMatch                  bool          `json:"is_match"`
	RecordCollectionDataPath string        `json:"record_collection_datapath"`
	Title                    string        `json:"record_title"`
	Author                   string        `json:"record_author"`
	Archive                  string        `json:"record_archive"`
	Summary                  string        `json:"record_summary"`
	LinksContactInformation  string        `json:"links_contact_information"`
	ContactInformation       string        `json:"contact_information"`
	DebugNotes               AGDEBUG       `json:"debug_notes"`
}

func (agr *Record) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%s%s", agr.MusicianId, agr.Query)
}

//func (agr Record) String() string {
//	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Archive.href)
//}

func (agr *Record) String() string {
	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Archive)
}

func (agr *Record) ToJson() string {
	return fmt.Sprintf("{\"ag_record_id\": %q, \n\"musician_id\": %q, \n\"query\": %q, \n}",
		agr.Id, agr.MusicianId, agr.Query)
}

//func (agr Record) ToCsv() string {
//	return fmt.Sprintf("%s; %s; %s; %s", agr.Id, agr.MusicianId, agr.Query, agr.Archive.href)
//}

func (agr *Record) ToCsv() string {
	return fmt.Sprintf("%q; %q; %q; %d; %q; %q; %q; %q; %q; %q; %q; %q\n",
		agr.Id,
		agr.MusicianId,
		agr.Query,
		agr.ResultCount,
		agr.RecordCollectionDataPath,
		agr.Title,
		agr.Author,
		agr.Archive,
		agr.Summary,
		agr.LinksContactInformation,
		agr.ContactInformation,
		agr.DebugNotes)
}

func (agr Record) Hash() utils.HashSum {
	hashfunc := md5.New()
	data := agr.PrimaryKey()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return utils.HashSum(fmt.Sprintf("%x", hashsum))
}

func NewArchiveGridRecord(musicianId utils.HashSum, query MusicianQuery) (archiveGridRecord *Record) {
	archiveGridRecord = new(Record)
	archiveGridRecord = &Record{
		MusicianId:  musicianId,
		Query:       query,
		ResultCount: -1,
		//RecordCollectionDataPath:                           ArchiveGridRecordSTRINGNULL,
		//Title:                     AGRecordTitle,
		//Author:                    AGRecordAuthor,
		//Archive:                   AGRecordArchive,
		//Summary:                   AGRecordSummary,
		//LinksContactInformation: AGRecordLinksContactInformation,
	}

	archiveGridRecord.Id = archiveGridRecord.Hash()
	return archiveGridRecord
}

func (agr *Record) Destroy() {
	agr.Id = ""
	agr.MusicianId = ""
	agr.Query = MusicianQuery{}
	agr.ResultCount = 0
	agr.IsMatch = false
	agr.RecordCollectionDataPath = ""
	agr.Title = ""
	agr.Author = ""
	agr.Archive = ""
	agr.Summary = ""
	agr.LinksContactInformation = ""
	agr.ContactInformation = ""
	agr.DebugNotes = AGDEBUG(0)
	return
}

func (agr *Record) Set(record, title, author, archive, summary, link, contact string) {
	agr.IsMatch = false
	agr.RecordCollectionDataPath = record
	agr.Title = title
	agr.Author = author
	agr.Archive = archive
	agr.Summary = summary
	agr.LinksContactInformation = link
	agr.ContactInformation = contact
	agr.DebugNotes = AGDEBUG(FOUNDNOTVALIDATEDYET)
}

func (agr *Record) ContainsAnyFolded(phrases []string) (matches int) {
	if len(phrases) < 1 {
		return -1
	}

	for _, phrase := range phrases {

		p := strings.ToLower(phrase)
		//log.Printf("A PHRASE %s", p)
		//WaitForKeypress()
		if strings.Contains(strings.ToLower(agr.Title), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(agr.Author), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(agr.Archive), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(agr.Summary), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(agr.ContactInformation), p) {
			matches++
			//WaitForKeypress()

		}
	}
	return matches
}

//

type AGDomPaths struct {
	Record                   string
	RecordCollectionDataPath string // AGRecord.Dom
	Title                    string // AGRecordTitle.Dom
	Author                   string // AGRecordAuthor.Dom
	Archive                  string // AGRecordArchive.Dom
	Summary                  string // AGRecordSummary.Dom
	LinksContactInformation  string // AGRecordLinksContactInformation.Dom
	ContactInformation       string
	Results                  string
	ResultsNotEmpty          string
	ResultsEmpty             string
	ResultsSize              string
	ResultsSizeMessage       string
	ResultsNext              string
}

var AGDomPathsDefinition = AGDomPaths{
	Record:                   "div.record",
	RecordCollectionDataPath: "input[value]",                      // container->archivegrid collection data path  "div.record > input[value]",
	Title:                    "div.record_title > h3 > a[title]",  // h3>a href THEN $inner_text "div.record_title > h3 > a[title]"
	Author:                   "div.record_author span[itemprop]",  // span[itemprop="name"] THEN $inner_text "div.record_author span[itemprop]"
	Archive:                  "div.record_archive span[itemprop]", // span[itemprop="name"] THEN $inner_text  "div.record_archive span[itemprop]"
	Summary:                  "div.record_summary",                // THEN $inner_text
	LinksContactInformation:  "div.record_links > a[href]",        // a href ANDALSO    "div.record_links > a[href]"
	ContactInformation:       "div.record_links > a[title]",       // a  ANDALSO title   "div.record_links > a[title]"
	Results:                  "div.results",
	ResultsNotEmpty:          "div.results div.searchresult",
	ResultsEmpty:             "div.results div.alertresult",
	ResultsSize:              "resultsize", // "main h2 > span[id=resultsize]", // "main > h2", // "main h2 > span#resultsize"
	ResultsSizeMessage:       "div.navtable > div.navrow > div.navrowright > span",
	ResultsNext:              ".results .navtable .navrow a[title=\"View the Next page of results\"]", // get the href

}

//type AGResults struct {
//	Results            string
//	ResultsNotEmpty    string //div.results > div.searchresults
//	ResultsEmpty       string // div.results > div.alertresult
//	ResultsSize        string // span#resultsize
//	ResultsSizeMessage string
//	ResultsNext        string
//}

//var AGResultsDefinition = AGResults{
//	Results:            "div.results",
//	ResultsNotEmpty:    "div.results > div.searchresult",
//	ResultsEmpty:       "div.results > div.alertresult",
//	ResultsSize:        "main > h2", // "main h2 > span#resultsize"
//	ResultsSizeMessage: ".navrow span",
//	ResultsNext:        ".results .navtable .navrow a[title=\"View the Next page of results\"]", // get the href
//}

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

/*

main
	...
	span#resultsize
		$text

div.results
	div.alertresult
	div
		text " No ArchiveGrid collection descriptions match this search:"

div.results
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

//
////
//var AGDomPathsDefinition = AGDomPaths{
//	RecordCollectionDataPath:                           "div.record",                // container
//	Title:                     "div.record_title > h3 > a", // h3>a href ANDTHEN $inner_text
//	Author:                    "div.record_author",         // span THEN $inner_text
//	Archive:                   "div.record_archive",        // span THEN $inner_text
//	Summary:                   "div.record_summary",        // THEN $inner_text
//	LinksContactInformation: "div.record_links",          // a href ANDALSO title
//}
//
////
//
//var ARCHIVE_GRID_BASE_URL = "https://researchworks.oclc.org/archivegrid"
//var AG_BASE_URL, _ = url.Parse(ARCHIVE_GRID_BASE_URL)
//log.Printf("INFO: %v", AG_BASE_URL)
//
//// type Record struct {
//// 	RecId                            int
//// 	RecordCollectionDataPath                           AGRecord
//// 	Title                     AGRecordTitle
//// 	Author                    AGRecordAuthor
//// 	Archive                   AGRecordArchive
//// 	Summary                   AGRecordSummary
//// 	LinksContactInformation AGRecordLinksContactInformation
//// }
//
////

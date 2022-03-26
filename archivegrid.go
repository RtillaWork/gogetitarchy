package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

type AGDEBUG int

const (
	EMPTY AGDEBUG = iota
	TOOMANYRECORDS
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

//type ArchiveGridRecord struct {
//	Id                               HashSum                         `json:"id"`
//	MusicianId                       HashSum                         `json:"musician_id"`
//	Query                            MusicianQuery                   `json:"musician_query"`
//	IsFound                            bool                            `json:"is_found"`
//	Record                           AGRecord                        `json:"record"`
//	Record_title                     AGRecordTitle                   `json:"record_title"`
//	Record_author                    AGRecordAuthor                  `json:"record_author"`
//	Record_archive                   AGRecordArchive                 `json:"record_archive"`
//	Record_summary                   AGRecordSummary                 `json:"record_summary"`
//	Record_links_contact_information AGRecordLinksContactInformation `json:"record_links_contact_information"`
//	DebugNotes                       AGDEBUG                         `json:"debug_notes"`
//}

type ArchiveGridRecord struct {
	Id                               HashSum       `json:"id"`
	MusicianId                       HashSum       `json:"musician_id"`
	Query                            MusicianQuery `json:"musician_query"`
	IsFound                          bool          `json:"is_found"`
	IsMatch                          bool          `json:"is_match"`
	Record                           string        `json:"record"`
	Record_title                     string        `json:"record_title"`
	Record_author                    string        `json:"record_author"`
	Record_archive                   string        `json:"record_archive"`
	Record_summary                   string        `json:"record_summary"`
	Record_links_contact_information string        `json:"record_links_contact_information"`
	DebugNotes                       AGDEBUG       `json:"debug_notes"`
}

func (agr ArchiveGridRecord) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%s%s", agr.MusicianId, agr.Query)
}

//func (agr ArchiveGridRecord) String() string {
//	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Record_archive.href)
//}

func (agr ArchiveGridRecord) String() string {
	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Record_archive)
}

func (agr ArchiveGridRecord) ToJson() string {
	return fmt.Sprintf("{\"ag_record_id\": %q, \n\"musician_id\": %q, \n\"query\": %q, \n}",
		agr.Id, agr.MusicianId, agr.Query)
}

//func (agr ArchiveGridRecord) ToCsv() string {
//	return fmt.Sprintf("%s; %s; %s; %s", agr.Id, agr.MusicianId, agr.Query, agr.Record_archive.href)
//}

func (agr ArchiveGridRecord) ToCsv() string {
	return fmt.Sprintf("%q; %q; %q; %q; %q; %q; %q; %q; %q; %q; %q\n",
		agr.Id,
		agr.MusicianId,
		agr.Query,
		agr.IsFound,
		agr.Record,
		agr.Record_title,
		agr.Record_author,
		agr.Record_archive,
		agr.Record_summary,
		agr.Record_links_contact_information,
		agr.DebugNotes)
}

func (agr ArchiveGridRecord) Hash() HashSum {
	hashfunc := md5.New()
	data := agr.PrimaryKey()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return HashSum(fmt.Sprintf("%x", hashsum))
}

func NewArchiveGridRecord(musicianId HashSum, query MusicianQuery) (archiveGridRecord ArchiveGridRecord) {
	archiveGridRecord = ArchiveGridRecord{
		MusicianId: musicianId,
		Query:      query,
		IsFound:    false,
		//Record:                           ArchiveGridRecordSTRINGNULL,
		//Record_title:                     AGRecordTitle,
		//Record_author:                    AGRecordAuthor,
		//Record_archive:                   AGRecordArchive,
		//Record_summary:                   AGRecordSummary,
		//Record_links_contact_information: AGRecordLinksContactInformation,
	}

	archiveGridRecord.Id = archiveGridRecord.Hash()
	return archiveGridRecord
}

func (agr ArchiveGridRecord) Set() {

}

//

type AGDomPaths struct {
	Record                           string // AGRecord.Dom
	Record_title                     string // AGRecordTitle.Dom
	Record_author                    string // AGRecordAuthor.Dom
	Record_archive                   string // AGRecordArchive.Dom
	Record_summary                   string // AGRecordSummary.Dom
	Record_links_contact_information string // AGRecordLinksContactInformation.Dom
	Results                          string
	ResultsNotEmpty                  string
	ResultsEmpty                     string
	ResultsSize                      string
	ResultsSizeMessage               string
	ResultsNext                      string
}

var AGDomPathsDefinition = AGDomPaths{
	Record:                           "div.record",                // container
	Record_title:                     "div.record_title > h3 > a", // h3>a href ANDTHEN $inner_text
	Record_author:                    "div.record_author",         // span THEN $inner_text
	Record_archive:                   "div.record_archive",        // span THEN $inner_text
	Record_summary:                   "div.record_summary",        // THEN $inner_text
	Record_links_contact_information: "div.record_links",          // a href ANDALSO title
	Results:                          "div.results",
	ResultsNotEmpty:                  "div.results > div.searchresult",
	ResultsEmpty:                     "div.results > div.alertresult",
	ResultsSize:                      "main > h2", // "main h2 > span#resultsize"
	ResultsSizeMessage:               ".navrow span",
	ResultsNext:                      ".results .navtable .navrow a[title=\"View the Next page of results\"]", // get the href

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
//	Record:                           "div.record",                // container
//	Record_title:                     "div.record_title > h3 > a", // h3>a href ANDTHEN $inner_text
//	Record_author:                    "div.record_author",         // span THEN $inner_text
//	Record_archive:                   "div.record_archive",        // span THEN $inner_text
//	Record_summary:                   "div.record_summary",        // THEN $inner_text
//	Record_links_contact_information: "div.record_links",          // a href ANDALSO title
//}
//
////
//
//var ARCHIVE_GRID_BASE_URL = "https://researchworks.oclc.org/archivegrid"
//var AG_BASE_URL, _ = url.Parse(ARCHIVE_GRID_BASE_URL)
//log.Printf("INFO: %v", AG_BASE_URL)
//
//// type ArchiveGridRecord struct {
//// 	RecId                            int
//// 	Record                           AGRecord
//// 	Record_title                     AGRecordTitle
//// 	Record_author                    AGRecordAuthor
//// 	Record_archive                   AGRecordArchive
//// 	Record_summary                   AGRecordSummary
//// 	Record_links_contact_information AGRecordLinksContactInformation
//// }
//
////

package archivegrid

import (
	"crypto/md5"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils/hash"
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

const ArchiveGridRecordSTRINGNULL = "NODATAFOUND"

type RecordHash hash.HashSum

type Record struct {
	Id                       RecordHash            `json:"id"`
	MusicianId               musician.MusicianHash `json:"musician_id"`
	QueryId                  MusicianQueryHash     `json:"query_id"`
	Query                    MusicianQuery         `json:"musician_query"`
	ResultCount              int                   `json:"result_count"`
	IsMatch                  bool                  `json:"is_match"`
	RecordCollectionDataPath string                `json:"record_collection_datapath"`
	Title                    string                `json:"record_title"`
	Author                   string                `json:"record_author"`
	Archive                  string                `json:"record_archive"`
	Summary                  string                `json:"record_summary"`
	LinksContactInformation  string                `json:"links_contact_information"`
	ContactInformation       string                `json:"contact_information"`
	DebugNotes               AGDEBUG               `json:"debug_notes"`
}

func (rec *Record) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%s%s", rec.MusicianId, rec.Query)
}

//func (rec Record) String() string {
//	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", rec.Id, rec.MusicianId, rec.Query, rec.Archive.href)
//}

func (rec *Record) String() string {
	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", rec.Id, rec.MusicianId, rec.Query, rec.Archive)
}

func (rec *Record) ToJson() string {
	return fmt.Sprintf("{\"ag_record_id\": %q, \n\"musician_id\": %q, \n\"query\": %q, \n}",
		rec.Id, rec.MusicianId, rec.Query)
}

//func (rec Record) ToCsv() string {
//	return fmt.Sprintf("%s; %s; %s; %s", rec.Id, rec.MusicianId, rec.Query, rec.Archive.href)
//}

func (rec *Record) ToCsv() string {
	return fmt.Sprintf("%q; %q; %q; %d; %q; %q; %q; %q; %q; %q; %q; %q\n",
		rec.Id,
		rec.MusicianId,
		rec.Query,
		rec.ResultCount,
		rec.RecordCollectionDataPath,
		rec.Title,
		rec.Author,
		rec.Archive,
		rec.Summary,
		rec.LinksContactInformation,
		rec.ContactInformation,
		rec.DebugNotes)
}

func (rec Record) Hash() RecordHash {
	hashfunc := md5.New()
	data := rec.PrimaryKey()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return RecordHash(fmt.Sprintf("%x", hashsum))
}

func NewArchiveGridRecord(musicianId musician.MusicianHash, query MusicianQuery) (archiveGridRecord *Record) {
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

func (rec *Record) Destroy() {
	rec.Id = ""
	rec.MusicianId = ""
	rec.Query = MusicianQuery{}
	rec.ResultCount = 0
	rec.IsMatch = false
	rec.RecordCollectionDataPath = ""
	rec.Title = ""
	rec.Author = ""
	rec.Archive = ""
	rec.Summary = ""
	rec.LinksContactInformation = ""
	rec.ContactInformation = ""
	rec.DebugNotes = AGDEBUG(0)
	return
}

func (rec *Record) Set(record, title, author, archive, summary, link, contact string) {
	rec.IsMatch = false
	rec.RecordCollectionDataPath = record
	rec.Title = title
	rec.Author = author
	rec.Archive = archive
	rec.Summary = summary
	rec.LinksContactInformation = link
	rec.ContactInformation = contact
	rec.DebugNotes = AGDEBUG(FOUNDNOTVALIDATEDYET)
}

func (rec *Record) ContainsAnyFolded(phrases []string) (matches int) {
	if len(phrases) < 1 {
		return -1
	}

	for _, phrase := range phrases {

		p := strings.ToLower(phrase)
		//log.Printf("A PHRASE %s", p)
		//WaitForKeypress()
		if strings.Contains(strings.ToLower(rec.Title), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(rec.Author), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(rec.Archive), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(rec.Summary), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(rec.ContactInformation), p) {
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

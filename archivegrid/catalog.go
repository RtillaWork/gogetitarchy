package archivegrid

import (
	"crypto/md5"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils"
	"io"
	"strings"
)

const CatalogSTRINGNULL = "NODATAFOUND"

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

type Catalog struct {
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

func (agr *Catalog) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%s%s", agr.MusicianId, agr.Query)
}

//func (agr Record) String() string {
//	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Archive.href)
//}

func (agr *Catalog) String() string {
	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Archive)
}

func (agr *Catalog) ToJson() string {
	return fmt.Sprintf("{\"ag_record_id\": %q, \n\"musician_id\": %q, \n\"query\": %q, \n}",
		agr.Id, agr.MusicianId, agr.Query)
}

//func (agr Record) ToCsv() string {
//	return fmt.Sprintf("%s; %s; %s; %s", agr.Id, agr.MusicianId, agr.Query, agr.Archive.href)
//}

func (agr *Catalog) ToCsv() string {
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

func (agr Catalog) Hash() utils.HashSum {
	hashfunc := md5.New()
	data := agr.PrimaryKey()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return utils.HashSum(fmt.Sprintf("%x", hashsum))
}

func NewCatalog(musicianId utils.HashSum, query MusicianQuery) (archiveGridRecord *Record) {
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

func (agr *Catalog) Destroy() {
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

func (agr *Catalog) Set(record, title, author, archive, summary, link, contact string) {
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

func (agr *Catalog) ContainsAnyFolded(phrases []string) (matches int) {
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

type AGCDomPaths struct {
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

var AGCDomPathsDefinition = AGDomPaths{
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

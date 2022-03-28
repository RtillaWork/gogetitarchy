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

func (cat *Catalog) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%s%s", cat.MusicianId, cat.Query)
}

//func (cat Record) String() string {
//	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", cat.Id, cat.MusicianId, cat.Query, cat.Archive.href)
//}

func (cat *Catalog) String() string {
	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", cat.Id, cat.MusicianId, cat.Query, cat.Archive)
}

func (cat *Catalog) ToJson() string {
	return fmt.Sprintf("{\"ag_record_id\": %q, \n\"musician_id\": %q, \n\"query\": %q, \n}",
		cat.Id, cat.MusicianId, cat.Query)
}

//func (cat Record) ToCsv() string {
//	return fmt.Sprintf("%s; %s; %s; %s", cat.Id, cat.MusicianId, cat.Query, cat.Archive.href)
//}

func (cat *Catalog) ToCsv() string {
	return fmt.Sprintf("%q; %q; %q; %d; %q; %q; %q; %q; %q; %q; %q; %q\n",
		cat.Id,
		cat.MusicianId,
		cat.Query,
		cat.ResultCount,
		cat.RecordCollectionDataPath,
		cat.Title,
		cat.Author,
		cat.Archive,
		cat.Summary,
		cat.LinksContactInformation,
		cat.ContactInformation,
		cat.DebugNotes)
}

func (cat Catalog) Hash() utils.HashSum {
	hashfunc := md5.New()
	data := cat.PrimaryKey()
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

func (cat *Catalog) Destroy() {
	cat.Id = ""
	cat.MusicianId = ""
	cat.Query = MusicianQuery{}
	cat.ResultCount = 0
	cat.IsMatch = false
	cat.RecordCollectionDataPath = ""
	cat.Title = ""
	cat.Author = ""
	cat.Archive = ""
	cat.Summary = ""
	cat.LinksContactInformation = ""
	cat.ContactInformation = ""
	cat.DebugNotes = AGDEBUG(0)
	return
}

func (cat *Catalog) Set(record, title, author, archive, summary, link, contact string) {
	cat.IsMatch = false
	cat.RecordCollectionDataPath = record
	cat.Title = title
	cat.Author = author
	cat.Archive = archive
	cat.Summary = summary
	cat.LinksContactInformation = link
	cat.ContactInformation = contact
	cat.DebugNotes = AGDEBUG(FOUNDNOTVALIDATEDYET)
}

func (cat *Catalog) ContainsAnyFolded(phrases []string) (matches int) {
	if len(phrases) < 1 {
		return -1
	}

	for _, phrase := range phrases {

		p := strings.ToLower(phrase)
		//log.Printf("A PHRASE %s", p)
		//WaitForKeypress()
		if strings.Contains(strings.ToLower(cat.Title), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(cat.Author), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(cat.Archive), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(cat.Summary), p) {
			matches++
			//WaitForKeypress()

		}
		if strings.Contains(strings.ToLower(cat.ContactInformation), p) {
			matches++
			//WaitForKeypress()

		}
	}
	return matches
}

//

type CatDomPaths struct {
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

var CatDomPathsDefinition = CatDomPaths{
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

/*
root: main > div.container > div.row > div.col-12-md

	div.row div.col-md-12
		h2 text=$RecordName
		h4.catalogrecordcreator text=$RecordCreator_MAYBE_WITH_DATE

	div.row.pad-above > div.col-md-8
		div.catalogrecordcontact
		h5
		div.catalogrecordfield(1)
		div.catalogrecordfield(2) text=$Description
		div.catalogrecordfield(3) text=$Info
		div.catalogrecordfield(4)
		div.catalogrecordfield(5)
		div.catalogrecordfield(6)
		div.catalogrecordfield(7)
		h5
		div.catalogrecordfield(8)
			a[href]=$OnlineAidLink  text=$OnlineAidDescription
		div.catalogrecordfield(9)
			a[href]=$WorldCatOCLCLink    #https://www.worldcat.org/oclc/890209766

	div.col-md-4#sidebar





*/

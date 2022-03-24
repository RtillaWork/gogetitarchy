package main

import (
	"crypto/md5"
	"fmt"
	"io"
)

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

const ArchiveGridRecordSTRINGNULL = "NODATAFOUND"

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

func (agr ArchiveGridRecord) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%s%s", agr.MusicianId, agr.Query)
}

func (agr ArchiveGridRecord) String() string {
	return fmt.Sprintf("{ RECORDID%sMUSICIAN%s__%s_in_%s }", agr.Id, agr.MusicianId, agr.Query, agr.Record_archive.href)
}

func (agr ArchiveGridRecord) ToJson() string {
	return fmt.Sprintf("{\"ag_record_id\": %s, \n\"musician_id\": %s, \n\"query\": %s, \n}", agr.Id, agr.MusicianId, agr.Query)
}

func (agr ArchiveGridRecord) ToCsv() string {
	return fmt.Sprintf("%s; %s; %s; %s", agr.Id, agr.MusicianId, agr.Query, agr.Record_archive.href)
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

//
var AGDomPathsDefinition = AGDomPaths{
	Record:                           "div.record",                // container
	Record_title:                     "div.record_title > h3 > a", // h3>a href ANDTHEN $inner_text
	Record_author:                    "div.record_author",         // span THEN $inner_text
	Record_archive:                   "div.record_archive",        // span THEN $inner_text
	Record_summary:                   "div.record_summary",        // THEN $inner_text
	Record_links_contact_information: "div.record_links",          // a href ANDALSO title
}

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

package query

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/archivegrid"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"github.com/RtillaWork/gogetitarchy/utils/hash"
	"io"
	"log"
	"net/url"
	"time"
)

var ARCHIVE_GRID_URL_TEMPLATE []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%s&limit=100",
	// REPLACE 100 by QUERY_LIMIT
	//"https://researchworks.oclc.org/archivegrid/?q=%s&limit=10",
}

type QUERYDEBUG int

const (
	TOOMANYRESULTS = iota
	NORESULTS
	ACCEPTABLERESULTS
	ERROR
)
const QUERY_LIMIT = 100

type MusicianQueryHash hash.HashSum

//type Queries map[utils.HashSum]*Query
type Queries map[musician.MusicianHash]*Query

type Query struct {
	Id         MusicianQueryHash     `json:"query_id"` // for now init to same as MusicianId one musician one query
	MusicianId musician.MusicianHash `json:"musician_id"`
	Url        string                `json:"url"`
	Timestamp  time.Time             `json:"timestamp"`   // should be initialized to a NEVERQUERIED value
	ResultSize int                   `json:"result_size"` // should be initialized to a NEVERQUERIED value
	Matches    int                   `json:"Matches"`
	DebugNotes QUERYDEBUG            `json:"debug_notes"`
}

func (mq *Query) String() string {
	return string(mq.Url)
}

func NewQuery(id musician.MusicianHash, url string) (newMusicianQuery *Query) {
	newMusicianQuery = new(Query)
	newMusicianQuery = &Query{
		Id:         MusicianQueryHash(id),
		MusicianId: id,
		Url:        url,
		// Timestamp should be initialized to a NEVERQUERIED value
		ResultSize: -1, //should be initialized to a NEVERQUERIED value
		Matches:    -1,
	}
	newMusicianQuery.Id = newMusicianQuery.Hash()

	return newMusicianQuery

	//return &Query{
	//	Id: id,
	//	//MusicianId: m.Id,
	//	Url: url,
	//	// Timestamp should be initialized to a NEVERQUERIED value
	//	ResultSize: -1, //should be initialized to a NEVERQUERIED value
	//}
}

func (mq *Query) ToJson() string {
	jsoned, err := json.Marshal(*mq)
	errors.FailOn(err, "Musician::ToJson json.Marshal")
	return fmt.Sprintf("%s", string(jsoned))
}

func (mq *Query) Hash() MusicianQueryHash {
	hashfunc := md5.New()
	// NOTE: assume Musician::String() is unique. Needs assertion, or else expand the Sum() contents
	data := mq.ToJson()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return MusicianQueryHash(fmt.Sprintf("%x", hashsum))
}

func (mq *Query) SetResultCount(count int) {
	mq.Timestamp = time.Now()
	mq.ResultSize = count
}

func (mq *Query) SetResultCountFunc(f func(query Query) (int, error)) {
	resultsize, err := archivegrid.ScanQueryResultSize(*mq)
	if err != nil {
		mq.SetResultCount(-1)
		mq.DebugNotes = QUERYDEBUG(ERROR)

		log.Printf("RESULT SIZE resultsSize == 0 || err != nil %d", resultsize)

		return
	} else if resultsize == 0 {
		mq.SetResultCount(0)
		mq.DebugNotes = QUERYDEBUG(NORESULTS)

		log.Printf("RESULT SIZE resultsSize == 0 || err != nil %d", resultsize)

		return
	} else if resultsize > archivegrid.TOOMANYRESULTSVALUE {
		log.Printf("RESULT SIZE resultsSize > TOOMANYRESULTSVALUE %d", resultsize)
		// too many to process for now, take note and pass, set ResultSize false as flag nor record as non nilfor now
		mq.SetResultCount(0)
		mq.DebugNotes = QUERYDEBUG(TOOMANYRESULTS)

	} else {
		log.Printf("RESULT SIZE ok supposed to process the other OnHtml for AG DOM elements %d", resultsize)
		mq.SetResultCount(resultsize)
		mq.DebugNotes = QUERYDEBUG(ACCEPTABLERESULTS)
	}
}

func (mq *Query) Destroy() {
	mq.Id = ""
	//mq.MusicianId
	mq.Url = ""
	mq.Timestamp = time.Time{}
	mq.ResultSize = 0
	mq.DebugNotes = QUERYDEBUG(0)
	return
}

func BuildQueries(ms musician.MusiciansMap) (mq Queries) {
	mq = Queries{}

	for _, m := range ms {
		query := buildQuery(m, ARCHIVE_GRID_URL_TEMPLATE[0], musician.NamesVariation(musician.FULL))
		mq[m.Id] = query
	}

	return mq
}

func buildQuery(m *musician.Musician, template string, variation musician.NamesVariation) *Query {
	querydata := url.QueryEscape(m.QueryFragment(variation))
	fullquery := fmt.Sprintf(template, querydata)

	return NewQuery(m.Id, fullquery)

	//return Query{
	//	Id: m.Id,
	//	//MusicianId: m.Id,
	//	Url: fullquery,
	//	// Timestamp should be initialized to a NEVERQUERIED value
	//	ResultSize: -1, //should be initialized to a NEVERQUERIED value
	//}
}

//func BuildQuery(m Musician, template string) Query {
//	query := url.QueryEscape(fmt.Sprintf(ARCHIVE_GRID_URL_TEMPLATE[0], m.FullName()))
//	queries := []string{query}
//	return MusicianQueries{m.id: queries}
//}

///////////

// var ARCHIVE_GRID_URL_PATTERNS []string = []string{
// 	"https://researchworks.oclc.org/archivegrid/?q=%s&limit=100",
// }

///////////////

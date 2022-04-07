package archivegrid

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
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

//type MusiciansQueries map[utils.HashSum]*MusicianQuery
type MusiciansQueries map[musician.MusicianHash]*MusicianQuery

type MusicianQuery struct {
	Id         MusicianQueryHash     `json:"query_id"` // for now init to same as MusicianId one musician one query
	MusicianId musician.MusicianHash `json:"musician_id"`
	Url        string                `json:"url"`
	Timestamp  time.Time             `json:"timestamp"`   // should be initialized to a NEVERQUERIED value
	ResultSize int                   `json:"result_size"` // should be initialized to a NEVERQUERIED value
	Matches    int                   `json:"Matches"`
	DebugNotes QUERYDEBUG            `json:"debug_notes"`
}

func (mq *MusicianQuery) String() string {
	return string(mq.Url)
}

func NewMusicianQuery(id musician.MusicianHash, url string) (newMusicianQuery *MusicianQuery) {
	newMusicianQuery = new(MusicianQuery)
	newMusicianQuery = &MusicianQuery{
		Id:         MusicianQueryHash(id),
		MusicianId: id,
		Url:        url,
		// Timestamp should be initialized to a NEVERQUERIED value
		ResultSize: -1, //should be initialized to a NEVERQUERIED value
		Matches:    -1,
	}
	newMusicianQuery.Id = newMusicianQuery.Hash()

	return newMusicianQuery

	//return &MusicianQuery{
	//	Id: id,
	//	//MusicianId: m.Id,
	//	Url: url,
	//	// Timestamp should be initialized to a NEVERQUERIED value
	//	ResultSize: -1, //should be initialized to a NEVERQUERIED value
	//}
}

func (mq *MusicianQuery) ToJson() string {
	jsoned, err := json.Marshal(*mq)
	errors.FailOn(err, "Musician::ToJson json.Marshal")
	return fmt.Sprintf("%s", string(jsoned))
}

func (mq *MusicianQuery) Hash() MusicianQueryHash {
	hashfunc := md5.New()
	// NOTE: assume Musician::String() is unique. Needs assertion, or else expand the Sum() contents
	data := mq.ToJson()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return MusicianQueryHash(fmt.Sprintf("%x", hashsum))
}

func (mq *MusicianQuery) SetResultCount(count int) {
	mq.Timestamp = time.Now()
	mq.ResultSize = count
}

func (mq *MusicianQuery) SetResultCountFunc(f func(query MusicianQuery) (int, error)) {
	resultsize, err := ScanQueryResultSize(*mq)
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
	} else if resultsize > TOOMANYRESULTSVALUE {
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

func (mq *MusicianQuery) Destroy() {
	mq.Id = ""
	//mq.MusicianId
	mq.Url = ""
	mq.Timestamp = time.Time{}
	mq.ResultSize = 0
	mq.DebugNotes = QUERYDEBUG(0)
	return
}

func BuildQueries(ms musician.MusiciansMap) (mq MusiciansQueries) {
	mq = MusiciansQueries{}

	for _, m := range ms {
		query := buildQuery(m, ARCHIVE_GRID_URL_TEMPLATE[0], musician.NamesVariation(musician.FULL))
		mq[m.Id] = query
	}

	return mq
}

func buildQuery(m *musician.Musician, template string, variation musician.NamesVariation) *MusicianQuery {
	querydata := url.QueryEscape(m.QueryFragment(variation))
	fullquery := fmt.Sprintf(template, querydata)

	return NewMusicianQuery(m.Id, fullquery)

	//return MusicianQuery{
	//	Id: m.Id,
	//	//MusicianId: m.Id,
	//	Url: fullquery,
	//	// Timestamp should be initialized to a NEVERQUERIED value
	//	ResultSize: -1, //should be initialized to a NEVERQUERIED value
	//}
}

//func BuildQuery(m Musician, template string) MusicianQuery {
//	query := url.QueryEscape(fmt.Sprintf(ARCHIVE_GRID_URL_TEMPLATE[0], m.FullName()))
//	queries := []string{query}
//	return MusicianQueries{m.id: queries}
//}

///////////

// var ARCHIVE_GRID_URL_PATTERNS []string = []string{
// 	"https://researchworks.oclc.org/archivegrid/?q=%s&limit=100",
// }

///////////////

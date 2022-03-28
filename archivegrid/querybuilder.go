package archivegrid

import (
	"fmt"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils"
	"net/url"
	"time"
)

var ARCHIVE_GRID_URL_TEMPLATE []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%s&limit=100",
	//"https://researchworks.oclc.org/archivegrid/?q=%s&limit=10",
}

type QUERYDEBUG int

const (
	TOOMANYRESULTS = iota
	NORESULTS
	ACCEPTABLERESULTS
)

type MusicianQuery struct {
	Id utils.HashSum `json:"query_id"` // for now init to same as MusicianId one musician one query
	//MusicianId HashSum  `json:"musician_id"`
	Url        string     `json:"url"`
	Timestamp  time.Time  `json:"timestamp"`   // should be initialized to a NEVERQUERIED value
	ResultSize int        `json:"result_size"` // should be initialized to a NEVERQUERIED value
	Matches    int        `json:"Matches"`
	DebugNotes QUERYDEBUG `json:"debug_notes"`
}

func (mq *MusicianQuery) String() string {
	return string(mq.Url)
}

func NewMusicianQuery(id utils.HashSum, url string) (newMusicianQuery *MusicianQuery) {
	newMusicianQuery = new(MusicianQuery)
	newMusicianQuery = &MusicianQuery{
		Id: id,
		//MusicianId: m.Id,
		Url: url,
		// Timestamp should be initialized to a NEVERQUERIED value
		ResultSize: -1, //should be initialized to a NEVERQUERIED value
		Matches:    -1,
	}

	return newMusicianQuery

	//return &MusicianQuery{
	//	Id: id,
	//	//MusicianId: m.Id,
	//	Url: url,
	//	// Timestamp should be initialized to a NEVERQUERIED value
	//	ResultSize: -1, //should be initialized to a NEVERQUERIED value
	//}
}

func (mq *MusicianQuery) SetResultCount(count int) {
	mq.Timestamp = time.Now()
	mq.ResultSize = count
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

//type MusiciansQueries map[HashSum][]MusicianQuery
type MusiciansQueries map[utils.HashSum]*MusicianQuery

func BuildQueries(ms musician.MusiciansMap) (mq MusiciansQueries) {
	mq = MusiciansQueries{}

	for _, m := range ms {
		query := buildQuery(m, ARCHIVE_GRID_URL_TEMPLATE[0], musician.MusicianNamesVariation(musician.FULL))
		mq[m.Id] = query
	}

	return mq
}

func buildQuery(m *musician.Musician, template string, variation musician.MusicianNamesVariation) *MusicianQuery {
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
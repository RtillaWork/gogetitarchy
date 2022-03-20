package main

import (
	"fmt"
	"net/url"
)

var ARCHIVE_GRID_URL_TEMPLATE []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%s&limit=100",
}

type MusicianQuery string
type MusiciansQueries map[HashSum][]MusicianQuery

func BuildQueries(ms MusiciansMap) MusiciansQueries {
	mq := MusiciansQueries{}

	for _, m := range ms {
		query := BuildQuery(m, ARCHIVE_GRID_URL_TEMPLATE[0], MusicianNamesVariation(FIRSTNAMELASTNAME))
		mq[m.Id] = []MusicianQuery{query}
	}

	return mq
}

func BuildQuery(m Musician, template string, variation MusicianNamesVariation) MusicianQuery {
	querydata := m.NameFmt(variation)
	fullquery := url.QueryEscape(fmt.Sprintf(template, querydata))
	return MusicianQuery(fullquery)
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

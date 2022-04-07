package archivegrid

import "github.com/RtillaWork/gogetitarchy/musician"

const TOOMANYRESULTSVALUE int = 15000

var ALLOWED_DOMAINS []string = []string{"researchworks.oclc.org", "archives.chadwyck.com", "www.newspapers.com"}
var ARCHIVE_GRID_URL_PATTERNS []string = []string{
	"https://researchworks.oclc.org/archivegrid/?q=%22Albert+Quincy+Porter%22",
}

type MusiciansData map[musician.MusicianHash][]*Record

type FilteredMusiciansData map[musician.MusicianHash][]*Record

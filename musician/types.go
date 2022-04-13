package musician

import (
	"github.com/RtillaWork/gogetitarchy/utils/hash"
	"time"
)

type MusiciansMap map[MusicianHash]*Musician
type MusiciansDb struct {
	Musicians *MusiciansMap
	Dict      *DataDict
}

var Defaults Musician

type MusicianHash hash.HashSum

//const INT64_NULL = 9223372036854775807 // Max int64
const AGE_NULL = 0
const STRING_NULL = "STRINGNULL"
const NAMES_DEFAULT_SEP = " "
const LAST_NAME_SEP = ","
const INITIALS_SEP = ". " // I.N.I.T._NAMES
const NOTES_SEP_OPEN = "("
const NOTES_SEP_CLOSE = ")"
const FIELDS_SEP = ",: -"
const CSV_SEP = ";"

// an impossible time for the Domain, to signify a null
var TIME_NULL time.Time = time.Date(2022, time.March, 01, 00, 00, 00, 00, time.UTC)

type Musician struct { // nils, 0s are not valid to represent missing information
	// TODO assertion: creating a Musician -> no field is nil
	Id          MusicianHash      `json:"id"`
	RawName     string            `json:"raw_name"`
	FName       string            `json:"first_name"`
	LName       string            `json:"last_name"`
	MName       string            `json:"middle_name"`
	Notes       string            `json:"notes"`
	Confidence  int               `json:"confidence"`
	TimeCreated int64             `json:"encounter"`
	Fields      map[string]string `json:"fields"`
	Tags        []string          `json:"tags"`

	// FIELDS:
	// FNAME, MNAME, LNAME, MISCS, DATEBEGIN, DATEEND, DATEOTHER, DATES
	// Fields []string
	//L, F  || F M L || F M. L || F L || F "M" L
	//Military Unit:
	//Estimated Birth Year: y0 - y1
	//Year: y
	//Enlistment Rank: rank
	///Enlistment: y - rank
	///Enlistment: d m y - rankorbrigade - rank
	///Enlistment: d m y - rankorbrigade - rank
	//Branch: Union Army
	//Civil war (union): unit ( Union )
	///Civil war (union): Army - 97th US Colored Infantry - D,A - United States of America
	//Discharge: rank
	//Company: letter
	//Military Unit: unit, unit
	///Military Unit: unit
	//Birth: y - city, county
	//Death: d m y - city, state
	//DateOfBirth   time.Time `json:"dateofbirth"`
	//EstimatedBirthYear
	//DateOfDeath   time.Time `json:"dateofdeath"`
	//PleaceOfBirth string    `json:"placeofbirth"`
	//PlaceOfDeath  string    `json:"placeofdeath"`
	//Age           byte      `json:"age"`
	//Bio           string    `json:"bio"` // other
	// Army string
	// Enlistement
	// EnlistementDate
	// Discharge
	// Rank string
	// Branch
	// Company
	// MiscDate

}

type NamesVariation int

const (
	FULL NamesVariation = iota
	L
	FL
	LFM
)

func init() {
	Defaults = Musician{
		Id:          MusicianHash(""),
		RawName:     "NULL_RAWNAME",
		FName:       "NULL_FIRSTNAME",
		MName:       "NULL_MIDDLENAME",
		LName:       "NULL_LASTNAME",
		Notes:       "NULL_NOTES",
		Confidence:  -100,
		TimeCreated: 0,
		Fields: map[string]string{
			"FIRSTNAME":   "NULL_FIRSTNAME",
			"MIDDLENAME":  "NULL_MIDDLENAME",
			"LASTNAME":    "NULL_LASTNAME",
			"MISCELLANEA": "NULL_MISCELLANEA", // TODO RENAME TO MISCELLANEA because MISC might already be in the text
			"DATEBEGIN":   "NULL_DATEBEGIN",
			"DATEEND":     "NULL_DATEEND",
			"DATEOTHER":   "NULL_DATEOTHER",
			"DATECSV":     "NULL_DATEBCSV", //"xxxx;yyyy;....zzzz"

		}, //map[string]string{},
		Tags: []string{},
	}

}

//newMusician = new(Musician)
//newMusician.Id = MusicianHash(STRING_NULL)
//newMusician.FName = STRING_NULL
//newMusician.MName = STRING_NULL
//newMusician.LName = STRING_NULL
//newMusician.Notes = STRING_NULL
//newMusician.Fields = map[string]string{}
//newMusician.Tags = []string{}
//ok = false

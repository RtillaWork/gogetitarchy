package musician

import (
	"github.com/RtillaWork/gogetitarchy/utils/hash"
	"time"
)

type MusiciansMap map[MusicianHash]*Musician

//const INT64_NULL = 9223372036854775807 // Max int64
const AGE_NULL = 0
const STRING_NULL = "STRINGNULL"
const NAMES_DEFAULT_SEP = " "
const LAST_NAME_SEP = ","
const INITIALS_SEP = ". " // I.N.I.T._NAMES
const NOTES_SEP_OPEN = "("
const NOTES_SEP_CLOSE = ")"
const FIELDS_SEP = ",: -"

// an impossible time for the Domain, to signify a null
var TIME_NULL time.Time = time.Date(2022, time.March, 01, 00, 00, 00, 00, time.UTC)
var Defaults Musician

type MusicianHash hash.HashSum

type Musician struct { // nils, 0s are not valid to represent missing information
	// TODO assertion: creating a Musician -> no field is nil
	Id         MusicianHash      `json:"id"`
	FirstName  string            `json:"first_name"`
	LastName   string            `json:"last_name"`
	MiddleName string            `json:"middle_name"`
	Notes      string            `json:"notes"`
	Encounter  uint8             `json:"encounter"`
	Fields     map[string]string `json:"fields"`
	Tags       []string          `json:"tags"`
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

func init() {
	Defaults = Musician{
		Id:         MusicianHash(""),
		FirstName:  "NULL_FIRSTNAME",
		MiddleName: "NULL_MIDDLENAME",
		LastName:   "NULL_LASTNAME",
		Notes:      "NULL_NOTES",
		Encounter:  0,
		Fields:     map[string]string{},
		Tags:       []string{},
	}

}

//newMusician = new(Musician)
//newMusician.Id = MusicianHash(STRING_NULL)
//newMusician.FirstName = STRING_NULL
//newMusician.MiddleName = STRING_NULL
//newMusician.LastName = STRING_NULL
//newMusician.Notes = STRING_NULL
//newMusician.Fields = map[string]string{}
//newMusician.Tags = []string{}
//ok = false

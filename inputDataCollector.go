package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"strings"
	"time"
)

//const INT64_NULL = 9223372036854775807 // Max int64
const AGE_NULL = 0
const STRING_NULL = "STRING_NULL"

const NAMES_DEFAULT_SEP = " "
const LAST_NAME_SEP = ","
const INITIALS_SEP = ". " // I.N.I.T._NAMES
const NOTES_SEP_OPEN = "("
const NOTES_SEP_CLOSE = ")"

// an impossible time for the Domain, to signify a null
var TIME_NULL time.Time = time.Date(2022, time.March, 01, 00, 00, 00, 00, time.UTC)

type HashSum string

func (h HashSum) String() string {
	return string(h)
}

type Musician struct { // nils, 0s or  mean no information
	// TODO assertion: creating a Musician -> no field is nil
	// MD5 on aMusician.String()
	Id         HashSum `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName string  `json:"middle_name"`
	Notes      string  `json:"notes"`
	//DateOfBirth   time.Time `json:"dateofbirth"`
	//DateOfDeath   time.Time `json:"dateofdeath"`
	//PleaceOfBirth string    `json:"placeofbirth"`
	//PlaceOfDeath  string    `json:"placeofdeath"`
	//Age           byte      `json:"age"`
	//Bio           string    `json:"bio"` // other

	// Army string
	// Rank string
}

var MusicianNULL = Musician{
	"NULL_HASH",
	STRING_NULL,
	STRING_NULL,
	STRING_NULL,
	STRING_NULL,
	//TIME_NULL,
	//TIME_NULL,
	//STRING_NULL,
	//STRING_NULL,
	//AGE_NULL,
	//STRING_NULL,

	// Army string
	// Rank string
}

func (m Musician) String() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("%s_%s_%s", first, middle, last)
}

func (m Musician) FullNameTuple() (string, bool, string, bool, string, bool) { //  firstname, middlename, lastname
	firstname := STRING_NULL
	middlename := STRING_NULL
	lastname := STRING_NULL

	isFirstNamePresent := m.FirstName != STRING_NULL
	isMiddleNamePresent := m.MiddleName != STRING_NULL
	isLastNamePresent := m.LastName != STRING_NULL
	FailNotOK(isLastNamePresent, "Musician#FullNameTuple NO LASTNAME")
	lastname = m.LastName

	if isFirstNamePresent {
		firstname = m.FirstName
	}

	if isMiddleNamePresent {
		middlename = m.MiddleName
	}
	return firstname, isFirstNamePresent, middlename, isMiddleNamePresent, lastname, isLastNamePresent
}

func (m Musician) FullName() string {
	first, isFirstPresent, middle, isMiddlePresent, last, _ := m.FullNameTuple()
	if isFirstPresent {
		first = first + NAMES_DEFAULT_SEP
	}
	if isMiddlePresent {
		middle = middle + NAMES_DEFAULT_SEP
	}

	return fmt.Sprintf("%s%s%s", first, middle, last)
}

type MusicianNamesVariation int

const (
	FULLNAME MusicianNamesVariation = iota
	LASTNAME
	FIRSTNAMELASTNAME
	LASTNAMEFIRSTNAMEMIDDLENAME
)

func (m Musician) NameFmt(v MusicianNamesVariation) string {
	formattedName := ""
	switch v {
	case MusicianNamesVariation(FULLNAME):
		formattedName = m.FullName()
	case MusicianNamesVariation(LASTNAME):
		_, _, _, _, last, _ := m.FullNameTuple()
		formattedName = fmt.Sprintf("%s", last)
	case MusicianNamesVariation(FIRSTNAMELASTNAME):
		first, isFirstPresent, _, _, last, _ := m.FullNameTuple()
		if !isFirstPresent {
			first = ""
		}
		formattedName = fmt.Sprintf("%s %s", first, last)
	case MusicianNamesVariation(LASTNAMEFIRSTNAMEMIDDLENAME):
		first, isFirstPresent, middle, isMiddlePresent, last, _ := m.FullNameTuple()
		if !isFirstPresent {
			first = ""
		}
		if !isMiddlePresent {
			middle = ""
		}
		formattedName = fmt.Sprintf("%s, %s %s.", last, first, middle)
	default:
		formattedName = m.FullName()

	}
	return formattedName
}

func (m Musician) Hash() HashSum {
	hashfunc := md5.New()
	// NOTE: assume Musician::String() is unique. Needs assertion, or else expand the Sum() contents
	hashsum := hashfunc.Sum([]byte(m.String()))
	return HashSum(hashsum)
}

func NewMusician(data string) Musician {
	aMusician := Musician{}

	notes, oknotes, names, okmore := ExtractNotes(data)
	FailNotOK(okmore, "NewMusician Try to ExctractNotes( FAILED TO FIND NAMES")

	if oknotes {
		aMusician.Notes = notes
	}

	firstname, middlename, lastname, ok := ExtractNames(names)
	FailNotOK(ok, "NewMusician try to ExtractNames( FAILED FOR UNKNOWN REASONS")

	aMusician.FirstName = firstname
	aMusician.MiddleName = middlename
	aMusician.LastName = lastname
	aMusician.id = aMusician.Hash()

	return aMusician
}

func ExtractNotes(data string) (string, bool, string, bool) {
	// returns (notes if any, truncated data without notes if any, oknotes, okmore)
	notes := STRING_NULL
	foundnotes := false
	rest := data
	foundmore := false

	fields := strings.Fields(data)
	switch len(fields) {
	case 0:
		//return notes,  false, rest, false
	case 1:
		// NOTE: assume if only one field
		if strings.HasPrefix(fields[0], NOTES_SEP_OPEN) && strings.HasSuffix(fields[0], NOTES_SEP_CLOSE) {
			// CASE (ABC) -> notes
			// NOTE: should panic as it means no names in the scanned line
			notes = fields[0]
			//return notes, false, rest,  false
		} else if !strings.ContainsAny(fields[0], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
			// CASE ABC -> lastname
			rest = fields[0]
			foundmore = true
			//return notes, false,  rest, true
		} else {
			// CASE badly formed ex (ABC\n or \ABC)
			//return notes, false, data, false
		}

	case 2:
		//
		if strings.HasPrefix(fields[1], NOTES_SEP_OPEN) && strings.HasSuffix(fields[1], NOTES_SEP_CLOSE) {
			// CASE XYZ (ABC) -> notes
			// NOTE: should panic as it means no names in the scanned line
			notes = fields[1]
			foundnotes = true
			rest = fields[0]
			foundmore = true
			//return notes, true, rest, true
		} else if !strings.ContainsAny(fields[1], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
			// CASE XYZ ABC  -> rest
			foundmore = true
			//return notes, false, rest,  true
		} else {
			// CASE badly formed
			//return  notes,  false, data,false
		}

	case 3:
		//
		if strings.HasPrefix(fields[2], NOTES_SEP_OPEN) && strings.HasSuffix(fields[2], NOTES_SEP_CLOSE) {
			// CASE X YZ (ABC) -> notes
			// NOTE: should panic as it means no names in the scanned line
			notes = fields[2]
			foundnotes = true
			rest = fields[0] + NAMES_DEFAULT_SEP + fields[1]
			foundmore = true
			//return notes,  true, rest, true
		} else if !strings.ContainsAny(fields[2], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
			// CASE X YZ ABC  -> rest
			rest = fields[0] + NAMES_DEFAULT_SEP + fields[1] + NAMES_DEFAULT_SEP + fields[2]
			foundmore = true
			//return notes, false, rest, true
		} else {
			// CASE badly formed
			//return  notes, false, data, false
		}

	case 4:
		//
		if strings.HasPrefix(fields[3], NOTES_SEP_OPEN) && strings.HasSuffix(fields[3], NOTES_SEP_CLOSE) {
			// CASE X Y Z (ABC) -> notes
			// NOTE: should panic as it means no names in the scanned line
			notes = fields[3]
			foundnotes = true
			rest = fields[0] + NAMES_DEFAULT_SEP + fields[1] + NAMES_DEFAULT_SEP + fields[2]
			foundmore = true
			//return notes,true,  rest, true
		} else {
			// CASE badly formed X Y Z ABC  -> ??
			//return  notes, false, data, false
		}

	default:

	}
	return notes, foundnotes, rest, foundmore
}

func ExtractNames(data string) (string, string, string, bool) {
	// returns (firstname if any, middlename if any, lastname if any, ok)
	firstname := STRING_NULL
	middlename := STRING_NULL
	lastname := STRING_NULL
	ok := false
	//rest := data

	fields := strings.Fields(data)
	switch len(fields) {
	case 0:
		//return _, _ ,_ , false
	case 1:
		// NOTE: assume if only one field
		if strings.HasSuffix(fields[0], LAST_NAME_SEP) {
			// CASE ABC,
			// NOTE: should panic as it means no names in the scanned line
			lastname = strings.TrimSuffix(fields[0], LAST_NAME_SEP)
			ok = true

		} else if !strings.ContainsAny(fields[0], LAST_NAME_SEP+NAMES_DEFAULT_SEP) {
			// CASE ABC -> lastname
			lastname = fields[0]
			ok = true

		} else {
			// CASE badly formed

		}

	case 2:
		//
		if strings.HasSuffix(fields[0], LAST_NAME_SEP) {
			// CASE XYZ, ABC
			lastname = strings.TrimSuffix(fields[0], LAST_NAME_SEP)
			firstname = fields[1]
			ok = true

		} else if !strings.ContainsAny(fields[0], LAST_NAME_SEP+NAMES_DEFAULT_SEP) {
			// CASE XYZ ABC
			firstname = fields[0]
			lastname = fields[1]
			ok = true
		} else {
			// CASE badly formed

		}

	case 3:
		//
		if strings.HasSuffix(fields[0], LAST_NAME_SEP) {
			// CASE XYZ, ABC I.
			lastname = strings.TrimSuffix(fields[0], LAST_NAME_SEP)
			firstname = fields[1]
			middlename = fields[2]
			ok = true

		} else if !strings.ContainsAny(fields[0], LAST_NAME_SEP+NAMES_DEFAULT_SEP) {
			// CASE XYZ ABC
			firstname = fields[0]
			middlename = fields[1]
			lastname = fields[2]
			ok = true
		} else {
			// CASE badly formed

		}

	default:

	}
	return firstname, middlename, lastname, ok
}

type MusiciansMap map[HashSum]Musician

func ReadMusicianData(inFileName string) MusiciansMap {

	inFile, err := os.Open(inFileName)
	FailOn(err, "opening inFile for reading...")
	defer inFile.Close()

	s := bufio.NewScanner(inFile)
	scanline := s.Text()
	aMusician := NewMusician(scanline)
	musicians := map[HashSum]Musician{aMusician.Hash(): aMusician}
	return musicians

}

//func ExtractDataFromString(data string) (string, string, string, string, bool){
//	firstname := STRING_NULL
//	middlename := STRING_NULL
//	lastname := STRING_NULL
//	notes := STRING_NULL
//
//	fields := strings.Fields(data)
//	switch len(fields) {
//	case 0:
//		return firstname, middlename, lastname, notes, false
//	case 1:
//		// NOTE: assume if only one field
//		if strings.HasPrefix(fields[0], NOTES_SEP_OPEN) && strings.HasSuffix(fields[0], NOTES_SEP_CLOSE) {
//			// CASE (ABCD) -> notes
//			return firstname, middlename, lastname, fields[0], true
//		} else if !strings.ContainsAny(fields[0], NOTES_SEP_OPEN + NOTES_SEP_CLOSE ){
//			// CASE ABCD -> lastname
//			return firstname, middlename, fields[0], notes, true
//
//		} else {
//			// CASE badly formed
//			return firstname, middlename, lastname, fields[0], false
//		}
//
//	case 2:
//		// NOTE: assume 2 fields
//		if strings.HasPrefix(fields[1], NOTES_SEP_OPEN) && strings.HasSuffix(fields[1], NOTES_SEP_CLOSE) {
//			// CASE (ABCD) -> notes
//			return firstname, middlename, lastname, fields[0], true
//		} else if !strings.ContainsAny(fields[0], NOTES_SEP_OPEN + NOTES_SEP_CLOSE ){
//			// CASE ABCD -> lastname
//			return firstname, middlename, fields[0], notes, true
//
//		} else {
//			// CASE badly formed
//			return firstname, middlename, lastname, fields[0], false
//		}
//
//	}
//
//
//	for _, field := range fields {
//
//	}
//
//}

package main

import (
	"crypto/md5"
	"fmt"
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

type HashSum string

func (h HashSum) String() string {
	return string(h)
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
	aMusician.Id = aMusician.Hash()

	return aMusician
}

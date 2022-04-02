package musician

import (
	"crypto/md5"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"io"
	"strings"
)

func (m *Musician) String() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("%s_%s_%s", first, middle, last)
}

func (m *Musician) PrimaryKey() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("PRIMARYKEY=%s%s%s", first, middle, last)
}

func (m *Musician) ToCsv() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	id := m.Id
	return fmt.Sprintf("%q; %q; %q; %q", id, first, middle, last)
}

func (m *Musician) ToJson() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	id := m.Id
	return fmt.Sprintf("{ \"id\": %q,\n \"first_name\": %q,\n \"middle_name\": %q,\n \"last_name\": %q\n}", id, first, middle, last)
}

func (m *Musician) QueryFragment(v MusicianNamesVariation) string {
	notes := ""
	if m.Notes != STRING_NULL {
		notes = m.Notes
	}
	return fmt.Sprintf("%s %s", m.NameFmt(v), notes)
}

//
func (m *Musician) FullNameTuple() (firstname string, isFirstNamePresent bool, middlename string, isMiddleNamePresent bool, lastname string, isLastNamePresent bool) { //  firstname, middlename, lastname
	//firstname := STRING_NULL
	//middlename := STRING_NULL
	//lastname := STRING_NULL

	firstname = STRING_NULL
	middlename = STRING_NULL
	lastname = STRING_NULL
	//isFirstNamePresent := m.FirstName != STRING_NULL
	//isMiddleNamePresent := m.MiddleName != STRING_NULL
	//isLastNamePresent := m.LastName != STRING_NULL
	isFirstNamePresent = m.FirstName != STRING_NULL
	isMiddleNamePresent = m.MiddleName != STRING_NULL
	isLastNamePresent = m.LastName != STRING_NULL
	errors.FailNotOK(isLastNamePresent, "Musician#FullNameTuple NO LASTNAME")
	lastname = m.LastName

	if isFirstNamePresent {
		firstname = m.FirstName
	}

	if isMiddleNamePresent {
		middlename = m.MiddleName
	}
	return firstname, isFirstNamePresent, middlename, isMiddleNamePresent, lastname, isLastNamePresent
}

//
func (m *Musician) FullName() string {
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
	FULL MusicianNamesVariation = iota
	LAST
	FIRSTLAST
	LASTFIRSTMIDDLE
)

func (m *Musician) NameFmt(v MusicianNamesVariation) (formattedName string) {
	formattedName = ""
	switch v {
	case MusicianNamesVariation(FULL):
		first, isFirstPresent, middle, isMiddlePresent, last, _ := m.FullNameTuple()
		if !isFirstPresent {
			first = ""
		}
		if !isMiddlePresent {
			middle = ""
		}
		formattedName = fmt.Sprintf("%s %s %s", first, middle, last)
	case MusicianNamesVariation(LAST):
		_, _, _, _, last, _ := m.FullNameTuple()
		formattedName = fmt.Sprintf("%s", last)
	case MusicianNamesVariation(FIRSTLAST):
		first, isFirstPresent, _, _, last, _ := m.FullNameTuple()
		if !isFirstPresent {
			first = ""
		}
		formattedName = fmt.Sprintf("%s %s", first, last)
	case MusicianNamesVariation(LASTFIRSTMIDDLE):
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

func (m *Musician) Hash() MusicianHash {
	hashfunc := md5.New()
	// NOTE: assume Musician::String() is unique. Needs assertion, or else expand the Sum() contents
	data := m.PrimaryKey()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return MusicianHash(fmt.Sprintf("%x", hashsum))
}

func NewMusician(data string) (newMusician *Musician, ok bool) {
	newMusician = new(Musician)
	newMusician.Id = MusicianHash(STRING_NULL)
	newMusician.FirstName = STRING_NULL
	newMusician.MiddleName = STRING_NULL
	newMusician.LastName = STRING_NULL
	newMusician.Notes = STRING_NULL
	newMusician.Fields = map[string]string{}
	newMusician.Tags = []string{}
	ok = false

	notes, oknotes, names, okmore := ExtractNotes(data)
	//FailNotOK(okmore, "NewMusician Try to ExctractNotes( FAILED TO FIND NAMES")
	if !okmore {
		return newMusician, false
	}

	if oknotes {
		newMusician.Notes = notes
	}

	firstname, middlename, lastname, ok := ExtractNames(names)
	errors.FailNotOK(ok, "NewMusician try to ExtractNames( FAILED FOR UNKNOWN REASONS")

	newMusician.FirstName = firstname
	newMusician.MiddleName = middlename
	newMusician.LastName = lastname
	newMusician.Id = newMusician.Hash()
	ok = true

	return newMusician, ok
}

func (m *Musician) GetDates(interval uint8) []string {
	return []string{}
}

func (m *Musician) buildField() {

	m.Fields["FIRSTNAME"] = strings.ToUpper(m.FirstName)
	m.Fields["MIDDLENAME"] = strings.ToUpper(m.MiddleName)
	m.Fields["LASTNAME"] = strings.ToUpper(m.LastName)
	//	from Notes with
	//FIELD: TEXT\n all ToUpper
	// plus struct fields

}

func (m *Musician) buildTags() {
	tags := []string{}

	for k, v := range m.Fields {
		tags = strings.Split(strings.ToUpper(k+v), FIELDS_SEP) // TODO replace by SplitFunc or Regex
	}
	m.Tags = tags
}

//var MusicianNULL = Musician{
//	"NULL_HASH",
//	STRING_NULL,
//	STRING_NULL,
//	STRING_NULL,
//	STRING_NULL,
//	//TIME_NULL,
//	//TIME_NULL,
//	//STRING_NULL,
//	//STRING_NULL,
//	//AGE_NULL,
//	//STRING_NULL,
//	// Army string
//	// Rank string
//}

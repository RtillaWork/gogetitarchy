package musician

import (
	"crypto/md5"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"io"
	"strings"
)

func NewMusicianFrom(data string) (newMusician *Musician, ok bool) {
	fname, mname, lname, notes, ok := ExtractNamesNotesFrom(data)
	errors.FailNotOK(ok, "NewMusicianFrom try to ExtractNames( FAILED FOR UNKNOWN REASONS")
	newMusician = New(fname, mname, lname, notes, 1)
	return newMusician, true
}

func New(fname, mname, lname, notes string, encounter uint8) (newMusician *Musician) {
	newMusician = new(Musician)
	*newMusician = Defaults
	newMusician.FName = fname
	newMusician.MName = mname
	newMusician.LName = lname
	newMusician.Notes = notes
	newMusician.Encounter = encounter
	newMusician.Id = newMusician.Hash()

	return newMusician
}

func (m *Musician) PrimaryKey() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("PRIMARYKEY=%s%s%s%s%d", first, middle, last, m.Notes, m.Encounter)
}

func (m *Musician) String() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("%s_%s_%s_%d", first, middle, last, m.Encounter)
}

func (m *Musician) ToCsv() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	id := m.Id
	return fmt.Sprintf("%q; %q; %q; %q; %q, %q", id, first, middle, last, m.Notes, m.Encounter)
}

func (m *Musician) ToJson() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	id := m.Id
	return fmt.Sprintf("{ \"id\": %q,\n \"first_name\": %q,\n \"middle_name\": %q,\n \"last_name\": %q\n}", id, first, middle, last)
}

func (m *Musician) QueryFragment(v NamesVariation) string {
	notes := ""
	if m.Notes != STRING_NULL {
		notes = m.Notes
	}
	return fmt.Sprintf("%s %s", m.NameFmt(v), notes)
}

//
func (m *Musician) FullNameTuple() (fname string, isfnamepresent bool, mname string, ismnamepresent bool, lname string, islnamepresent bool) { //  fname, mname, lname
	fname = Defaults.FName
	mname = Defaults.MName
	lname = Defaults.LName

	isfnamepresent = m.FName != Defaults.FName
	ismnamepresent = m.MName != Defaults.MName
	islnamepresent = m.LName != Defaults.LName
	errors.FailNotOK(islnamepresent, "Musician#FullNameTuple NO LASTNAME")
	lname = m.LName

	if isfnamepresent {
		fname = m.FName
	}

	if ismnamepresent {
		mname = m.MName
	}
	return fname, isfnamepresent, mname, ismnamepresent, lname, islnamepresent
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

func (m *Musician) NameFmt(v NamesVariation) (formattedName string) {
	formattedName = ""
	switch v {
	case NamesVariation(FULL):
		first, isFirstPresent, middle, isMiddlePresent, last, _ := m.FullNameTuple()
		if !isFirstPresent {
			first = ""
		}
		if !isMiddlePresent {
			middle = ""
		}
		formattedName = fmt.Sprintf("%s %s %s", first, middle, last)
	case NamesVariation(L):
		_, _, _, _, last, _ := m.FullNameTuple()
		formattedName = fmt.Sprintf("%s", last)
	case NamesVariation(FL):
		first, isFirstPresent, _, _, last, _ := m.FullNameTuple()
		if !isFirstPresent {
			first = ""
		}
		formattedName = fmt.Sprintf("%s %s", first, last)
	case NamesVariation(LFM):
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

func (m *Musician) GetDates(interval uint8) []string {
	return []string{}
}

func (m *Musician) buildField() {

	m.Fields["FIRSTNAME"] = strings.ToUpper(m.FName)
	m.Fields["MIDDLENAME"] = strings.ToUpper(m.MName)
	m.Fields["LASTNAME"] = strings.ToUpper(m.LName)
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

// OLD

//func NewMusicianFrom(data string) (newMusician *Musician, ok bool) {
//
//	notes, oknotes, names, okmore := ExtractNotes(data)
//	//FailNotOK(okmore, "NewMusicianFrom Try to ExctractNotes( FAILED TO FIND NAMES")
//	if !okmore {
//		return newMusician, false
//	}
//
//	if oknotes {
//		newMusician.Notes = notes
//	}
//
//	fname, middlename, lastname, ok := ExtractNames(names)
//	errors.FailNotOK(ok, "NewMusicianFrom try to ExtractNames( FAILED FOR UNKNOWN REASONS")
//
//	newMusician = New(fname, middlename, lastname, notes, 1)
//	return newMusician, true
//}

package musician

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"io"
	"log"
	"strings"
)

func NewMusicianFrom(data string) (newMusician *Musician, ok bool) {
	fname, mname, lname, notes, ok := ExtractNamesNotesFrom(data)
	if !ok {
		//errors.Assert(ok, "NewMusicianFrom try to ExtractNames( FAILED FOR UNKNOWN REASONS")
		log.Printf("NewMusicianFrom try to ExtractNames( FAILED FOR UNKNOWN REASONS\n")
		return &Defaults, false
	}

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
	return fmt.Sprintf("PRIMARYKEY=%x", m.Hash())
}

func (m *Musician) String() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("%s_%s_%s_%s_%d", m.Id, first, middle, last, m.Encounter)
}

func (m *Musician) ToCsv() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	id := m.Id
	// TODO enumrate .Fields by going through DataDictionary and...
	// TODO ... accumulate key-value orelse key-NOVALUE
	return fmt.Sprintf("%q; %q; %q; %q; %q, %q", id, first, middle, last, m.Notes, m.Encounter)
}

func (m *Musician) ToJson() string {
	jsoned, err := json.Marshal(*m)
	errors.FailOn(err, "Musician::ToJson json.Marshal")
	return fmt.Sprintf("%s", string(jsoned))
}

func (m *Musician) Hash() MusicianHash {
	hashfunc := md5.New()
	// NOTE: assume Musician::String() is unique. Needs assertion, or else expand the Sum() contents
	data := m.ToJson()
	io.WriteString(hashfunc, data)
	hashsum := hashfunc.Sum(nil)
	return MusicianHash(fmt.Sprintf("%x", hashsum))
}

// Support funcs

func (m *Musician) QueryFragment(v NamesVariation) string {
	notes := ""
	// TODO INCLUDE NOTES OR FIELDS OR GoodSetPhrases in Query variations
	//if m.Notes != Defaults.Notes {
	//	notes = m.Notes
	//}
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
	errors.Assert(islnamepresent, "Musician#FullNameTuple NO LASTNAME")
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

// MusiciansDb

func NewMusiciansDb(musicians MusiciansMap) (musiciansdb *MusiciansDb) {
	BuildTheDataDict(musicians)
	musiciansdb = &MusiciansDb{musicians, TheDataDict}

	return musiciansdb
	// TODO repopulate musiciansMap with the same in common KEYS (assigning "" to non existent)
	// TODO build musicians TAGS
	// TODO separate dates and add dates data to fields
}

// MusiciansDb utilities to create Dict and stats

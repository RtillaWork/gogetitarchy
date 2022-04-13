package musician

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"io"
	"log"
	"strings"
	"time"
)

func NewMusicianFrom(rawname string) (newMusician *Musician, ok bool) {
	fname, mname, lname, notes, ok := ExtractNamesNotesFrom(rawname)
	if !ok {
		//errors.Assert(ok, "NewMusicianFrom try to ExtractNames( FAILED FOR UNKNOWN REASONS")
		log.Printf("NewMusicianFrom try to ExtractNames( FAILED FOR UNKNOWN REASONS\n")
		return &Defaults, false
	}

	newMusician = New(rawname, fname, mname, lname, notes)
	newMusician.AddToFields(nil)
	return newMusician, true
}

func New(rawname, fname, mname, lname, notes string) (newMusician *Musician) {
	newMusician = new(Musician)
	*newMusician = Defaults
	newMusician.RawName = utils.NormalizeField(rawname)
	newMusician.FName = utils.NormalizeField(fname)
	newMusician.MName = utils.NormalizeField(mname)
	newMusician.LName = utils.NormalizeField(lname)
	newMusician.Notes = utils.NormalizeField(notes)
	newMusician.TimeCreated = time.Now().UnixNano()
	newMusician.Id = newMusician.Hash()
	log.Printf("FROM NEW %#v musician\n", newMusician)
	return newMusician
}

func (m *Musician) PrimaryKey() string {
	return fmt.Sprintf("PRIMARYKEY=%x", m.Hash())
}

func (m *Musician) String() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	return fmt.Sprintf("%s_%s_%s_%s_%d", m.Id, first, middle, last, m.TimeCreated)
}

func (m *Musician) ToCsv() string {
	first, _, middle, _, last, _ := m.FullNameTuple()
	id := m.Id
	// TODO enumrate .Fields by going through DataDictionary and...
	// TODO ... accumulate key-value orelse key-NOVALUE
	return fmt.Sprintf("%q; %q; %q; %q; %q, %q", id, first, middle, last, m.Notes, m.TimeCreated)
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

func (m *Musician) AddToFields(fields map[string]string) {
	if fields == nil {
		m.Fields["FIRSTNAME"] = utils.NormalizeValue(m.FName)
		m.Fields["MIDDLENAME"] = utils.NormalizeValue(m.FName)
		m.Fields["LASTNAME"] = utils.NormalizeValue(m.FName)
	} else {
		for k, v := range fields {
			m.Fields[utils.NormalizeKey(k)] = utils.NormalizeValue(v)
		}
	}

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
	musiciansdb = &MusiciansDb{&musicians, &TheDataDict}

	return musiciansdb
	// TODO repopulate musiciansMap with the same in common KEYS (assigning "" to non existent)
	// TODO build musicians TAGS
	// TODO separate dates and add dates data to fields
}

// MusiciansDb utilities to create Dict and stats

// Utilities: Contains series of func

//
func (ms MusiciansMap) CountRawName(mayberawname string) (count int) {
	if len(ms) == 0 {
		return -1
	}
	count = 0
	maybename := utils.NormalizeField(mayberawname)
	for _, m := range ms {
		if maybename == m.RawName {
			count++
		}
	}
	return count
}

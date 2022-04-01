package musician

import (
	"bufio"
	"github.com/RtillaWork/gogetitarchy/utils"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"log"
	"os"
	"strings"
)

// BlockDelim Some interesting block elements contain `:` as fields separators
const BlockDelim = "Civil War (Union)" // must be the second line, following the soldier's name
var skipThese = []string{BlockDelim, "MEMORIAL", ""}

type RawMusicianBlock struct {
	Names           string            `json:"Names"`
	ConnectionCount string            `json:"connection_count"`
	Notes           map[string]string `json:"notes"`
}

type MusiciansMap map[MusicianHash]*Musician

func ImportData(inFileName string, delim string) MusiciansMap {

	inFile, err := os.Open(inFileName)
	errors.FailOn(err, "opening inFile for reading...")
	defer inFile.Close()

	musicians := make(MusiciansMap)

	s := bufio.NewScanner(inFile)

	//func initBlock() {
	//	lines = []string{}
	//	lines = append(lines, prevline) // first string would always be name (preceding delim)
	//}
	lines := []string{}
	for state, prevline := 0, ""; s.Scan(); {
		line := s.Text()
		log.Printf("s.Text() %s\n", line)
		switch state {
		case 0:
			{
				if line == delim {
					lines = append([]string{}, prevline) // first string would always be name (preceding delim)
					state = 1
					continue

				}

			}

		case 1:
			{
				if line == delim {
					if len(lines[0]) != 0 {
						musician, ok := ReadMusicianData(lines)
						if ok {
							musicians[musician.Id] = musician
						}
					} else {
						log.Printf("\n = = ERROR READING FOR FILE: line:{ %v } prevline:{ %v}\n\n", line, prevline)
					}
					lines = append([]string{}, prevline) // first string would always be name (preceding delim)
					continue

				} else {
					if utils.IsLikelyValidData(strings.TrimSpace(prevline), skipThese) {
						lines = append(lines, prevline)
					}
				}

			}

		}

		prevline = line
		//log.Printf("prevline %s\n", prevline)
		//utils.WaitForKeypress()
	}
	return musicians

}

func ReadMusicianData(ablock []string) (musician *Musician, ok bool) {
	log.Printf("### ablock[0] %s\n", ablock[0])
	//utils.WaitForKeypress()
	musician, ok = NewMusician(ablock[0])
	if !ok {
		return musician, false
	}
	//errors.FailNotOK(ok, "\n\nSCANNING BAD line: %s ONLT FOUND NOTES, NO NAMES\n\n")
	musician.Id = musician.Hash()
	ExtractFields(ablock)
	//log.Printf("\nSCANNING SUCCESS aMusican: {  %v  }\n\n", aMusician.Hash())

	return musician, true

}

func ReadMusiciansNames(inFileName string) MusiciansMap {

	inFile, err := os.Open(inFileName)
	errors.FailOn(err, "opening inFile for reading...")
	defer inFile.Close()

	musicians := make(MusiciansMap)

	s := bufio.NewScanner(inFile)
	for line := ""; s.Scan(); {
		line = s.Text()
		//log.Printf("SCANNING line: %s\n", line)
		aMusician, ok := NewMusician(line)
		if !ok {
			continue
			log.Printf("\n\nSCANNING BAD line: %s\n\n", line)
		}

		musicians[aMusician.Hash()] = aMusician
		//log.Printf("\nSCANNING SUCCESS aMusican: {  %v  }\n\n", aMusician.Hash())
	}
	return musicians

}

// TODO: replace by regex

func ExtractNotes(data string) (notes string, foundnotes bool, rest string, foundmore bool) {
	// returns (notes if any, truncated data without notes if any, oknotes, okmore)
	//notes := STRING_NULL
	//foundnotes := false
	//rest := data
	//foundmore := false

	notes = STRING_NULL
	foundnotes = false
	rest = data
	foundmore = false

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

func ExtractNames(data string) (firstname string, middlename string, lastname string, ok bool) {
	// returns (firstname if any, middlename if any, lastname if any, ok)
	//firstname := STRING_NULL
	//middlename := STRING_NULL
	//lastname := STRING_NULL
	//ok := false
	////rest := data

	firstname = STRING_NULL
	middlename = STRING_NULL
	lastname = STRING_NULL
	ok = false
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

func ExtractFields(data []string) (fields map[string]string) {
	fields = make(map[string]string)
	for i, d := range data {
		//k := strings.Split(d, ":")
		v := strings.Split(d, ":")[0]
		fields[string(i)] = v
	}
	return fields
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

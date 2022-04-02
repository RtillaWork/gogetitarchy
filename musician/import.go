package musician

import (
	"bufio"
	"github.com/RtillaWork/gogetitarchy/utils"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"log"
	"os"
	"regexp"
	"strings"
)

// BlockDelimDef Some interesting block elements contain `:` as fields separators
const BlockDelimDef = "Civil War (Union)" // must be the second line, following the soldier's name
const block_FIELD_SEP = ":"
const block_DATE_SEP = "-"

var skipThese = []string{BlockDelimDef, "MEMORIAL", ""}

// ImportData builds a MusiciansMap from a textfile where names section precedes a delimiter
// it reads the musician block content (partially unstructured)
func ImportData(inFileName string, delim string) (musicians MusiciansMap) {

	inFile, err := os.Open(inFileName)
	errors.FailOn(err, "opening inFile for reading...")
	defer inFile.Close()

	s := bufio.NewScanner(inFile)

	blklines := []string{}
	for initial, curln, prevln := true, "", ""; s.Scan(); prevln = curln {
		curln = s.Text()
		// NOTE DEBUG
		log.Printf("for prevline %s\n", prevln)
		log.Printf("for curln %s\n", curln)
		log.Printf("blklines %#v\n", blklines)
		log.Printf("initial %#v\n", initial)
		// END NOTE DEBUG

		if initial && curln == delim {
			initial = false
			blklines[0] = prevln // prevlin == names
			log.Printf("if initial blklines %#v\n", blklines)
			curln = prevln // to skip the next coniditon during the transition from initial true to false
		}

		if !initial && curln == delim {
			amusician, ok := ReadMusicianData(blklines)
			if ok {
				musicians[amusician.Id] = amusician
				log.Printf("ENTRY ADDED to RawMusicians \n")

			} else {
				log.Printf("ENTRY IGNORED UNDERTERMINATE REASON \n")
				log.Printf("\n = = ERROR READING FOR FILE: line:{ %v } prevline:{ %v}\n\n", curln, prevln)

			}
			blklines = []string{}
			blklines[0] = prevln // prevlin == names // prevlin == names
			log.Printf("if not initial blklines %#v\n", blklines)
		}
		blklines = append(blklines, prevln)
		utils.WaitForKeypress()

	}

	return musicians

}

// ReadMusicianData creates a Musician struct data from a partially unstructured block of []string
// it expects that block[0] is t least present with names
func ReadMusicianData(ablock []string) (amusician *Musician, ok bool) {
	errors.Assert(len(ablock) != 0, "ReadMusicianData []ablock is nil or empty\n")
	log.Printf("### ablock[0] %#v\n", ablock[0])
	utils.WaitForKeypress()
	amusician, ok = NewMusicianFrom(ablock[0])
	if !ok {
		return amusician, false
	}
	//errors.Assert(ok, "\n\nSCANNING BAD line: %s ONLT FOUND NOTES, NO NAMES\n\n")
	amusician.Id = amusician.Hash()
	if len(ablock) > 1 {
		amusician.Fields = ExtractFields(ablock[1:])
	}
	//log.Printf("\nSCANNING SUCCESS aMusican: {  %v  }\n\n", aMusician.Hash())
	return amusician, true
}

// ExtractNamesNotesFrom returns (firstname if any, middlename if any, lastname if any, ok if lastname)
// L, F  || F M L || F M. L || F L || F "M" L (NOTES)
func ExtractNamesNotesFrom(data string) (fname string, mname string, lname string, notes string, ok bool) {
	errors.Assert(len(data) != 0, "ExtractFrom data is empty")
	fname, mname, lname, notes = Defaults.FName, Defaults.MName, Defaults.LName, Defaults.Notes
	ok = false
	// split names away from notes through `(`, if exists
	names, notes := "", Defaults.Notes
	switch s := strings.Split(strings.TrimSpace(data), NOTES_SEP_OPEN); len(s) {
	case 0:
		errors.Assert(false, "ExtractFrom switch Split error data likely nil/empty")
	case 1:
		if strings.Contains(s[0], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
			errors.Assert(false, "ExtractFrom Contains error data likely conmtains only notes but no names")
		} else {
			names = s[0]
		}
	case 2:
		names = strings.TrimSpace(s[0])
		notes = strings.TrimSpace(strings.Trim(s[1], NOTES_SEP_OPEN+NOTES_SEP_CLOSE))
	default:
		errors.Assert(false, "ExtractFrom data Split returned too many fields separated by `(`")
	}

	// Only one will match. len(result) = number of names + 1
	r0 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)[\W\s]*$`)                               // L
	r1 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+),\s*([A-Za-z]+)[\W\s]*$`)                // L, F
	r2 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+([A-Za-z]+)[\W\s]*$`)                 // F L
	r3 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+([A-Za-z]+)\s+([A-Za-z]+)[\W\s]*$`)   // F M L
	r4 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+([A-Za-z]\.)\s+([A-Za-z]+)[\W\s]*$`)  // F M. L
	r5 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+("[A-Za-z]+")\s+([A-Za-z]+)[\W\s]*$`) // F "M" L

	s0 := r0.FindAllStringSubmatch(names, -1)
	s1 := r1.FindAllStringSubmatch(names, -1)
	s2 := r2.FindAllStringSubmatch(names, -1)
	s3 := r3.FindAllStringSubmatch(names, -1)
	s4 := r4.FindAllStringSubmatch(names, -1)
	s5 := r5.FindAllStringSubmatch(names, -1)
	switch {
	case len(s0[0]) == 2:
		lname = s0[0][1]
		mname = Defaults.MName
		fname = Defaults.FName

	case len(s1[0]) == 3:
		lname = s0[0][1]
		mname = Defaults.MName
		fname = s0[0][2]

	case len(s2[0]) == 3:
		lname = s0[0][2]
		mname = Defaults.MName
		fname = s0[0][1]

	case len(s3[0]) == 4:
		lname = s0[0][3]
		mname = s0[0][2]
		fname = s0[0][1]

	case len(s4[0]) == 4:
		lname = s0[0][3]
		mname = s0[0][2]
		fname = s0[0][1]

	case len(s5[0]) == 4:
		lname = s0[0][3]
		mname = s0[0][2]
		fname = s0[0][1]

	default:
		// Errors
		lname = Defaults.LName
		mname = Defaults.MName
		fname = Defaults.FName
	}

	if len(lname) > 1 {
		ok = true
	} else {
		ok = false
	}
	return fname, mname, lname, notes, ok
}

// ExtractFields tries to build Musician.Fields map from valid data in the previously read block []string
func ExtractFields(data []string) (fields map[string]string) {
	fields = make(map[string]string)
	//log.Printf("Raw Block Data i:{ %v }\n %s\n", data, data)
	for i, d := range data {
		// NOTE DEBUG
		log.Printf("### ablock[%d] %v\n", i, d)
		// NOTE END DEBUG

		if !utils.IsLikelyValidData(d, skipThese) {
			//continue
			log.Printf("DEBUG IS UNLIKELY VALID DATA")
		}

		s := strings.Split(strings.TrimSpace(d), block_FIELD_SEP)
		// NOTE DEBUG
		log.Printf("### s[%d] %v\n", i, s)
		// NOTE END DEBUG
		if len(s) == 0 {
			continue
		} else if len(s[0]) == 0 {
			s = s[1:]
		} else {
			s = s[0:]
		}
		var k, v string
		switch l := len(s); l {
		case 0:
			continue
		case 1:
			k = strings.ToUpper(s[0])
			v = s[0]
		case 2:
			k = strings.ToUpper(s[0])
			v = s[1]
		default:
			k = strings.ToUpper(s[0])
			v = strings.Join(s[1:], block_FIELD_SEP)
		}

		fields[k] = v
		//log.Printf("BLOCK i: %v { %v }\n %s\n", i, fields, fields)
	}
	utils.WaitForKeypress()
	// NOTE DEBUG
	for k, v := range fields {
		log.Printf("BLOCK: k:  { %v } v:   %s\n", k, v)
	}
	// END NOTE DEBUG
	return fields
}

// OLDER functions

//func ReadMusiciansNames(inFileName string) MusiciansMap {
//
//	inFile, err := os.Open(inFileName)
//	errors.FailOn(err, "opening inFile for reading...")
//	defer inFile.Close()
//
//	musicians := make(MusiciansMap)
//
//	s := bufio.NewScanner(inFile)
//	for line := ""; s.Scan(); {
//		line = s.Text()
//		//log.Printf("SCANNING line: %s\n", line)
//		aMusician, ok := NewMusicianFrom(line)
//		if !ok {
//			continue
//			log.Printf("\n\nSCANNING BAD line: %s\n\n", line)
//		}
//
//		musicians[aMusician.Hash()] = aMusician
//		//log.Printf("\nSCANNING SUCCESS aMusican: {  %v  }\n\n", aMusician.Hash())
//	}
//	return musicians
//}

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

////  old

//
//func ImportData(inFileName string, delim string) (*MusiciansMap) {
//
//	inFile, err := os.Open(inFileName)
//	errors.FailOn(err, "opening inFile for reading...")
//	defer inFile.Close()
//
//	s := bufio.NewScanner(inFile)
//
//	//func initBlock() {
//	//	lines = []string{}
//	//	lines = append(lines, prevline) // first string would always be name (preceding delim)
//	//}
//
//	lines := []string{}
//	for state, prevline := 0, ""; s.Scan(); {
//		line := s.Text()
//
//		// NOTE DEBUG
//		log.Printf("prevline %s\n", prevline)
//		log.Printf("s.Text() %s\n", line)
//		// END NOTE DEBUG
//		switch state {
//		case 0:
//			{
//				if line == delim {
//					lines = append([]string{}, prevline) // first string would always be name (preceding delim)
//					state = 1
//					continue
//
//				}
//
//			}
//
//		case 1:
//			{
//				if line == delim {
//					if len(lines[0]) != 0 {
//						musician, ok := ReadMusicianData(lines)
//						if ok {
//							musicians[musician.Id] = musician
//						}
//					} else {
//						log.Printf("\n = = ERROR READING FOR FILE: line:{ %v } prevline:{ %v}\n\n", line, prevline)
//					}
//					lines = append([]string{}, prevline) // first string would always be name (preceding delim)
//					continue
//
//				} else {
//					lines = append(lines, prevline)
//
//				}
//
//			}
//
//		}
//
//		prevline = line
//
//		//utils.WaitForKeypress()
//	}
//	return MusiciansMap{}
//
//}

//
//func ReadMusicianData(ablock []string) (musician *Musician, ok bool) {
//	//log.Printf("### ablock[0] %s\n", ablock[0])
//
//	//utils.WaitForKeypress()
//	musician, ok = NewMusicianFrom(ablock[0])
//	if !ok {
//		return musician, false
//	}
//	//errors.Assert(ok, "\n\nSCANNING BAD line: %s ONLT FOUND NOTES, NO NAMES\n\n")
//	musician.Id = musician.Hash()
//	if len(ablock) > 1 {
//		ExtractFields(ablock[1:])
//	}
//	//log.Printf("\nSCANNING SUCCESS aMusican: {  %v  }\n\n", aMusician.Hash())
//
//	return musician, true
//
//}

// OLD

//// TODO: replace by regex 02APR2222 INPROGRESS
//
//func ExtractNotes(data string) (notes string, foundnotes bool, rest string, foundmore bool) {
//	// returns (notes if any, truncated data without notes if any, oknotes, okmore)
//	//notes := STRING_NULL
//	//foundnotes := false
//	//rest := data
//	//foundmore := false
//
//	notes = STRING_NULL
//	foundnotes = false
//	rest = data
//	foundmore = false
//
//	fields := strings.Fields(data)
//	switch len(fields) {
//	case 0:
//		//return notes,  false, rest, false
//	case 1:
//		// NOTE: assume if only one field
//		if strings.HasPrefix(fields[0], NOTES_SEP_OPEN) && strings.HasSuffix(fields[0], NOTES_SEP_CLOSE) {
//			// CASE (ABC) -> notes
//			// NOTE: should panic as it means no names in the scanned line
//			notes = fields[0]
//			//return notes, false, rest,  false
//		} else if !strings.ContainsAny(fields[0], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
//			// CASE ABC -> lastname
//			rest = fields[0]
//			foundmore = true
//			//return notes, false,  rest, true
//		} else {
//			// CASE badly formed ex (ABC\n or \ABC)
//			//return notes, false, data, false
//		}
//
//	case 2:
//		//
//		if strings.HasPrefix(fields[1], NOTES_SEP_OPEN) && strings.HasSuffix(fields[1], NOTES_SEP_CLOSE) {
//			// CASE XYZ (ABC) -> notes
//			// NOTE: should panic as it means no names in the scanned line
//			notes = fields[1]
//			foundnotes = true
//			rest = fields[0]
//			foundmore = true
//			//return notes, true, rest, true
//		} else if !strings.ContainsAny(fields[1], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
//			// CASE XYZ ABC  -> rest
//			foundmore = true
//			//return notes, false, rest,  true
//		} else {
//			// CASE badly formed
//			//return  notes,  false, data,false
//		}
//
//	case 3:
//		//
//		if strings.HasPrefix(fields[2], NOTES_SEP_OPEN) && strings.HasSuffix(fields[2], NOTES_SEP_CLOSE) {
//			// CASE X YZ (ABC) -> notes
//			// NOTE: should panic as it means no names in the scanned line
//			notes = fields[2]
//			foundnotes = true
//			rest = fields[0] + NAMES_DEFAULT_SEP + fields[1]
//			foundmore = true
//			//return notes,  true, rest, true
//		} else if !strings.ContainsAny(fields[2], NOTES_SEP_OPEN+NOTES_SEP_CLOSE) {
//			// CASE X YZ ABC  -> rest
//			rest = fields[0] + NAMES_DEFAULT_SEP + fields[1] + NAMES_DEFAULT_SEP + fields[2]
//			foundmore = true
//			//return notes, false, rest, true
//		} else {
//			// CASE badly formed
//			//return  notes, false, data, false
//		}
//
//	case 4:
//		//
//		if strings.HasPrefix(fields[3], NOTES_SEP_OPEN) && strings.HasSuffix(fields[3], NOTES_SEP_CLOSE) {
//			// CASE X Y Z (ABC) -> notes
//			// NOTE: should panic as it means no names in the scanned line
//			notes = fields[3]
//			foundnotes = true
//			rest = fields[0] + NAMES_DEFAULT_SEP + fields[1] + NAMES_DEFAULT_SEP + fields[2]
//			foundmore = true
//			//return notes,true,  rest, true
//		} else {
//			// CASE badly formed X Y Z ABC  -> ??
//			//return  notes, false, data, false
//		}
//
//	default:
//
//	}
//	return notes, foundnotes, rest, foundmore
//}
//
//func ExtractNames(data string) (firstname string, middlename string, lastname string, ok bool) {
//	// returns (firstname if any, middlename if any, lastname if any, ok)
//	//firstname := STRING_NULL
//	//middlename := STRING_NULL
//	//lastname := STRING_NULL
//	//ok := false
//	////rest := data
//
//	firstname = STRING_NULL
//	middlename = STRING_NULL
//	lastname = STRING_NULL
//	ok = false
//	//rest := data
//
//	fields := strings.Fields(data)
//	switch len(fields) {
//	case 0:
//		//return _, _ ,_ , false
//	case 1:
//		// NOTE: assume if only one field
//		if strings.HasSuffix(fields[0], LAST_NAME_SEP) {
//			// CASE ABC,
//			// NOTE: should panic as it means no names in the scanned line
//			lastname = strings.TrimSuffix(fields[0], LAST_NAME_SEP)
//			ok = true
//
//		} else if !strings.ContainsAny(fields[0], LAST_NAME_SEP+NAMES_DEFAULT_SEP) {
//			// CASE ABC -> lastname
//			lastname = fields[0]
//			ok = true
//
//		} else {
//			// CASE badly formed
//
//		}
//
//	case 2:
//		//
//		if strings.HasSuffix(fields[0], LAST_NAME_SEP) {
//			// CASE XYZ, ABC
//			lastname = strings.TrimSuffix(fields[0], LAST_NAME_SEP)
//			firstname = fields[1]
//			ok = true
//
//		} else if !strings.ContainsAny(fields[0], LAST_NAME_SEP+NAMES_DEFAULT_SEP) {
//			// CASE XYZ ABC
//			firstname = fields[0]
//			lastname = fields[1]
//			ok = true
//		} else {
//			// CASE badly formed
//
//		}
//
//	case 3:
//		//
//		if strings.HasSuffix(fields[0], LAST_NAME_SEP) {
//			// CASE XYZ, ABC I.
//			lastname = strings.TrimSuffix(fields[0], LAST_NAME_SEP)
//			firstname = fields[1]
//			middlename = fields[2]
//			ok = true
//
//		} else if !strings.ContainsAny(fields[0], LAST_NAME_SEP+NAMES_DEFAULT_SEP) {
//			// CASE XYZ ABC
//			firstname = fields[0]
//			middlename = fields[1]
//			lastname = fields[2]
//			ok = true
//		} else {
//			// CASE badly formed
//
//		}
//
//	default:
//
//	}
//	return firstname, middlename, lastname, ok
//}

//// TESTS
//
//https://go.dev/play/p/caqAFfFJLWq
//package main
//
//import (
//"fmt"
//"regexp"
//)
//
//func main() {
//r0 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)[\W\s]*$`)                                 // L
//r1 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+),\s*([A-Za-z]+)[\W\s]*$`)                  // L, F
//r2 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+([A-Za-z]+)[\W\s]*$`)                   // F L
//r3 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+([A-Za-z]+)\s+([A-Za-z]+)[\W\s]*$`)     // F M L
//r4 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+([A-Za-z]\.)\s+([A-Za-z]+)[\W\s]*$`)    // F M. L
//r5 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)\s+(\"[A-Za-z]+\")\s+([A-Za-z]+)[\W\s]*$`) // F "M" L
//
////n0 := "  Last     "
////n1 := "  Last,    First"
////n2 := ",   First    Last"
////n3 := "First   Middle   Last"
////n4 := "First M. Last"
//n5 := "First \"Moo\" Last"
//
//fmt.Printf("0 %#v\n", r0.FindAllStringSubmatch(n5, -1))
//fmt.Printf("1 %#v\n", r1.FindAllStringSubmatch(n5, -1))
//fmt.Printf("2 %#v\n", r2.FindAllStringSubmatch(n5, -1))
//fmt.Printf("3 %#v\n", r3.FindAllStringSubmatch(n5, -1))
//fmt.Printf("4 %#v\n", r4.FindAllStringSubmatch(n5, -1))
//fmt.Printf("5 %#v\n", r5.FindAllStringSubmatch(n5, -1))
//
//}

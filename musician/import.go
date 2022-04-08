package musician

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"log"
	"os"
	"regexp"
	"strings"
)

// BlockDelimDef1 Some interesting block elements contain `:` as fields separators
const BlockDelimDef1 = "Civil War (Union)"       // must be the second line, following the soldier's name
const BlockDelimDef2 = "Civil War (Confederate)" // must be the second line, following the soldier's name
const block_FIELD_SEP = ":"
const block_DATE_SEP = "-"

var skipThese = []string{BlockDelimDef1, BlockDelimDef2, "MEMORIAL", ""}

// ImportData builds a MusiciansMap from a textfile where names section precedes a delimiter
// it reads the musician block content (partially unstructured)
func ImportData(inFileName string, delim1 string, delim2 string) (musicians MusiciansMap) {
	totalcount := 0
	musicians = make(MusiciansMap)

	inFile, err := os.Open(inFileName)
	errors.FailOn(err, "opening inFile for reading...")
	defer inFile.Close()

	s := bufio.NewScanner(inFile)

	blklines := []string{}
	for initial, curln, prevln := true, "", ""; s.Scan(); prevln = curln {
		curln = s.Text()
		//// NOTE DEBUG
		//log.Printf("for prevline %s\n", prevln)
		//log.Printf("for curln %s\n", curln)
		//log.Printf("blklines %#v\n", blklines)
		////log.Printf("initial %#v\n", initial)
		//// END NOTE DEBUG

		if initial && (curln == delim1 || curln == delim2) {
			initial = false
			blklines[0] = prevln // prevlin == names
			//log.Printf("if initial blklines %#v\n", blklines)
			continue // to skip the next coniditon during the transition from initial true to false
		}

		if !initial && (curln == delim1 || curln == delim2) {
			amusician, ok := ReadMusicianData(blklines)
			if ok {
				musicians[amusician.Id] = amusician
				totalcount++
				log.Printf("Musician ENTRY count %d ADDED to RawMusicians %v \n\n", totalcount, amusician.ToJson())

			} else {
				log.Printf("ENTRY %v IGNORED UNDERTERMINATE REASON \n", amusician.ToJson())
				log.Printf("\n = = ERROR READING FOR FILE: line:{ %v } prevline:{ %v}\n\n", curln, prevln)
				utils.WaitForKeypress()

			}
			blklines = []string{}
			//log.Printf("if not initial   prevline %s\n", prevln)
			//log.Printf("if not initial   curln %s\n", curln)
			//log.Printf("if not initial  blklines %#v\n", blklines)
			blklines = []string{prevln} // prevlin == names
			// TODO DELETEME blklines[0] = prevln // prevlin == names
			log.Printf("if not initial blklines after %#v\n", blklines)
		}
		blklines = append(blklines, prevln)
		//utils.WaitForKeypress()

	}

	log.Printf("\nTotalCount %d = musicians.len %d", totalcount, len(musicians))
	utils.WaitForKeypress()
	return musicians

}

// ReadMusicianData creates a Musician struct data from a partially unstructured block of []string
// it expects that block[0] is at least present with names
func ReadMusicianData(ablock []string) (amusician *Musician, ok bool) {
	errors.Assert(len(ablock) != 0, "ReadMusicianData []ablock is nil or empty\n")
	//log.Printf("### ablock[0] %#v\n", ablock[0])
	//utils.WaitForKeypress()
	amusician, ok = NewMusicianFrom(ablock[0])
	if !ok {
		return amusician, false
	}
	//errors.Assert(ok, "\n\nSCANNING BAD line: %s ONLT FOUND NOTES, NO NAMES\n\n")
	amusician.Id = amusician.Hash()
	if len(ablock) > 1 {
		amusician.AddFields(ExtractFields(ablock[1:]))
	}
	//log.Printf("\nSCANNING SUCCESS aMusican: {  %v  }\n\n", aMusician.Hash())
	return amusician, true
}

// ExtractNamesNotesFrom returns (firstname if any, middlename if any, lastname if any, ok if lastname)
// L, F  || F M L || F M. L || F L || F "M" L (NOTES)
func ExtractNamesNotesFrom(data string) (fname string, mname string, lname string, notes string, ok bool) {

	if len(data) == 0 {
		//errors.Assert(len(data) != 0, "ExtractNamesNotesFrom data is empty")
		log.Printf("ExtractNamesNotesFrom data is empty. returning defaults and false\n")
		return Defaults.FName, Defaults.MName, Defaults.LName, Defaults.Notes, false
	}

	fname, mname, lname, notes, ok = Defaults.FName, Defaults.MName, Defaults.LName, Defaults.Notes, false

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
	r6 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+),\s*([A-Za-z]+)\s+([A-Za-z])[\W\s]*$`)   // L, F M
	r7 := regexp.MustCompile(`(?is)^[\W\s]*([A-Za-z]+)[\W\s]*`)                                // L , gar

	s0 := r0.FindAllStringSubmatch(names, -1)
	s1 := r1.FindAllStringSubmatch(names, -1)
	s2 := r2.FindAllStringSubmatch(names, -1)
	s3 := r3.FindAllStringSubmatch(names, -1)
	s4 := r4.FindAllStringSubmatch(names, -1)
	s5 := r5.FindAllStringSubmatch(names, -1)
	s6 := r6.FindAllStringSubmatch(names, -1)
	s7 := r7.FindAllStringSubmatch(names, -1)
	switch {
	case len(s0) > 0 && len(s0[0]) == 2:
		lname = s0[0][1]
		mname = Defaults.MName
		fname = Defaults.FName
		ok = true
	case len(s1) > 0 && len(s1[0]) == 3:
		lname = s1[0][1]
		mname = Defaults.MName
		fname = s1[0][2]
		ok = true
	case len(s2) > 0 && len(s2[0]) == 3:
		lname = s2[0][2]
		mname = Defaults.MName
		fname = s2[0][1]
		ok = true
	case len(s3) > 0 && len(s3[0]) == 4:
		lname = s3[0][3]
		mname = s3[0][2]
		fname = s3[0][1]
		ok = true
	case len(s4) > 0 && len(s4[0]) == 4:
		lname = s4[0][3]
		mname = s4[0][2]
		fname = s4[0][1]
		ok = true
	case len(s5) > 0 && len(s5[0]) == 4:
		lname = s5[0][3]
		mname = s5[0][2]
		fname = s5[0][1]
		ok = true
	case len(s6) > 0 && len(s6[0]) == 4:
		lname = s6[0][1]
		mname = s6[0][3]
		fname = s6[0][2]
		ok = true
	case len(s7) > 0 && len(s7[0]) == 2:
		lname = s7[0][1]
		mname = Defaults.MName
		fname = Defaults.FName
		ok = true
	default:
		// Errors
		log.Printf("####### WARNING UNDEFINED REGREX FOR names \n")
		log.Printf("REGEX s 0 %#v\n", s0)
		log.Printf("REGEX s 1 %#v\n", s1)
		log.Printf("REGEX s 2 %#v\n", s2)
		log.Printf("REGEX s 3 %#v\n", s3)
		log.Printf("REGEX s 4 %#v\n", s4)
		log.Printf("REGEX s 5 %#v\n", s5)
		log.Printf("REGEX s 6 %#v\n", s6)
		errors.Assert(false, "UNDEFINED REGREX FOR names")
		lname = Defaults.LName
		mname = Defaults.MName
		fname = Defaults.FName
		ok = false
	}

	//if len(lname) > 1 {
	//	ok = true
	//} else {
	//	ok = false
	//}
	return fname, mname, lname, notes, ok
}

// ExtractFields tries to build Musician.Fields map from valid data in the previously read block []string
func ExtractFields(data []string) (fields map[string]string) {
	fields = make(map[string]string)
	datecsv := ""
	//log.Printf("Raw Block Data i:{ %v }\n %s\n", data, data)
	for _, d := range data { // _=i
		// NOTE DEBUG
		//log.Printf("### \n\nablock[%d] %v\n", i, d)
		// NOTE END DEBUG

		if utils.IsUnwantedInput(d, skipThese) {
			continue
			//log.Printf("DEBUG IS UNLIKELY VALID DATA")
		}

		//
		k, v, ok := breakLineByFields(d, block_FIELD_SEP)
		if !ok {
			continue
		} else if k != v {
			v1, v2, ok := breakLineByFields(v, block_DATE_SEP)
			if ok && v1 != v2 {
				datecsv = fmt.Sprintf("%s:%s%s%s%s%s", k, v1, CSV_SEP, k, v2, CSV_SEP)
			}

		}

		//s := strings.Split(strings.TrimSpace(d), block_FIELD_SEP)
		//// NOTE DEBUG
		////log.Printf("### s[%d] %v\n", i, s)
		//// NOTE END DEBUG
		//if len(s) == 0 {
		//	continue
		//} else if len(s[0]) == 0 {
		//	s = s[1:]
		//} else {
		//	s = s[0:]
		//}
		//var k, v string
		//switch l := len(s); l {
		//case 0:
		//	continue
		//case 1:
		//	k = "MISCELLANEA"
		//	v += utils.NormalizeValue(s[0]) //s[0]
		//	// PREVIOUSLY k = strings.ToUpper(s[0])
		//	// PREVIOUSLY v = s[0]
		//case 2:
		//	k = utils.NormalizeKey(s[0])   // strings.ToUpper(s[0])
		//	v = utils.NormalizeValue(s[1]) // s[1]
		//default:
		//	k = utils.NormalizeKey((s[0])) // strings.ToUpper(s[0])
		//	v = utils.NormalizeValue(strings.Join(s[1:], block_FIELD_SEP))
		//}

		fields[k] = v
		fields["DATECSV"] = datecsv
		//log.Printf("BLOCK i: %v { %v }\n %s\n", i, fields, fields)

	}
	//utils.WaitForKeypress()
	//// NOTE DEBUG
	//for k, v := range fields {
	//	log.Printf("BLOCK: k:  { %v } v:   %s\n", k, v)
	//}
	//// END NOTE DEBUG
	return fields
}

// Utility funcs

// break raw line by meaninful field separators, successsively/recursively
func breakLineByFields(d string, sep string) (key string, value string, ok bool) {
	s := strings.Split(strings.TrimSpace(d), sep)
	// NOTE DEBUG
	//log.Printf("### s[%d] %v\n", i, s)
	// NOTE END DEBUG
	if len(s) == 0 {
		return "", "", false
	} else if len(s[0]) == 0 {
		s = s[1:]
	} else {
		s = s[0:]
	}
	var k, v string
	switch l := len(s); l {
	case 0:
		return "", "", false
	case 1:
		k = "MISCELLANEA"
		v += (utils.NormalizeValue(s[0]) + block_FIELD_SEP) //s[0]
		// PREVIOUSLY k = strings.ToUpper(s[0])
		// PREVIOUSLY v = s[0]
	case 2:
		k = utils.NormalizeKey(s[0])   // strings.ToUpper(s[0])
		v = utils.NormalizeValue(s[1]) // s[1]
	default:
		k = utils.NormalizeKey((s[0]))                         // strings.ToUpper(s[0])
		v = utils.NormalizeValue(strings.Join(s[1:], CSV_SEP)) //  block_FIELD_SEP
	}

	key = k
	value = v
	ok = true
	return key, value, ok
}

// ReadData
func ReadData(data []byte) (musicians MusiciansMap) {
	musicians = make(MusiciansMap)

	musicianslist := make([]Musician, 0)

	err := json.Unmarshal(data, &musicianslist)
	errors.FailOn(err, "ReadData Import Musicians")

	for _, m := range musicianslist {
		musicians[m.Id] = &m
	}
	return musicians
}

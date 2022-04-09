package utils

import (
	"bufio"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"log"
	"os"
	"regexp"
	"strings"
)

type ValidScore int8

const (
	False ValidScore = -127
	True  ValidScore = 127
)

const xNAMES_DEFAULT_SEP = " "
const xLAST_NAME_SEP = ","
const xINITIALS_SEP = ". " // I.N.I.T._NAMES
const xNOTES_SEP_OPEN = "("
const xNOTES_SEP_CLOSE = ")"
const xFIELDS_SEP = ",: -"
const xCSV_SEP = ";"

var rWord = regexp.MustCompile(`[A-Zaz]+`)
var rBlank = regexp.MustCompile(`^[\s]+$`)

// regexp improvement and generalization of data validation
func IsUnwantedInput(s string, badset []string) (ok bool) {
	// tests if s is blank
	if rBlank.MatchString(s) {
		return true
	} else {
		return false
	}

	// test if is in badset
	// TODO replace with regex template
	ok = false
	for _, badstr := range badset {
		if strings.TrimSpace(strings.ToUpper(s)) == strings.TrimSpace(strings.ToUpper(badstr)) {
			ok = true
			break
		}
	}

	return ok
}

// regexp improvement and generalization of data validation
func IsWanteddInput(s string, goodset []string) (ok bool) {
	//// tests if s is blank
	//if rBlank.MatchString(s) {
	//	return true
	//} else {
	//	return false
	//}
	//
	//// test if is in badset
	//ok = false
	//for _, badstr := range badset {
	//	if strings.TrimSpace(strings.ToUpper(s)) == strings.TrimSpace(strings.ToUpper(badstr)) {
	//		ok = true
	//		break
	//	}
	//}

	//return ok
	return true
}

// case insensitive, tests a string against some criteria and returns a score BYTE_MAX = 255 <=> 100%
func LikelyValidData(s string, nopes []string) (score ValidScore) {
	score = ValidScore(False)
	if LikelyInAList(s, nopes) == ValidScore(False) && LikelyAWord(s) == ValidScore(True) {
		return ValidScore(True)
	}
	return score
}

// case insensitive, tests if a string is a word and returns a score BYTE_MAX = 255 <=> 100%
func LikelyAWord(s string) (score ValidScore) {
	score = ValidScore(False)
	r := regexp.MustCompile(`(?i)\w+`)

	if r.MatchString(s) {
		score = ValidScore(True)
	} else {
		score = ValidScore(False)
	}
	return score
}

// case insensitive, tests if a string is in a list and returns a score BYTE_MAX = 127 <=> 100%
func LikelyInAList(s string, list []string) (score ValidScore) {
	score = ValidScore(False)
	for _, l := range list {
		if strings.Contains(s, l) {
			score = ValidScore(True)
			break
		}
	}
	return score
}

//
// boolean versions

// case insensitive, tests a string against some criteria and returns a score BYTE_MAX = 255 <=> 100%
func IsLikelyValidData(s string, nopes []string) (ok bool) {
	ok = false
	if LikelyValidData(s, nopes) == ValidScore(True) {
		ok = true
	}
	return ok
}

// case insensitive, tests if a string is a word and returns a score BYTE_MAX = 255 <=> 100%
func IsLikelyAWord(s string) (ok bool) {
	ok = false

	if LikelyAWord(s) == ValidScore(True) {
		ok = true
	}
	return ok
}

// case insensitive, tests if a string is in a list and returns a score BYTE_MAX = 127 <=> 100%
func IsInAList(s string, list []string) (ok bool) {
	ok = false
	if LikelyInAList(s, list) == ValidScore(True) {
		ok = true
	}
	return ok
}

func IsAName(s string) (ok bool) {

	return
}

func ImportPhrases(filename string) (phrases []string) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil
	}

	r := bufio.NewScanner(f)

	for r.Scan() {
		phrases = append(phrases, strings.Trim(r.Text(), "\" "))

	}
	for i, phrase := range phrases {
		log.Printf("Phrase #%d: %s\n", i, phrase)
	}
	WaitForKeypress()
	return phrases
}

// Specific data types validations

// IsANamesLine
// L, F  || F M L || F M. L || F L || F "M" L (NOTES)
func IsANamesLine(data string) (fname string, mname string, lname string, notes string, ok bool) {

	if len(data) == 0 {
		//errors.Assert(len(data) != 0, "ExtractNamesNotesFrom data is empty")
		log.Printf("ExtractNamesNotesFrom data is empty. returning defaults and false\n")
		return "NULL_FNAME", "NULL_MNAME", "NULL_LNAME", "NULL_NOTES", false
	}

	fname, mname, lname, notes, ok = "NULL_FNAME", "NULL_MNAME", "NULL_LNAME", "NULL_NOTES", false

	// split names away from notes through `(`, if exists
	names, notes := "", "NULL_NOTES"
	switch s := strings.Split(strings.TrimSpace(data), xNOTES_SEP_OPEN); len(s) {
	case 0:
		errors.Assert(false, "ExtractFrom switch Split error data likely nil/empty")
	case 1:
		if strings.Contains(s[0], xNOTES_SEP_OPEN+xNOTES_SEP_CLOSE) {
			errors.Assert(false, "ExtractFrom Contains error data likely conmtains only notes but no names")
		} else {
			names = s[0]
		}
	case 2:
		names = strings.TrimSpace(s[0])
		notes = strings.TrimSpace(strings.Trim(s[1], xNOTES_SEP_OPEN+xNOTES_SEP_CLOSE))
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
		mname = "NULL_MNAME"
		fname = "NULL_FNAME"
		ok = true
	case len(s1) > 0 && len(s1[0]) == 3:
		lname = s1[0][1]
		mname = "NULL_MNAME"
		fname = s1[0][2]
		ok = true
	case len(s2) > 0 && len(s2[0]) == 3:
		lname = s2[0][2]
		mname = "NULL_MNAME"
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
		mname = "NULL_MNAME"
		fname = "NULL_FNAME"
		ok = true
	default:
		// Errors
		log.Printf("####### IsANamesLine WARNING UNDEFINED REGREX FOR names: %#v \n", names)
		log.Printf("REGEX s 0 %#v\n", s0)
		log.Printf("REGEX s 1 %#v\n", s1)
		log.Printf("REGEX s 2 %#v\n", s2)
		log.Printf("REGEX s 3 %#v\n", s3)
		log.Printf("REGEX s 4 %#v\n", s4)
		log.Printf("REGEX s 5 %#v\n", s5)
		log.Printf("REGEX s 6 %#v\n", s6)
		errors.Assert(false, " sANamesLine  UNDEFINED REGREX FOR names")
		lname = "NULL_LNAME"
		mname = "NULL_MNAME"
		fname = "NULL_FNAME"
		ok = false
	}

	//if len(lname) > 1 {
	//	ok = true
	//} else {
	//	ok = false
	//}
	return fname, mname, lname, notes, ok
}

// Utilities to make data valid

// NormalizeStr converts a string to Uppercase and remove spaces around
// returns the changed string and true if success, otherwise false if in our out string is invalid
func NormalizeStr(in string) (out string) { //, err error) {
	//if IsUnwantedInput(in, goodset, badset) {
	//	return "", false
	//}
	out = strings.ToUpper(strings.TrimSpace(in))
	return out

}

func NormalizeKey(in string) (out string) { //, err error) {
	//if IsUnwantedInput(in, goodset, badset) {
	//	return "", false
	//}
	out = NormalizeStr(in)
	return out

}

func NormalizeValue(in string) (out string) { //, err error) {
	//if IsUnwantedInput(in, goodset, badset) {
	//	return "", false
	//}
	out = strings.ToLower(strings.TrimSpace(in))
	return out
}

func NormalizeField(in string) (out string) { //, err error) {
	//if IsUnwantedInput(in, goodset, badset) {
	//	return "", false
	//}
	out = strings.ToTitle(strings.TrimSpace(in))
	return out
}

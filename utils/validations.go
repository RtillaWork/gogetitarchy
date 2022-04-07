package utils

import (
	"bufio"
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

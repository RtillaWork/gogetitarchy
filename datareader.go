package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type MusiciansMap map[HashSum]*Musician

func ReadMusicianData(inFileName string) MusiciansMap {

	inFile, err := os.Open(inFileName)
	FailOn(err, "opening inFile for reading...")
	defer inFile.Close()

	musicians := make(map[HashSum]*Musician)

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

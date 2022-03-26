package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func FailOn(err error, desc string) {
	const SEP = "\n\n"
	if err != nil {
		log.Printf(SEP+"ERR: %s"+SEP, err)
		panic(err)
	} else {
		log.Printf(SEP+"INFO: %s"+SEP, desc)
	}
}

func FailNotOK(ok bool, desc string) {
	const SEP = "\n\n"
	if !ok {
		log.Printf(SEP+"ERR: FAILED ON NOT OK: %s"+SEP, desc)
		panic(errors.New(desc))
	} else {
		log.Printf(SEP+"INFO: OK'ed %s"+SEP, desc)
	}
}

const inFileName = "../inFile.txt"

func main() {
	musicians := ReadMusicianData(inFileName)
	//exportAllMusicians(musicians, "")
	musiciansQueries := BuildQueries(musicians)
	//exportAllqueries(musicians, musiciansQueries, "")

	musiciansResponseData := CrawlArchiveGrid(musicians, musiciansQueries)
	exportAllResponseData(musicians, musiciansResponseData, "")

}

func exportAllMusicians(musicians MusiciansMap, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.Open(filename); err != nil {
		log.Printf("Eeeor opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for _, m := range musicians {
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  LAST: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FirstName, m.LastName, m.MiddleName, m.Notes)
		//log.Println(m.ToCsv())
		fmt.Fprintf(outfile, "%d; %s", counter, m.ToCsv())
		counter++
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func exportAllqueries(ms MusiciansMap, mqs MusiciansQueries, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.Open(filename); err != nil {
		log.Printf("Eeeor opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for m, mq := range mqs {
		//log.Printf("\n COUNTER: %d Musician{%s}\nQuery{%s}\n\n", counter, ms[m], mq)
		fmt.Fprintf(outfile, "%d; %q; %q", counter, ms[m], mq)
		counter++
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func exportAllResponseData(ms MusiciansMap, mrd MusiciansData, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.Open(filename); err != nil {
		log.Printf("Eeeor opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for k, _ := range mrd {

		//log.Printf("%s\n", ms[k].ToCsv())
		fmt.Fprintf(outfile, "%d; %q", counter, ms[k].ToCsv())
		counter++
	}

	log.Printf("TOTAL DATA FOUND ABOUT ALL MUSICANS: %d\n", counter)
}

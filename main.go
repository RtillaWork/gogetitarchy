package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func FailOn(err error, desc string) {
	const SEP = "\n\n"
	if err != nil {
		//log.Printf(SEP+"ERR: %s"+SEP, err)
		panic(err)
	} else {
		//log.Printf(SEP+"INFO: %s"+SEP, desc)
	}
}

func FailNotOK(ok bool, desc string) {
	const SEP = "\n\n"
	if !ok {
		//log.Printf(SEP+"ERR: FAILED ON NOT OK: %s"+SEP, desc)
		panic(errors.New(desc))
	} else {
		//log.Printf(SEP+"INFO: OK'ed %s"+SEP, desc)
	}
}

const inFileName = "../inFile.txt"

func main() {

	// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
	musicians := ReadMusicianData(inFileName)
	//if len(os.Args) == 2 {
	//	exportAllMusicians(musicians, os.Args[1])
	//} else {
	//	exportAllMusicians(musicians, "")
	//}
	musiciansQueries := BuildQueries(musicians)
	//exportAllqueries(musicians, musiciansQueries, "")

	var phrases []string = nil
	if len(os.Args) == 2 {
		phrases = importPhrases(os.Args[1])
	} else { // DEBUG TEMPORARY
		phrases = importPhrases("./phrases.csv")
	}
	musiciansResponseData, ok := CrawlArchiveGrid(musicians, musiciansQueries, 10, phrases)
	if ok {
		exportAllResponseData(musicians, musiciansResponseData, "")
	} else {
		log.Println("CrawlArchiveGrid returned not ok")
	}

}

func importPhrases(filename string) (phrases []string) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil
	}

	r := bufio.NewScanner(f)

	for r.Scan() {
		phrases = append(phrases, r.Text())

	}
	for i, phrase := range phrases {
		log.Printf("Phrase #%d: %s\n", i, phrase)
	}
	WaitForKeypress()
	return phrases
}

func exportAllMusicians(musicians MusiciansMap, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.Open("OUT_MUSICIANS_" + filename); err != nil {
		log.Printf("Error opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for _, m := range musicians {
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  LAST: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FirstName, m.LastName, m.MiddleName, m.Notes)
		//log.Println(m.ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n===================")
		}
		fmt.Fprintf(outfile, "\n%d; %s\n", counter, m.ToCsv())
		counter++
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func exportAllqueries(ms MusiciansMap, mqs MusiciansQueries, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.Open("OUT_QUERIES_" + filename); err != nil {
		log.Printf("Error opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for m, mq := range mqs {
		//log.Printf("\n COUNTER: %d Musician{%s}\nQuery{%s}\n\n", counter, ms[m], mq)
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n===================")
		}
		fmt.Fprintf(outfile, "\n%d; %q; %q\n", counter, ms[m], mq)
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
	for h, records := range mrd {

		//log.Printf("%s\n", ms[k].ToCsv())
		fmt.Fprintf(outfile, "\n\n======================================\n%d; %q\n", counter, ms[h].ToCsv())
		counter++
		for hh, record := range records {
			//fmt.Fprintf(outfile, "%s >> %q\n", hh, record.ToCsv())

			fmt.Fprintf(outfile, "%s >> %q\n", hh, record.ToCsv())
		}
	}

	log.Printf("TOTAL DATA FOUND ABOUT ALL MUSICANS: %d\n", counter)
}

func WaitForKeypress() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press key... ")
	_, err := reader.ReadString('\n')
	FailOn(err, "WaitForKeypress Failed")
	fmt.Println("RESUMING...")
}

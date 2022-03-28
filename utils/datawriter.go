package utils

import (
	"fmt"
	"github.com/RtillaWork/gogetitarchy/archivegrid"
	"github.com/RtillaWork/gogetitarchy/musician"
	"log"
	"os"
	"strings"
)

func BuildOutput(musiciansData archivegrid.MusiciansData) {

	//outFileName := "outFile.csv"
	//outFile, err := os.Create(outFileName)
	//if err != nil {
	//	log.Fatalf("LOG: ERROR CREATING %s, err %q", outFileName, err)
	//	return
	//}
	//defer outFile.Close()
	//
	//writer := csv.NewWriter(outFile)
	//defer writer.Flush()
}

func ExportAllMusicians(musicians musician.MusiciansMap, filename string) {
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

func ExportAllqueries(ms musician.MusiciansMap, mqs MusiciansQueries, filename string) {
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

func ExportAllResponseData(ms musician.MusiciansMap, mrd archivegrid.MusiciansData, filename string) {
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

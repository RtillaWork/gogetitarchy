package archivegrid

import (
	"fmt"
	"github.com/RtillaWork/gogetitarchy/musician"
	"log"
	"os"
	"strings"
)

func ExportAllqueries(ms musician.MusiciansMap, mqs MusiciansQueries, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.OpenFile("OUT_QUERIES_"+filename, os.O_WRONLY, 0777); err != nil {
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

func ExportAllResponseData(ms musician.MusiciansMap, mrd MusiciansData, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.OpenFile(filename, os.O_WRONLY, 0777); err != nil {
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

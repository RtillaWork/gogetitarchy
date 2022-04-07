package archivegrid

import (
	"fmt"
	"github.com/RtillaWork/gogetitarchy/musician"
	"log"
	"os"
)

func ExportAllqueries(ms musician.MusiciansMap, mqs MusiciansQueries, filename string) {
	var outfile *os.File

	if filename == "" {
		outfile = os.Stdout
	} else if h, err := os.OpenFile("OUT_QUERIES_"+filename, os.O_WRONLY, 0777); err != nil {
		log.Printf("Error opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}

	if outfile == os.Stdout {
		fmt.Fprintf(outfile, "\n=== BEGIN QUERIES ========")
	} else {
		fmt.Fprintf(outfile, "[\n")
	}
	counter := 1
	for _, mq := range mqs {
		//log.Printf("\n COUNTER: %d Musician{%s}\nQuery{%s}\n\n", counter, ms[m], mq)
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n==== QUERY %d==========", counter)
		}
		//fmt.Fprintf(outfile, "\n%d; %q; %q\n", counter, ms[m], mq)
		fmt.Fprintf(outfile, "%s,\n", mq.ToJson())
		counter++
	}
	if outfile == os.Stdout {
		fmt.Fprintf(outfile, "\n=== END QUERIES ========")
	} else {
		fmt.Fprintf(outfile, "]\n")
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func ExportAllResponseData(ms musician.MusiciansMap, mrd MusiciansData, filename string) {
	var outfile *os.File
	if filename == "" {
		outfile = os.Stdout
	} else if h, err := os.OpenFile(filename, os.O_WRONLY, 0777); err != nil {
		log.Printf("Eeeor opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}

	//
	if outfile == os.Stdout {
		fmt.Fprintf(outfile, "\n=== BEGIN RESPONSE DATA ========")
	} else {
		fmt.Fprintf(outfile, "[\n")
	}
	counter := 1
	for h, records := range mrd {

		//log.Printf("%s\n", ms[k].ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n\n======== Musician=================\n%d; %q\n", counter, ms[h].ToCsv())
		}
		//} else {
		//	fmt.Fprintf(outfile, "[\n")
		//}
		counter++
		for _, record := range records {
			//fmt.Fprintf(outfile, "%s >> %q\n", hh, record.ToCsv())

			//fmt.Fprintf(outfile, "%s >> %q\n", hh, record.ToCsv())
			fmt.Fprintf(outfile, "%s,\n", record.ToJson())
		}
	}
	if outfile == os.Stdout {
		fmt.Fprintf(outfile, "\n=== END RESPONSE DATA ========")
	} else {
		fmt.Fprintf(outfile, "]\n")
	}

	log.Printf("TOTAL DATA FOUND ABOUT ALL MUSICANS: %d\n", counter)
}

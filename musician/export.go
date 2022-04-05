package musician

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func ExportAll(musicians MusiciansMap, filename string) {
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
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  L: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FName, m.LName, m.MName, m.Notes)
		//log.Println(m.ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n===================")
		}
		fmt.Fprintf(outfile, "\n%d; %s\n", counter, m.ToCsv())
		counter++
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func ExportJson(musicians MusiciansMap, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".json") {
		outfile = os.Stdout
	} else if h, err := os.Open("OUT_MUSICIANS_" + filename); err != nil {
		log.Printf("Error opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for _, m := range musicians {
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  L: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FName, m.LName, m.MName, m.Notes)
		//log.Println(m.ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n===================")
		}
		fmt.Fprintf(outfile, "\n%d; %s\n", counter, m.ToJson())
		counter++
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func ExportCsv(musicians MusiciansMap, filename string) {
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
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  L: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FName, m.LName, m.MName, m.Notes)
		//log.Println(m.ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n===================")
		}
		fmt.Fprintf(outfile, "\n%d; %s\n", counter, m.ToCsv())
		counter++
	}
	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
}

func ExportDataDict(dict DataDict, filename string) {
	var outfile *os.File
	if filename == "" || !strings.HasSuffix(filename, ".csv") {
		outfile = os.Stdout
	} else if h, err := os.Open("OUT_MUSICIANS_DATADICT" + filename); err != nil {
		log.Printf("Error opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		outfile = h
	}
	counter := 1
	for k, vs := range dict.Fields {
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  L: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FName, m.LName, m.MName, m.Notes)
		//log.Println(m.ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n== KEY ==============")
		}
		fmt.Fprintf(outfile, "\n\n\n### KEY: %s  ### [%5d]\n", k, counter)
		fmt.Fprintf(outfile, "KEY: %s STATS ### [%5d]\n", k, dict.KeyStats[k])
		for _, v := range vs {
			fmt.Fprintf(outfile, "VALUE: %s  ( VALUE STATS: [%5d] \n", v, dict.ValuesStats[v])
		}

		counter++
	}
	log.Printf("\n\n\n Counted Number of Keys: %d\n\n", counter)
}

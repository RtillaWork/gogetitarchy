package musician

import (
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils"
	"log"
	"os"
	"strings"
)

func ExportJson(musicians MusiciansMap, filename string) {
	var outfile *os.File
	if filename == "" {
		outfile = os.Stdout
		log.Printf("Opening file outfile=OS.STDOUT")
	} else if h, err := os.OpenFile(filename, os.O_WRONLY, 0666); err != nil {
		log.Printf("Error opening file: %s \n%v\n", outfile, err)
		outfile = os.Stdout
	} else {
		log.Printf("Opening file OUT_time.now_+%s", filename)
		outfile = h
		utils.WaitForKeypress()
	}
	counter := 1
	for _, m := range musicians {
		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  L: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FName, m.LName, m.MName, m.Notes)
		//log.Println(m.ToCsv())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n=== BEGIN RECORD %d ==========", counter)
		} else {
			fmt.Fprintf(outfile, "[\n")
		}
		fmt.Fprintf(outfile, "%s,\n", m.ToJson())
		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n===END RECORD ==========")
		} else {
			fmt.Fprintf(outfile, "]")
		}
		counter++
	}
	log.Printf("\n\n\n Json records exported for musicians: %d\n\n", counter)
	//utils.WaitForKeypress()
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
	if filename == "" {
		outfile = os.Stdout
	} else if h, err := os.OpenFile(filename, os.O_WRONLY, 0666); err != nil {
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
			fmt.Fprintf(outfile, "\n== KEY %s [counter %d] ==============", k, counter)
		} else {
			fmt.Fprintf(outfile, "{\n")
		}
		fmt.Fprintf(outfile, "\n\n\n### KEY: %s  ### [%d]\n", k, counter)
		fmt.Fprintf(outfile, "KEY: %s STATS ### [%d]\n", k, dict.KeyStats[k])
		for _, v := range vs {
			fmt.Fprintf(outfile, "VALUE: %s  ( VALUE STATS: [%d] \n", v, dict.ValuesStats[v])
		}

		if outfile == os.Stdout {
			fmt.Fprintf(outfile, "\n== END KEY %s [counter %d] ==============", k, counter)
		} else {
			fmt.Fprintf(outfile, "}\n")
		}

		counter++
	}
	log.Printf("\n\n\n Counted Number of Keys: %d\n\n", counter)
	//utils.WaitForKeypress()
}

// OLD

//func ExportAll(musicians MusiciansMap, filename string) {
//	var outfile *os.File
//	if filename == "" || !strings.HasSuffix(filename, ".csv") {
//		outfile = os.Stdout
//	} else if h, err := os.Open("OUT_MUSICIANS_" + filename); err != nil {
//		log.Printf("Error opening file: %s \n%v\n", outfile, err)
//		outfile = os.Stdout
//	} else {
//		outfile = h
//	}
//	counter := 1
//	for _, m := range musicians {
//		//log.Printf("{KEY: %s ,,,, VALUE: {FIRST: %s  L: %s   MIDDLE:  %s   NOTES: %s  }", k, m.FName, m.LName, m.MName, m.Notes)
//		//log.Println(m.ToCsv())
//		if outfile == os.Stdout {
//			fmt.Fprintf(outfile, "\n===================")
//		}
//		fmt.Fprintf(outfile, "\n%d; %s\n", counter, m.ToCsv())
//		counter++
//	}
//	log.Printf("\n\n\n SIZE of musicians: %d\n\n", counter)
//}

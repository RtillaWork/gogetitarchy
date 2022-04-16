package main

import (
	"flag"
	"github.com/RtillaWork/gogetitarchy/archivegrid"
	"github.com/RtillaWork/gogetitarchy/archivegrid/query"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/testing"
	"github.com/RtillaWork/gogetitarchy/utils"
	"log"
	"os"
	"strconv"
	"time"
)

// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
//const inRawFileNameDefault = "../inFile.txt"
var InRawFileNameDefault = "../infantry_RAW3_in_nospaces.txt" // "../infantry_raw_in.txt"
var OutRawRebuiltfilenameDefault = OutMusiciansFilenameDefault + "_OUT_RAW_REBUILT.txt"
var FilterPhrasesFilenameDefault = "../phrases.csv"
var OutMusiciansFilenameDefault = "/home/webdev/_ARCHIVEGRID/musiciansdefault" + strconv.FormatInt(time.Now().Unix(), 10)
var InMusiciansFilenameDefault = "/home/webdev/_ARCHIVEGRID/musiciansdefault"
var OutMusiciansDbFilenameDefault = OutMusiciansFilenameDefault + "_DB_"
var OutTheDataDictFilenameDefault = OutMusiciansDbFilenameDefault + "_DATADICT_"
var OutMusiciansQueryFilenameDefault = OutMusiciansFilenameDefault + "_QUERIES_"
var OutResponseDataFilenameDefault = OutMusiciansFilenameDefault + "_RESPONSERECORDS_"
var OutExtensionDefault = ".json" // or ".csv"
var testModeDefault = false

func main() {
	InRawFilename := flag.String("inRaw", InRawFileNameDefault, "Input Raw Musicians filename")
	OutMusiciansFilename := flag.String("outMusicians", OutMusiciansFilenameDefault, "Output Musicians filename")
	InMusiciansFilename := flag.String("InMusicians", InMusiciansFilenameDefault, "In Musicians filename")
	//OutMusiciansDbFilename := flag.String("outMusiciansDbFilename", OutMusiciansDbFilenameDefault, "Output MusiciansDb filename")
	OutTheDataDictFilename := flag.String("outTheDatadict", OutTheDataDictFilenameDefault, "Output Data dictionary filename in json")
	OutMusiciansQueryFilename := flag.String("outQueries", OutMusiciansQueryFilenameDefault, "Output queries json")
	OutResponseDataFilename := flag.String("outResponse", OutResponseDataFilenameDefault, "Output response data in json")
	OutExtension := flag.String("outformat", OutExtensionDefault, "Output format json or csv(;). Default json")
	testMode := flag.Bool("testMode", testModeDefault, "compare computed with saved (default true)")
	FilterKeywordsFilename := flag.String("filterPhrases", FilterPhrasesFilenameDefault, "Input filter-in GoodSetKeywords in csv format")

	flag.Parse()

	GoodSetKeywords := utils.ImportPhrases(*FilterKeywordsFilename)

	//
	var musicians musician.MusiciansMap
	if d, err := os.ReadFile(*InMusiciansFilename + *OutExtension); err == nil {
		log.Printf("Musicians file %s found, reading...\n", *InMusiciansFilename)
		musicians = musician.ReadData(d)
		log.Printf("Musicians file %s found, read %d musicians\n", *InMusiciansFilename, len(musicians))
		utils.WaitForKeypress()
	}
	if musicians == nil {
		if d, err := os.ReadFile(*OutMusiciansFilename + *OutExtension); err != nil {
			log.Printf("Musicians file %s not found, importing...\n", *OutMusiciansFilename)
			musicians = musician.Import(*InRawFilename, musician.BlockDelimDef1, musician.BlockDelimDef2)
			musician.ExportJson(musicians, *OutMusiciansFilename+*OutExtension)
		} else {
			log.Printf("Musicians file %s found, reading...\n", *OutMusiciansFilename)
			musicians = musician.ReadData(d)
			log.Printf("Musicians file %s found, read %d musicians\n", *OutMusiciansFilename, len(musicians))
			utils.WaitForKeypress()
		}

		if *testMode {
			var testmusiciansA, testmusiciansB musician.MusiciansMap
			testmusiciansA = musician.ImportData(*InRawFilename, musician.BlockDelimDef1, musician.BlockDelimDef2)
			if d, err := os.ReadFile(*OutMusiciansFilename); err != nil {
				log.Printf("NO file %s to test against", *OutMusiciansFilename)
			} else {
				testmusiciansB = musician.ReadData(d)
			}
			testing.CompareMusicians(&testmusiciansA, &testmusiciansB)
		}
	}

	//musiciansdb := musician.NewMusiciansDb(musicians)
	musician.ExportDataDict(musician.TheDataDict, *OutTheDataDictFilename+*OutExtension)

	//
	musiciansQueries := query.BuildQueries(musicians)
	musiciansResponseData, ok := archivegrid.CrawlArchiveGrid(musicians, musiciansQueries, 3900, GoodSetKeywords)
	archivegrid.ExportAllqueries(musicians, musiciansQueries, *OutMusiciansQueryFilename+*OutExtension)

	if ok {
		archivegrid.ExportAllResponseData(musicians, musiciansResponseData, *OutResponseDataFilename+*OutExtension)
		markedResponseData := archivegrid.MusiciansData{}
		for mh, rs := range musiciansResponseData {
			for _, r := range rs {
				if r.IsMatch {
					markedResponseData[mh] = append(markedResponseData[mh], r)
				}
			}
		}
		archivegrid.ExportAllResponseData(musicians, markedResponseData, "MARKEDDATA"+*OutResponseDataFilename+*OutExtension)
	} else {
		log.Println("CrawlArchiveGrid returned not ok")
	}

}

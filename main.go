package main

import (
	"flag"
	"github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils"
	"log"
	"os"
	"strconv"
	"time"
)

// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
//const inRawFileNameDefault = "../inFile.txt"
var InRawFileNameDefault = "../infantry_raw_in.txt"
var OutRawRebuiltfilenameDefault = OutMusiciansFilenameDefault + "_OUT_RAW_REBUILT.txt"
var FilterPhrasesFilenameDefault = "../phrases.csv"
var OutMusiciansFilenameDefault = "/home/webdev/_ARCHIVEGRID/musiciansdefault" + strconv.FormatInt(time.Now().Unix(), 10)
var OutMusiciansDbFilenameDefault = OutMusiciansFilenameDefault + "_DB_"
var OutTheDataDictFilenameDefault = OutMusiciansDbFilenameDefault + "_DATADICT_"
var OutMusiciansQueryFilenameDefault = OutMusiciansFilenameDefault + "_QUERIES_"
var OutResponseDataFilenameDefault = OutMusiciansFilenameDefault + "_RESPONSERECORDS_"
var OutExtensionDefault = ".json" // or ".csv"
var testModeDefault = false

func main() {
	InRawFilename := flag.String("inRaw", InRawFileNameDefault, "Input Raw Musicians filename")
	//FilterPhrasesFilename := flag.String("filterPhrases", FilterPhrasesFilenameDefault, "Input filter-in GoodSetPhrases in csv format")
	OutMusiciansFilename := flag.String("outMusicians", OutMusiciansFilenameDefault, "Output Musicians filename")
	//OutMusiciansDbFilename := flag.String("outMusiciansDbFilename", OutMusiciansDbFilenameDefault, "Output MusiciansDb filename")
	//OutTheDataDictFilename := flag.String("outTheDatadict", OutTheDataDictFilenameDefault, "Output Data dictionary filename in json")
	//OutMusiciansQueryFilename := flag.String("outQueries", OutMusiciansQueryFilenameDefault, "Output queries json")
	//OutResponseDataFilename := flag.String("outResponse", OutResponseDataFilenameDefault, "Output response data in json")
	OutExtension := flag.String("outformat", OutExtensionDefault, "Output format json or csv(;). Default json")
	testMode := flag.Bool("testMode", testModeDefault, "compare computed with saved (default true)")
	flag.Parse()
	//GoodSetPhrases := utils.ImportPhrases(*FilterPhrasesFilename)

	//
	var musicians musician.MusiciansMap
	if d, err := os.ReadFile(*OutMusiciansFilename); err != nil {
		musicians = musician.ImportData(*InRawFilename, musician.BlockDelimDef1)
		musician.ExportJson(musicians, *OutMusiciansFilename+*OutExtension)
	} else {
		musicians = musician.ReadData(d)
		log.Printf("Musicians file %s found, imported %d musicians\n", *OutMusiciansFilename, len(musicians))
		utils.WaitForKeypress()
	}

	if *testMode {
		var testmusiciansA, testmusiciansB musician.MusiciansMap
		testmusiciansA = musician.ImportData(*InRawFilename, musician.BlockDelimDef1)
		if d, err := os.ReadFile(*OutMusiciansFilename); err != nil {
			log.Printf("NO file %s to test against", *OutMusiciansFilename)
		} else {
			testmusiciansB = musician.ReadData(d)
		}
		utils.CompareMusicians(&testmusiciansA, &testmusiciansB)
	}

	//musiciansdb := musician.NewMusiciansDb(musicians)
	//musician.ExportDataDict(*musiciansdb.Dict, *OutTheDataDictFilename+*OutExtension)
	//
	////
	//musiciansQueries := archivegrid.BuildQueries(musicians)
	//archivegrid.ExportAllqueries(musicians, musiciansQueries, *OutMusiciansQueryFilename)
	//
	////
	//musiciansResponseData, ok := archivegrid.CrawlArchiveGrid(musicians, musiciansQueries, 3, GoodSetPhrases)
	//if ok {
	//	archivegrid.ExportAllResponseData(musicians, musiciansResponseData, *OutResponseDataFilename)
	//} else {
	//	log.Println("CrawlArchiveGrid returned not ok")
	//}

}

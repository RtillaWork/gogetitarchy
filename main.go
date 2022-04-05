package main

import (
	"flag"
	"github.com/RtillaWork/gogetitarchy/musician"
)

// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
//const inRawFileNameDefault = "../inFile.txt"
const InRawFileNameDefault = "../infantry_raw_in.txt"
const FilterPhrasesDefault = "../phrases.csv"
const OutMusiciansFilenameDefault = "../out_musicians_default"
const OutMusiciansDbFilenameDefault = "../out_musiciansdb_default"
const OutTheDataDictFilenameDefault = OutMusiciansDbFilenameDefault + "_DataDict"
const OutExtensionDefault = ".json" // or ".csv"

func main() {
	InRawFilename := flag.String("inRaw", InRawFileNameDefault, "Input Raw Musicians filename")
	FilterPhrases := flag.String("filterPhrases", FilterPhrasesDefault, "Input filter-in phrases in csv format")
	OutMusiciansFilename := flag.String("outMusicians", OutMusiciansFilenameDefault, "Output Musicians filename")
	//OutMusiciansDbFilename := flag.String("outMusiciansDbFilename", OutMusiciansDbFilenameDefault, "Output MusiciansDb filename")
	OutTheDataDictFilename := flag.String("outTheDatadict", OutTheDataDictFilenameDefault, "Output Data dictionary filename in json")
	OutExtension := flag.String("outformat", OutExtensionDefault, "Output format json or csv(;). Default json")
	flag.Parse()

	//musicians := musician.ReadMusiciansNames(inRawFileNameDefault)
	musicians := musician.ImportData(*InRawFilename, musician.BlockDelimDef)
	musiciansdb := musician.NewMusiciansDb(musicians)
	//if len(os.Args) == 2 {
	musician.ExportJson(musiciansdb.Musicians, *OutMusiciansFilename+*OutExtension)
	musician.ExportDataDict(musiciansdb.Dict, *OutTheDataDictFilename+*OutExtension)
	//} else {
	//
	//	musician.ExportJson(musiciansdb.Musicians, "")
	//	musician.ExportDataDict(musiciansdb.Dict, "")
	//}

	//musiciansQueries := archivegrid.BuildQueries(musicians)
	//exportAllqueries(musicians, musiciansQueries, "")
	//
	//var phrases []string = nil
	//if len(os.Args) == 2 {
	//	phrases = archivegrid.ImportPhrases(os.Args[1])
	//} else { // DEBUG TEMPORARY
	//	phrases = archivegrid.ImportPhrases("./phrases.csv")
	//}
	//musiciansResponseData, ok := archivegrid.CrawlArchiveGrid(musicians, musiciansQueries, 1, phrases)
	//if ok {
	//	archivegrid.ExportAllResponseData(musicians, musiciansResponseData, "")
	//} else {
	//	log.Println("CrawlArchiveGrid returned not ok")
	//}

}

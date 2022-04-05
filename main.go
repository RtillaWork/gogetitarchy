package main

import (
	"github.com/RtillaWork/gogetitarchy/musician"
	"os"
)

// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
//const inFileName = "../inFile.txt"
const inFileName = "../infantry_raw_in.txt"

func main() {

	//musicians := musician.ReadMusiciansNames(inFileName)
	musicians := musician.ImportData(inFileName, musician.BlockDelimDef)
	musiciansdb := musician.NewMusiciansDb(musicians)
	if len(os.Args) == 2 {
		musician.ExportJson(musiciansdb.Musicians, os.Args[1])
		musician.ExportDataDict(musiciansdb.Dict, os.Args[1])
	} else {

		musician.ExportJson(musiciansdb.Musicians, "")
		musician.ExportDataDict(musiciansdb.Dict, "")
	}

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

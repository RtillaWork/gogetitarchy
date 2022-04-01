package main

import (
	"github.com/RtillaWork/gogetitarchy/musician"
	"os"
)

const inFileName = "../inFile.txt"

func main() {

	//flag := flag.NewFlagSet("inrawfile")
	// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
	musicians := musician.ReadMusiciansNames(inFileName)
	if len(os.Args) == 2 {
		musician.ExportAllMusicians(musicians, os.Args[1])
	} else {
		musician.ExportAllMusicians(musicians, "")
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

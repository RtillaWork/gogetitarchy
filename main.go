package main

import (
	"github.com/RtillaWork/gogetitarchy/archivegrid"
	"github.com/RtillaWork/gogetitarchy/musician"
	"log"
	"os"
)

const inFileName = "../inFile.txt"

func main() {

	// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
	musicians := musician.ReadMusiciansNames(inFileName)
	//if len(os.Args) == 2 {
	//	exportAllMusicians(musicians, os.Args[1])
	//} else {
	//	exportAllMusicians(musicians, "")
	//}
	musiciansQueries := archivegrid.BuildQueries(musicians)
	//exportAllqueries(musicians, musiciansQueries, "")

	var phrases []string = nil
	if len(os.Args) == 2 {
		phrases = archivegrid.ImportPhrases(os.Args[1])
	} else { // DEBUG TEMPORARY
		phrases = archivegrid.ImportPhrases("./phrases.csv")
	}
	musiciansResponseData, ok := archivegrid.CrawlArchiveGrid(musicians, musiciansQueries, 1, phrases)
	if ok {
		archivegrid.ExportAllResponseData(musicians, musiciansResponseData, "")
	} else {
		log.Println("CrawlArchiveGrid returned not ok")
	}

}

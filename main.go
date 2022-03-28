package main

import (
	archivegrid "github.com/RtillaWork/gogetitarchy/archivegrid"
	musician "github.com/RtillaWork/gogetitarchy/musician"
	"github.com/RtillaWork/gogetitarchy/utils"
	"log"
	"os"
)

const inFileName = "../inFile.txt"

func main() {

	// archy INPHRASES IMPORTRAWMUSICIANS EXPORTJSONORCSVMUSICIANS
	musicians := musician.ReadMusicianData(inFileName)
	//if len(os.Args) == 2 {
	//	exportAllMusicians(musicians, os.Args[1])
	//} else {
	//	exportAllMusicians(musicians, "")
	//}
	musiciansQueries := archivegrid.BuildQueries(musicians)
	//exportAllqueries(musicians, musiciansQueries, "")

	var phrases []string = nil
	if len(os.Args) == 2 {
		phrases = utils.ImportPhrases(os.Args[1])
	} else { // DEBUG TEMPORARY
		phrases = utils.ImportPhrases("./phrases.csv")
	}
	musiciansResponseData, ok := archivegrid.CrawlArchiveGrid(musicians, musiciansQueries, 10, phrases)
	if ok {
		utils.ExportAllResponseData(musicians, musiciansResponseData, "")
	} else {
		log.Println("CrawlArchiveGrid returned not ok")
	}

}

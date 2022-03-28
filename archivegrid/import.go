package archivegrid

import (
	"bufio"
	"github.com/RtillaWork/gogetitarchy/utils"
	"log"
	"os"
	"strings"
)

func ImportPhrases(filename string) (phrases []string) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil
	}

	r := bufio.NewScanner(f)

	for r.Scan() {
		phrases = append(phrases, strings.Trim(r.Text(), "\" "))

	}
	for i, phrase := range phrases {
		log.Printf("Phrase #%d: %s\n", i, phrase)
	}
	utils.WaitForKeypress()
	return phrases
}

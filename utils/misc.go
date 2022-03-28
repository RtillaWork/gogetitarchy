package utils

import (
	"bufio"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
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
	WaitForKeypress()
	return phrases
}

func WaitForKeypress() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press key... ")
	_, err := reader.ReadString('\n')
	errors.FailOn(err, "WaitForKeypress Failed")
	fmt.Println("RESUMING...")
}

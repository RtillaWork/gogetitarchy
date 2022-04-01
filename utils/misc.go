package utils

import (
	"bufio"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"os"
)

var lastkey string = ""

func WaitForKeypress() {
	if lastkey == "a" {
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press a to resume all, any other key for next... ")
	key, err := reader.ReadString('\n')
	lastkey = key
	errors.FailOn(err, "WaitForKeypress Failed")
	fmt.Println("RESUMING...")
}

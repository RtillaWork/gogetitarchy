package utils

import (
	"bufio"
	"fmt"
	"github.com/RtillaWork/gogetitarchy/utils/errors"
	"os"
)

func WaitForKeypress() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press key... ")
	_, err := reader.ReadString('\n')
	errors.FailOn(err, "WaitForKeypress Failed")
	fmt.Println("RESUMING...")
}

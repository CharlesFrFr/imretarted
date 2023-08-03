package all

import (
	"fmt"
	"os"
)

func PrintGreen(strings []string) {
	if (os.Getenv("PRODUCTION") == "true") {
		return
	}

	for _, str := range strings {
		fmt.Print("\033[32m" + str + "\033[0m ")
	}
	fmt.Println()
}

func PrintRed(strings []string) {
	if (os.Getenv("PRODUCTION") == "true") {
		return
	}

	for _, str := range strings {
		fmt.Print("\033[31m" + str + "\033[0m ")
	}
	fmt.Println()
}
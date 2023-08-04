package all

import (
	"fmt"
	"os"
)

func PrintGreen(strings []any) {
	if (os.Getenv("PRODUCTION") == "true") {
		return
	}

	for _, str := range strings {
		fmt.Print("\033[32m" + fmt.Sprint(str) + "\033[0m ")
	}
	fmt.Println()
}

func PrintRed(strings []any) {
	if (os.Getenv("PRODUCTION") == "true") {
		return
	}

	for _, str := range strings {
		fmt.Print("\033[31m" + fmt.Sprint(str) + "\033[0m ")
	}
	fmt.Println()
}

func PrintYellow(strings []any) {
	if (os.Getenv("PRODUCTION") == "true") {
		return
	}

	for _, str := range strings {
		fmt.Print("\033[33m" + fmt.Sprint(str) + "\033[0m ")
	}
	fmt.Println()
}
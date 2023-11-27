package utils

import (
	"fmt"
	"log"
	"os"
)

func LogError(m string) {
	err := log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
	err.Println(m)
	fmt.Println()
}

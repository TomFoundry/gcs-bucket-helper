package gcp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	in *bufio.Reader
)

func init() {
	in = bufio.NewReader(os.Stdin)
}

func userInput(promptMessage string, inputValidators ...inputValidator) string {

	fmt.Println("- " + promptMessage)

OUTER:
	for {
		s, _ := in.ReadString('\n')

		s = strings.TrimSuffix(s, "\n")

		for _, validator := range inputValidators {
			if err := validator(s); err != nil {
				fmt.Println(err)
				fmt.Println("Please try again")
				continue OUTER
			}
		}

		return s
	}
}

type inputValidator func(string) error

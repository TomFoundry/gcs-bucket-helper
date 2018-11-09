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
			if err := validator.validatorFunc(s); err != nil {
				fmt.Println(validator.failMessage)
				continue OUTER
			}
		}

		return s
	}
}

type inputValidator struct {
	validatorFunc inputValidatorFunc
	failMessage   string
}

type inputValidatorFunc func(string) error

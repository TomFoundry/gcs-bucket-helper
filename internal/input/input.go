package input

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

// Recv ...
func Recv(promptMessage string, validators ...Validator) string {

	// Every line should have prefix "- "
	promptMessage = "- " + promptMessage
	promptMessage = strings.Replace(promptMessage, "\n", "\n- ", -1)

	fmt.Printf("%s\n", promptMessage)

OUTER:
	for {
		s, _ := in.ReadString('\n')

		s = strings.TrimSuffix(s, "\n")

		for _, validator := range validators {
			if err := validator(s); err != nil {
				fmt.Println(err)
				fmt.Println("Please try again")
				continue OUTER
			}
		}

		return s
	}
}

// Validator ...
type Validator func(string) error

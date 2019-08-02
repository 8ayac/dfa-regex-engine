package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/8ayac/dfa-regex-engine/dfaregex"
)

func main() {
	fmt.Printf("[\x1b[31m+\x1b[0m] input regex: ")
	regex := bufio.NewScanner(os.Stdin)
	regex.Scan()
	re := dfaregex.Compile(regex.Text())

	for {
		fmt.Printf("[\x1b[33m+\x1b[0m] input string to match: ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()

		if s.Text() == "<EXIT>" {
			fmt.Printf("Bye:)\n")
			os.Exit(0)
		}

		if s.Text() == "<REGEX>" {
			fmt.Printf("[\x1b[34m-\x1b[0m] input new regex: ")
			regex.Scan()
			re = dfaregex.Compile(regex.Text())
			fmt.Printf("\n")
			continue
		}

		if re.Match(s.Text()) {
			fmt.Printf("%s => \x1b[32mMatch!\x1b[0m\n", s.Text())
		} else {
			fmt.Printf("%s => \x1b[31mNot match.\x1b[0m\n", s.Text())
		}
		fmt.Printf("\n")
	}
}

package main

import (
	"fmt"
	"os"
	"bufio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}



func main() {
	
	args := os.Args[1:]
	
	if len(args) > 1 {
		fmt.Println("Usage: jlox [script]")
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(filename string) {
	fmt.Println(fmt.Sprintf("Opening file %v ...", filename))
	dat, err := os.ReadFile(filename);
	check(err)
	
	fmt.Printf("%d bytes read from file\n", len(dat))
}

func runPrompt() {
	fmt.Println("Running prompt")
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			fmt.Println("Exiting prompt...")
			return
		}
		run(line)
	}
}

func run(source string) {
	tokens, hasError, scanError := scanTokens(source)
	
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	
	if hasError {
		msg := fmt.Sprintf("Bad input with error:\n%v", scanError.Error())
		fmt.Println(msg)
	}
}

package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to the mastermind game !")
	fmt.Println("The computer will cose a combination of 4 numbers, each between 1 and 6.")
	fmt.Println("You have 12 tries to guess !")
	fmt.Println("...")
	solution := GenerateCode()
	fmt.Println("The Code has been generated !")
	try := 1
	for try <= 12 {
		fmt.Printf("Try %d:\n", try)
		code, err := ScanCode(os.Stdin)
		if err != nil {
			fmt.Printf("Ooops we didn't understand your answer ! (%s)\n", err.Error())
			fmt.Println("You can ty again !")
			continue
		}
		if code == solution {
			fmt.Printf("Bingo ! You won ! (in %d tries)\n", try)
			break
		}
		good, misplaced := CompareCode(code, solution)
		fmt.Println("Ooops ! You didn't find the solution !")
		fmt.Printf("Good number(s): %d, Misplaced number(s): %d\n", good, misplaced)

		try++
	}
	fmt.Println("You have lost !")
	fmt.Println("the Code was: ")
}

// Code represent a mmastermind code : 4 integers. (that should be between 1 and 6)
type Code [4]int

// ScanCode scan a code fro an io.Reader
func ScanCode(r io.Reader) (Code, error) {
	var s string
	_, err := fmt.Fscanln(r, &s)
	if err != nil {
		return Code{1, 1, 1, 1}, err
	}

	reg, err := regexp.Compile(`([1-6]){1},?`)
	if err != nil {
		return Code{1, 1, 1, 1}, err
	}
	ints := reg.FindAllString(s, -1)
	if len(ints) != 4 {
		return Code{1, 1, 1, 1}, errors.New("wrong input format ! Your input must be like this: '1,1,1,1', with numbers between 1 and 6")
	}

	var code Code
	for k, i := range ints {
		c, err := strconv.Atoi(strings.TrimSuffix(i, ","))
		if err != nil {
			return Code{1, 1, 1, 1}, err
		}
		code[k] = c
	}
	return code, nil
}

// GenerateCode generate a random code with int between 1 and 6.
func GenerateCode() Code {
	var c Code
	for k := range c {
		c[k] = rand.Intn(6) + 1
	}
	return c
}

// CompareCode compare 2 codes and give you how many guess are good, and how many are misplaced.
func CompareCode(try, solution Code) (good, misplaced int) {
	for k, v := range try {
		if solution[k] == v {
			good++
			// flag matching number so we won't compare them again
			try[k] = -1
			solution[k] = -1
			continue
		}
	}
MisplacedLoop:
	for k, v := range try {
		if v == -1 {
			continue
		}
		for k2, v2 := range solution {
			if k != k2 && v == v2 {
				misplaced++
				// flag matching number so we won't compare them again
				try[k] = -1
				solution[k2] = -1
				continue MisplacedLoop
			}
		}
	}
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

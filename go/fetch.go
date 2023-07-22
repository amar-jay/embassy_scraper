package main

import (
	"regexp"
	"strings"
)

// get all phone numbers from text
func Phone(text string) []string {
	var numbers []string
	// Sample text

	// Regular expression to match phone numbers with optional brackets and spaces
	regex := regexp.MustCompile(`\+?[0-9 ()]+`)

	// Find all matches in the text
	matches := regex.FindAllString(text, -1)

	// Loop over the matches and print them without brackets or spaces
	for _, match := range matches {
		match = strings.Trim(match, " ")
		delim := []byte{'(', ')', ' '}
		for _, d := range delim {
			match = strings.Replace(match, string(d), "", -1)
		}

		if len(match) < 10 {
			continue
		}

		numbers = append(numbers, match)
	}
	return numbers
}

// get email from text
func email(text string) []string {
	var addrs []string

	regex := regexp.MustCompile(`[a-zA-Z.]+@[a-zA-Z.]{3,}`)
	matches := regex.FindAllString(text, -1)

	for _, match := range matches {
		match = strings.Trim(match, " ")
		if len(match) < 5 {
			continue
		}

		addrs = append(addrs, match)
	}
	return addrs
}

// get address from text
func Address(text string) string {
	var addr string
	regex := regexp.MustCompile("Chancery:\n.*\n")
	addr = regex.FindString(text)
	addr = strings.Trim(addr, "Chancery:\n")
	addr = strings.Trim(addr, "\n ")
	return addr
}

// get ambassador name from text
func Ambassador(text string) string {
	var name string
	regex := regexp.MustCompile("Ambassador.*\n")
	name = regex.FindString(text)
	name = strings.Trim(name, "Ambassador\n")
	name = strings.Trim(name, "\n ")
	name = strings.ToTitle(name)
	return name
}

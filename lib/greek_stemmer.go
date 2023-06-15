package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Config struct {
	Step1Exceptions map[string]string "json:\"step_1_exceptions\""
	Step0Exceptions map[string]string "json:\"step_0_exceptions\""
	ProtectedWords  []string          "json:\"protected_words\""
}

func parseConfigFile(filepath string) (Config, error) {
	jsonFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to read the JSON file: %v", err)
	}
	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to unmarshal the JSON data: %v", err)
	}
	return config, err
}

func isgreek(x string) bool {
	alphabet := regexp.MustCompile("^[ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ]+$")
	return alphabet.MatchString(x)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Input: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	return strings.ReplaceAll(input, "\n", "")
}

func step1Regex(step1Exceptions map[string]string) (*regexp.Regexp, error) {
	suffixes := make([]string, 0, len(step1Exceptions))
	for suffix := range step1Exceptions {
		suffixes = append(suffixes, suffix)
	}
	pattern := "(.*)(?:" + strings.Join(suffixes, "|") + ")$"
	return regexp.Compile(pattern)
}

func main() {
	config, err := parseConfigFile("config/config.json")
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	word := getUserInput()
	stem := strings.Clone(word)

	if !isgreek(word) {
		log.Fatalf("Word is not valid: %v", word)
	}

	protected := config.ProtectedWords
	step0 := config.Step0Exceptions
	step1 := config.Step1Exceptions

    // If the word has less than 3 characters or is protected, return the word
	if utf8.RuneCountInString(stem) < 3 || contains(protected, stem) {
		fmt.Println("Stemmed word:", stem)
		os.Exit(0)
	}

    // Step 0 
	if val, ok := step0[stem]; ok {
		fmt.Println("Stemmed word:", val)
		os.Exit(0)
	}

	// Step 1
	step1Pattern, err := step1Regex(step1)
	if err != nil {
		log.Fatalf("Failed to compile step1 regular expression: %v", err)
	}

	stem = step1Pattern.ReplaceAllStringFunc(stem, func(match string) string {
		suffix := match
		replacement := step1[suffix]
		return replacement
	})

    // Step 2


	fmt.Println("Stemmed word: ", stem)
}

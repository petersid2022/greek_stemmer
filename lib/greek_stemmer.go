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
)

type Config struct {
	Step1Exceptions map[string]string "json:\"step_1_exceptions\""
	Step0Exceptions map[string]string "json:\"step_0_exceptions\""
	ProtectedWords  []string          "json:\"protected_words\""
}

func parseConfigFile(filepath string) Config {
	jsonFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read the JSON file: %v", err)
	}
	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal the JSON data: %v", err)
	}
	return config
}

func isgreek(x string) bool {
	alphabet := regexp.MustCompile("^[ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ]+$")
	return alphabet.MatchString(x)
}

func contains(word string, pool []string) bool {
	for _, w := range pool {
		if word == w {
			return true
		}
	}
	return false
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Input Greek capital letters only! : ")
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
	step2A_pattern := "(.*)(?:" + strings.Join(suffixes, "|") + ")$"
	return regexp.Compile(step2A_pattern)
}

func getMatchesFromInput(pattern string, stem string) []string {
    r := regexp.MustCompile(pattern)
	matches := r.FindStringSubmatch(stem)
    return matches
}

func main() {
	config := parseConfigFile("config/config.json")
	word := getUserInput()
    if len(word) < 3 && isgreek(word){
        fmt.Println(word)
        os.Exit(0)
    }

    stem := strings.Clone(word)
    protected := config.ProtectedWords
    step1 := config.Step1Exceptions
    step0 := config.Step0Exceptions

	if contains(stem, protected) || !isgreek(stem) {
		stem = "input word is protected or is not greek!"
	}

	// Step 0
	if val, ok := step0[stem]; ok {
		stem = val
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

	// Step 2a
	step2A_pattern := "^(.+?)(ΑΔΕΣ|ΑΔΩΝ)$"
    matches_2a := getMatchesFromInput(step2A_pattern,stem)
	if len(matches_2a) > 0 {
		st := matches_2a[1]
		if matched := regexp.MustCompile(`((ΟΚ|ΜΑΜ|ΜΑΝ|ΜΠΑΜΠ|ΠΑΤΕΡ|ΓΙΑΓΙ|ΝΤΑΝΤ|ΚΥΡ|ΘΕΙ|ΠΕΘΕΡ|ΜΟΥΣΑΜ|ΚΑΠΛΑΜ|ΠΑΡ|ΨΑΡ|ΤΖΟΥΡ|ΤΑΜΠΟΥΡ))$`).MatchString(st); matched {
			stem = st + "ΑΔ"
		}
	}

    // Step 2b
    step2B_pattern := "^(.+?)(ΕΔΕΣ|ΕΔΩΝ)$" 
    matches_2b := getMatchesFromInput(step2B_pattern,stem)
	if len(matches_2b) > 0 {
		st := matches_2b[1]
		if matched := regexp.MustCompile(`((ΟΠ|ΙΠ|ΕΜΠ|ΥΠ|ΓΗΠ|ΔΑΠ|ΚΡΑΣΠ|ΜΙΛ))$`).MatchString(st); matched {
			stem = st + "EΔ"
		}
	}

    // Step 2c
    step2C_pattern := "^(.+?)(ΟΥΔΕΣ|ΟΥΔΩΝ)$" 
    matches_2c := getMatchesFromInput(step2C_pattern,stem)
	if len(matches_2c) > 0 {
		st := matches_2c[1]
		if matched := regexp.MustCompile(`((ΑΡΚ|ΚΑΛΙΑΚ|ΠΕΤΑΛ|ΛΙΧ|ΠΛΕΞ|ΣΚ|Σ|ΦΛ|ΦΡ|ΒΕΛ|ΛΟΥΛ|ΧΝ|ΣΠ|ΤΡΑΓ|ΦΕ))$`).MatchString(st); matched {
			stem = st + "ΟΥΔ"
		}
	}

    // Step 2d
    step2D_pattern := "^(.+?)(ΕΩΣ|ΕΩΝ|ΕΑΣ|ΕΑ)$" 
    matches_2d := getMatchesFromInput(step2D_pattern,stem)
	if len(matches_2d) > 0 {
		st := matches_2d[1]
		if matched := regexp.MustCompile(`(^(Θ|Δ|ΕΛ|ΓΑΛ|Ν|Π|ΙΔ|ΠΑΡ|ΣΤΕΡ|ΟΡΦ|ΑΝΔΡ|ΑΝΤΡ))$`).MatchString(st); matched {
			stem = st + "Ε"
		}
	}

    // Step 3a 
    //step3A_pattern := "^(.+?)(ΕΙΟ|ΕΙΟΣ|ΕΙΟΙ|ΕΙΑ|ΕΙΑΣ|ΕΙΕΣ|ΕΙΟΥ|ΕΙΟΥΣ|ΕΙΩΝ)$"
    //matches_3a := getMatchesFromInput(step3A_pattern, stem)
    //st := matches_3a[1]
    //if len(st) > 4 {
    //    stem = st
    //}

    // Step 3b
    //step3B_pattern := "^(.+?)(ΙΟΥΣ|ΙΑΣ|ΙΕΣ|ΙΟΣ|ΙΟΥ|ΙΟΙ|ΙΩΝ|ΙΟΝ|ΙΑ|ΙΟ)$"
    //matches_3b := getMatchesFromInput(step3B_pattern, stem)

    fmt.Println(stem)
}

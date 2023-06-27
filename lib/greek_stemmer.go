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

func ends_on_vowel(word string) bool {
	pattern := "[ΑΕΗΙΟΥΩ]$"
	match, _ := regexp.MatchString(pattern, word)

	if match {
		return true
	}

	return false
}

func ends_on_vowel2(word string) bool {
	pattern := "[ΑΕΗΙΟΩ]$"
	match, _ := regexp.MatchString(pattern, word)

	if match {
		return true
	}

	return false
}

func main() {
	config := parseConfigFile("config/config.json")
	word := getUserInput()

	if len(word) < 3 && isgreek(word) {
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
	matches_2A := getMatchesFromInput(step2A_pattern, stem)
	if len(matches_2A) > 0 {
		st := matches_2A[1]
		if matched := regexp.MustCompile(`((ΟΚ|ΜΑΜ|ΜΑΝ|ΜΠΑΜΠ|ΠΑΤΕΡ|ΓΙΑΓΙ|ΝΤΑΝΤ|ΚΥΡ|ΘΕΙ|ΠΕΘΕΡ|ΜΟΥΣΑΜ|ΚΑΠΛΑΜ|ΠΑΡ|ΨΑΡ|ΤΖΟΥΡ|ΤΑΜΠΟΥΡ))$`).MatchString(st); matched {
			stem = st + "ΑΔ"
		}
	}

	// Step 2b
	step2B_pattern := "^(.+?)(ΕΔΕΣ|ΕΔΩΝ)$"
	matches_2B := getMatchesFromInput(step2B_pattern, stem)
	if len(matches_2B) > 0 {
		st := matches_2B[1]
		if matched := regexp.MustCompile(`((ΟΠ|ΙΠ|ΕΜΠ|ΥΠ|ΓΗΠ|ΔΑΠ|ΚΡΑΣΠ|ΜΙΛ))$`).MatchString(st); matched {
			stem = st + "EΔ"
		}
	}

	// Step 2c
	step2C_pattern := "^(.+?)(ΟΥΔΕΣ|ΟΥΔΩΝ)$"
	matches_2C := getMatchesFromInput(step2C_pattern, stem)
	if len(matches_2C) > 0 {
		st := matches_2C[1]
		if matched := regexp.MustCompile(`((ΑΡΚ|ΚΑΛΙΑΚ|ΠΕΤΑΛ|ΛΙΧ|ΠΛΕΞ|ΣΚ|Σ|ΦΛ|ΦΡ|ΒΕΛ|ΛΟΥΛ|ΧΝ|ΣΠ|ΤΡΑΓ|ΦΕ))$`).MatchString(st); matched {
			stem = st + "ΟΥΔ"
		}
	}

	// Step 2d
	step2D_pattern := "^(.+?)(ΕΩΣ|ΕΩΝ|ΕΑΣ|ΕΑ)$"
	matches_2D := getMatchesFromInput(step2D_pattern, stem)
	if len(matches_2D) > 0 {
		st := matches_2D[1]
		if matched := regexp.MustCompile(`(^(Θ|Δ|ΕΛ|ΓΑΛ|Ν|Π|ΙΔ|ΠΑΡ|ΣΤΕΡ|ΟΡΦ|ΑΝΔΡ|ΑΝΤΡ))$`).MatchString(st); matched {
			stem = st + "Ε"
		}
	}

	// Step 3a
	step3A_pattern := "^(.+?)(ΕΙΟ|ΕΙΟΣ|ΕΙΟΙ|ΕΙΑ|ΕΙΑΣ|ΕΙΕΣ|ΕΙΟΥ|ΕΙΟΥΣ|ΕΙΩΝ)$"
	matches_3A := getMatchesFromInput(step3A_pattern, stem)
	if len(matches_3A) > 0 {
		st := matches_3A[1]
		if len(st) > 4 {
			stem = st
		}
	}

	// Step 3b
	step3B_pattern := "^(.+?)(ΙΟΥΣ|ΙΑΣ|ΙΕΣ|ΙΟΣ|ΙΟΥ|ΙΟΙ|ΙΩΝ|ΙΟΝ|ΙΑ|ΙΟ)$"
	matches_3B := getMatchesFromInput(step3B_pattern, stem)
	if len(matches_3B) > 0 {
		st := matches_3B[1]
		big_regex := `^(ΑΓ|ΑΓΓΕΛ|ΑΓΡ|
                     ΑΕΡ|ΑΘΛ|ΑΚΟΥΣ|ΑΞ|ΑΣ|Β|ΒΙΒΛ|ΒΥΤ|Γ|ΓΙΑΓ|ΓΩΝ|Δ|ΔΑΝ|ΔΗΛ|ΔΗΜ|
                     ΔΟΚΙΜ|ΕΛ|ΖΑΧΑΡ|ΗΛ|ΗΠ|ΙΔ|ΙΣΚ|ΙΣΤ|ΙΟΝ|ΙΩΝ|ΚΙΜΩΛ|ΚΟΛΟΝ|ΚΟΡ|
                     ΚΤΗΡ|ΚΥΡ|ΛΑΓ|ΛΟΓ|ΜΑΓ|ΜΠΑΝ|ΜΠΕΤΟΝ|ΜΠΡ|ΝΑΥΤ|ΝΟΤ|ΟΠΑΛ|ΟΞ|ΟΡ|ΟΣ|
                     ΠΑΝ|ΠΑΝΑΓ|ΠΑΤΡ|ΠΗΛ|ΠΗΝ|ΠΛΑΙΣ|ΠΟΝΤ|ΡΑΔ|ΡΟΔ|ΣΚ|ΣΚΟΡΠ|ΣΟΥΝ|ΣΠΑΝ|
                     ΣΤΑΔ|ΣΥΡ|ΤΗΛ|ΤΕΤΡΑΔ|ΤΙΜ|ΤΟΚ|ΤΟΠ|ΤΡΟΧ|ΦΙΛ|ΦΩΤ|Χ|ΧΙΛ|ΧΡΩΜ|ΧΩΡ)$`
		match1, _ := regexp.MatchString(big_regex, st)
		if len(st) < 2 || ends_on_vowel(st) || match1 {
			stem = st + "Ι"
		}
		match2, _ := regexp.MatchString("^(ΠΑΛ)$", st)
		if match2 {
			stem = st + "ΑΙ"
		}
	}

	// Step 4
	step4_pattern := "^(.+?)(ΙΚΟΣ|ΙΚΟΝ|ΙΚΕΙΣ|ΙΚΟΙ|ΙΚΕΣ|ΙΚΟΥΣ|ΙΚΗ|ΙΚΗΣ|ΙΚΟ|ΙΚΑ|ΙΚΟΥ|ΙΚΩΝ|ΙΚΩΣ)$"
	matches_4 := getMatchesFromInput(step4_pattern, stem)
	if len(matches_4) > 0 {
		st := matches_4[1]
		big_regex := `^(ΑΔ|ΑΛ|ΑΜΑΝ|ΑΜΕΡ|ΑΜΜΟΧΑΛ|
                    ΑΝΗΘ|ΑΝΤΙΔ|ΑΠΛ|ΑΤΤ|ΑΦΡ|ΒΑΣ|ΒΡΩΜ|ΓΕΝ|ΓΕΡ|Δ|ΔΙΑΦΟΡ|ΔΙΚΑΝ|
                    ΔΥΤ|ΕΙΔ|ΕΝΔ|ΕΞΩΔ|ΗΘ|ΘΕΤ|ΚΑΛΛΙΝ|ΚΑΛΠ|ΚΑΤΑΔ|ΚΟΥΖΙΝ|ΚΡ|ΚΩΔ|
                    ΛΑΔ|ΛΟΓ|Μ|ΜΕΡ|ΜΟΝΑΔ|ΜΟΥΛ|ΜΟΥΣ|ΜΠΑΓΙΑΤ|ΜΠΑΝ|ΜΠΟΛ|ΜΠΟΣ|ΜΥΣΤ|
                    Ν|ΝΙΤ|ΞΙΚ|ΟΠΤ|ΠΑΝ|ΠΕΡΙΣΤΡΟΦ|ΠΕΤΣ|ΠΙΚΑΝΤ|ΠΙΤΣ|ΠΛΑΣΤ|ΠΛΙΑΤΣ|
                    ΠΟΝΤ|ΠΟΣΤΕΛΝ|ΠΡΩΤΟΔ|ΣΕΡΤ|ΣΗΜΑΝΤ|ΣΤΑΤ|ΣΥΝΑΔ|ΣΥΝΟΜΗΛ|ΤΕΛ|
                    ΤΕΧΝ|ΤΗΛΕΣΚΟΠ|ΤΡΟΠ|ΤΣΑΜ|ΥΠΟΔ|Φ|ΦΙΛΟΝ|ΦΥΛΟΔ|ΦΥΣ|ΧΑΣ|ΦΥΤ)
                    $`
		match1, _ := regexp.MatchString(big_regex, st)
		match2, _ := regexp.MatchString("(ΦΟΙΝ)$", st)
		if ends_on_vowel(st) || match1 || match2 {
			stem = st + "ΙΚ"
		} else if stem == "ΠΑΣΧΑΛΙΑΤ" {
			stem = "ΠΑΣΧΑ"
		}
	}

	// Step 5a
	if stem == "ΑΓΑΜΕ" {
		stem = "ΑΓΑΜ"
	}

	step5A_pattern := "^(.+?)(ΑΓΑΜΕ|ΗΣΑΜΕ|ΟΥΣΑΜΕ|ΗΚΑΜΕ|ΗΘΗΚΑΜΕ)$"
	matches_5A := getMatchesFromInput(step5A_pattern, stem)
	if len(matches_5A) > 0 {
		st := matches_5A[1]
		stem = st
	}

	step5A_2_pattern := "^(.+?)(ΑΜΕ)$"
	matches_5A_2 := getMatchesFromInput(step5A_2_pattern, stem)
	if len(matches_5A_2) > 0 {
		st := matches_5A_2[1]
		if matched := regexp.MustCompile(`(^(ΑΝΑΠ|ΑΠΟΘ|ΑΠΟΚ|ΑΠΟΣΤ|ΒΟΥΒ|ΞΕΘ|ΟΥΛ|ΠΕΘ|ΠΙΚΡ|ΠΟΤ|ΣΙΧ|Χ)$)$`).MatchString(st); matched {
			stem = st + "ΑΜ"
		}
	}

	// Step 5b
	step5B_pattern := "^(.+?)(ΑΓΑΝΕ|ΗΣΑΝΕ|ΟΥΣΑΝΕ|ΙΟΝΤΑΝΕ|ΙΟΤΑΝΕ|ΙΟΥΝΤΑΝΕ|ΟΝΤΑΝΕ|ΟΤΑΝΕ|ΟΥΝΤΑΝΕ|ΗΚΑΝΕ|ΗΘΗΚΑΝΕ)$"
	matches_5B := getMatchesFromInput(step5B_pattern, stem)
	if len(matches_5B) > 0 {
		st := matches_5B[1]
		if matched := regexp.MustCompile(`(^(ΤΡ|ΤΣ)$)$`).MatchString(st); matched {
			stem = st + "ΑΓΑΝ"
		}
	}

	step5B_2_pattern := "^(.+?)(ΑΝΕ)$"
	matches_5B_2 := getMatchesFromInput(step5B_2_pattern, stem)
	if len(matches_5B_2) > 0 {
		st := matches_5B_2[1]
		if matched := regexp.MustCompile(`(^(ΒΕΤΕΡ|ΒΟΥΛΚ|ΒΡΑΧΜ|Γ|ΔΡΑΔΟΥΜ|Θ|ΚΑΛΠΟΥΖ|ΚΑΣΤΕΛ|
	                           ΚΟΡΜΟΡ|ΛΑΟΠΛ|ΜΩΑΜΕΘ|Μ|ΜΟΥΣΟΥΛΜ|Ν|ΟΥΛ|Π|ΠΕΛΕΚ|
	                           ΠΛ|ΠΟΛΙΣ|ΠΟΡΤΟΛ|ΣΑΡΑΚΑΤΣ|ΣΟΥΛΤ|ΤΣΑΡΛΑΤ|ΟΡΦ|ΤΣΙΓΓ|
	                           ΤΣΟΠ|ΦΩΤΟΣΤΕΦ|Χ|ΨΥΧΟΠΛ|ΑΓ|ΟΡΦ|ΓΑΛ|ΓΕΡ|ΔΕΚ|ΔΙΠΛ|
	                           ΑΜΕΡΙΚΑΝ|ΟΥΡ|ΠΙΘ|ΠΟΥΡΙΤ|Σ|ΖΩΝΤ|ΙΚ|ΚΑΣΤ|ΚΟΠ|ΛΙΧ|
	                           ΛΟΥΘΗΡ|ΜΑΙΝΤ|ΜΕΛ|ΣΙΓ|ΣΠ|ΣΤΕΓ|ΤΡΑΓ|ΤΣΑΓ|Φ|ΕΡ|ΑΔΑΠ|
	                           ΑΘΙΓΓ|ΑΜΗΧ|ΑΝΙΚ|ΑΝΟΡΓ|ΑΠΗΓ|ΑΠΙΘ|ΑΤΣΙΓΓ|ΒΑΣ|ΒΑΣΚ|
	                           ΒΑΘΥΓΑΛ|ΒΙΟΜΗΧ|ΒΡΑΧΥΚ|ΔΙΑΤ|ΔΙΑΦ|ΕΝΟΡΓ|ΘΥΣ|
	                           ΚΑΠΝΟΒΙΟΜΗΧ|ΚΑΤΑΓΑΛ|ΚΛΙΒ|ΚΟΙΛΑΡΦ|ΛΙΒ|ΜΕΓΛΟΒΙΟΜΗΧ|
	                           ΜΙΚΡΟΒΙΟΜΗΧ|ΝΤΑΒ|ΞΗΡΟΚΛΙΒ|ΟΛΙΓΟΔΑΜ|ΟΛΟΓΑΛ|ΠΕΝΤΑΡΦ|
	                           ΠΕΡΗΦ|ΠΕΡΙΤΡ|ΠΛΑΤ|ΠΟΛΥΔΑΠ|ΠΟΛΥΜΗΧ|ΣΤΕΦ|ΤΑΒ|ΤΕΤ|
	                           ΥΠΕΡΗΦ|ΥΠΟΚΟΠ|ΧΑΜΗΛΟΔΑΠ|ΨΗΛΟΤΑΒ)$)$`).MatchString(st); matched || ends_on_vowel2(st){
			stem = st + "ΑΝ"
		}
	}


    // Step 5c
    

	fmt.Println(stem)
}

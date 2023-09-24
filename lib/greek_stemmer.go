package greek_stemmer_go 

import (
    "encoding/json"
	"fmt"
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

func getFilePathContents() Config{
    filePath := "../config/config.json"

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

    // Create a Config struct to unmarshal the JSON data into
	var config Config

	// Unmarshal the JSON data into the Config struct
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func longStemList(word string) string {
	stepPattern := `^(.+?)(Α|ΑΓΑΤΕ|ΑΓΑΝ|ΑΕΙ|ΑΜΑΙ|ΑΝ|ΑΣ|ΑΣΑΙ|ΑΤΑΙ|ΑΩ|Ε|ΕΙ|ΕΙΣ|ΕΙΤΕ|
    ΕΣΑΙ|ΕΣ|ΕΤΑΙ|Ι|ΙΕΜΑΙ|ΙΕΜΑΣΤΕ|ΙΕΤΑΙ|ΙΕΣΑΙ|ΙΕΣΑΣΤΕ|ΙΟΜΑΣΤΑΝ|ΙΟΜΟΥΝ|ΙΟΜΟΥΝΑ|
    ΙΟΝΤΑΝ|ΙΟΝΤΟΥΣΑΝ|ΙΟΣΑΣΤΑΝ|ΙΟΣΑΣΤΕ|ΙΟΣΟΥΝ|ΙΟΣΟΥΝΑ|ΙΟΤΑΝ|ΙΟΥΜΑ|ΙΟΥΜΑΣΤΕ|
    ΙΟΥΝΤΑΙ|ΙΟΥΝΤΑΝ|Η|ΗΔΕΣ|ΗΔΩΝ|ΗΘΕΙ|ΗΘΕΙΣ|ΗΘΕΙΤΕ|ΗΘΗΚΑΤΕ|ΗΘΗΚΑΝ|ΗΘΟΥΝ|ΗΘΩ|
    ΗΚΑΤΕ|ΗΚΑΝ|ΗΣ|ΗΣΑΝ|ΗΣΑΤΕ|ΗΣΕΙ|ΗΣΕΣ|ΗΣΟΥΝ|ΗΣΩ|Ο|ΟΙ|ΟΜΑΙ|ΟΜΑΣΤΑΝ|ΟΜΟΥΝ|ΟΜΟΥΝΑ|
    ΟΝΤΑΙ|ΟΝΤΑΝ|ΟΝΤΟΥΣΑΝ|ΟΣ|ΟΣΑΣΤΑΝ|ΟΣΑΣΤΕ|ΟΣΟΥΝ|ΟΣΟΥΝΑ|ΟΤΑΝ|ΟΥ|ΟΥΜΑΙ|ΟΥΜΑΣΤΕ|
    ΟΥΝ|ΟΥΝΤΑΙ|ΟΥΝΤΑΝ|ΟΥΣ|ΟΥΣΑΝ|ΟΥΣΑΤΕ|Υ||ΥΑ|ΥΣ|Ω|ΩΝ|ΟΙΣ)$`
	r := regexp.MustCompile(stepPattern)
	matches := r.FindAllStringSubmatch(word, -1)

	if matches == nil {
		return word
	}

	for _, match := range matches {
		st := match[1]
		suffix := match[2]

		if word == "ΠΑΣΧΑ" {
			return word
		}

		word = st

		if st == "ΣΠΟΡ" && suffix != "" {
			word += "Ο"
		} else if st == "ΠΑΣΧΑΛΙΝ" && suffix != "" {
			word = "ΠΑΣΧΑ"
		}
	}

	return word
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

func GreekStemmer(word string) string{
    config := getFilePathContents()

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
        stem = st
		if matched := regexp.MustCompile(`((ΟΚ|ΜΑΜ|ΜΑΝ|ΜΠΑΜΠ|ΠΑΤΕΡ|ΝΤΑΝΤ|ΚΥΡ|ΘΕΙ|ΠΕΘΕΡ|ΜΟΥΣΑΜ|ΚΑΠΛΑΜ|ΠΑΡ|ΨΑΡ|ΤΖΟΥΡ|ΤΑΜΠΟΥΡ))$`).MatchString(st); matched {
			stem += "ΑΔ"
		}
	}

	// Step 2b
	step2B_pattern := "^(.+?)(ΕΔΕΣ|ΕΔΩΝ)$"
	matches_2B := getMatchesFromInput(step2B_pattern, stem)
	if len(matches_2B) > 0 {
		st := matches_2B[1]
        stem = st
		if matched := regexp.MustCompile(`((ΟΠ|ΙΠ|ΕΜΠ|ΥΠ|ΓΗΠ|ΔΑΠ|ΚΡΑΣΠ|ΜΙΛ))$`).MatchString(st); matched {
			stem += "EΔ"
		}
	}

	// Step 2c
	step2C_pattern := "^(.+?)(ΟΥΔΕΣ|ΟΥΔΩΝ)$"
	matches_2C := getMatchesFromInput(step2C_pattern, stem)
	if len(matches_2C) > 0 {
		st := matches_2C[1]
        stem = st
		if matched := regexp.MustCompile(`((ΑΡΚ|ΚΑΛΙΑΚ|ΠΕΤΑΛ|ΛΙΧ|ΠΛΕΞ|ΣΚ|Σ|ΦΛ|ΦΡ|ΒΕΛ|ΛΟΥΛ|ΧΝ|ΣΠ|ΤΡΑΓ|ΦΕ))$`).MatchString(st); matched {
            stem += "ΟΥΔ"
        }
	}

	// Step 2d
	step2D_pattern := "^(.+?)(ΕΩΣ|ΕΩΝ|ΕΑΣ|ΕΑ)$"
	matches_2D := getMatchesFromInput(step2D_pattern, stem)
	if len(matches_2D) > 0 {
		st := matches_2D[1]
        stem = st
		if matched := regexp.MustCompile(`(^(Θ|Δ|ΕΛ|ΓΑΛ|Ν|Π|ΙΔ|ΠΑΡ|ΣΤΕΡ|ΟΡΦ|ΑΝΔΡ|ΑΝΤΡ))$`).MatchString(st); matched {
			stem += "Ε"
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
        stem = st
		big_regex := `^(ΑΓ|ΑΓΓΕΛ|ΑΓΡ|
                     ΑΕΡ|ΑΘΛ|ΑΚΟΥΣ|ΑΞ|ΑΣ|Β|ΒΙΒΛ|ΒΥΤ|Γ|ΓΙΑΓ|ΓΩΝ|Δ|ΔΑΝ|ΔΗΛ|ΔΗΜ|
                     ΔΟΚΙΜ|ΕΛ|ΖΑΧΑΡ|ΗΛ|ΗΠ|ΙΔ|ΙΣΚ|ΙΣΤ|ΙΟΝ|ΙΩΝ|ΚΙΜΩΛ|ΚΟΛΟΝ|ΚΟΡ|
                     ΚΤΗΡ|ΚΥΡ|ΛΑΓ|ΛΟΓ|ΜΑΓ|ΜΠΑΝ|ΜΠΕΤΟΝ|ΜΠΡ|ΝΑΥΤ|ΝΟΤ|ΟΠΑΛ|ΟΞ|ΟΡ|ΟΣ|
                     ΠΑΝ|ΠΑΝΑΓ|ΠΑΤΡ|ΠΗΛ|ΠΗΝ|ΠΛΑΙΣ|ΠΟΝΤ|ΡΑΔ|ΡΟΔ|ΣΚ|ΣΚΟΡΠ|ΣΟΥΝ|ΣΠΑΝ|
                     ΣΤΑΔ|ΣΥΡ|ΤΗΛ|ΤΕΤΡΑΔ|ΤΙΜ|ΤΟΚ|ΤΟΠ|ΤΡΟΧ|ΦΙΛ|ΦΩΤ|Χ|ΧΙΛ|ΧΡΩΜ|ΧΩΡ)$`
		match1, _ := regexp.MatchString(big_regex, st)
		if len(st) < 2 || ends_on_vowel(st) || match1 {
			stem += "Ι"
		}
		match2, _ := regexp.MatchString("^(ΠΑΛ)$", st)
		if match2 {
			stem += "ΑΙ"
		}
	}

	// Step 4
	step4_pattern := "^(.+?)(ΙΚΟΣ|ΙΚΟΝ|ΙΚΕΙΣ|ΙΚΟΙ|ΙΚΕΣ|ΙΚΟΥΣ|ΙΚΗ|ΙΚΗΣ|ΙΚΟ|ΙΚΑ|ΙΚΟΥ|ΙΚΩΝ|ΙΚΩΣ)$"
	matches_4 := getMatchesFromInput(step4_pattern, stem)
	if len(matches_4) > 0 {
		st := matches_4[1]
        stem = st
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
			stem += "ΙΚ"
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
        stem = st
		if matched := regexp.MustCompile(`(^(ΑΝΑΠ|ΑΠΟΘ|ΑΠΟΚ|ΑΠΟΣΤ|ΒΟΥΒ|ΞΕΘ|ΟΥΛ|ΠΕΘ|ΠΙΚΡ|ΠΟΤ|ΣΙΧ|Χ)$)$`).MatchString(st); matched {
			stem += "ΑΜ"
		}
	}

	// Step 5b
	step5B_pattern := "^(.+?)(ΑΓΑΝΕ|ΗΣΑΝΕ|ΟΥΣΑΝΕ|ΙΟΝΤΑΝΕ|ΙΟΤΑΝΕ|ΙΟΥΝΤΑΝΕ|ΟΝΤΑΝΕ|ΟΤΑΝΕ|ΟΥΝΤΑΝΕ|ΗΚΑΝΕ|ΗΘΗΚΑΝΕ)$"
	matches_5B := getMatchesFromInput(step5B_pattern, stem)
	if len(matches_5B) > 0 {
		st := matches_5B[1]
        stem = st
		if matched := regexp.MustCompile(`(^(ΤΡ|ΤΣ)$)$`).MatchString(st); matched {
			stem += "ΑΓΑΝ"
		}
	}

	step5B_2_pattern := "^(.+?)(ΑΝΕ)$"
	matches_5B_2 := getMatchesFromInput(step5B_2_pattern, stem)
	if len(matches_5B_2) > 0 {
		st := matches_5B_2[1]
        stem = st
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
	                           ΥΠΕΡΗΦ|ΥΠΟΚΟΠ|ΧΑΜΗΛΟΔΑΠ|ΨΗΛΟΤΑΒ)$)$`).MatchString(st); matched || ends_on_vowel2(st) {
			stem += "ΑΝ"
		}
	}

	// Step 5c
	step5C_pattern := "^(.+?)(ΗΣΕΤΕ)$"
	matches_5C := getMatchesFromInput(step5C_pattern, stem)
	if len(matches_5C) > 0 {
		st := matches_5C[1]
		stem = st
	}

	step5C_2_pattern := "^(.+?)(ΕΤΕ)$"
	matches_5C_2 := getMatchesFromInput(step5C_2_pattern, stem)
	if len(matches_5C_2) > 0 {
		st := matches_5C_2[1]
		stem = st
		if matched := regexp.MustCompile(`(^(ΟΔ|ΑΙΡ|ΦΟΡ|ΤΑΘ|ΔΙΑΘ|ΣΧ|ΕΝΔ|ΕΥΡ|ΤΙΘ|ΥΠΕΡΘ|ΡΑΘ|ΕΝΘ|ΡΟΘ|ΣΘ|ΠΥΡ|ΑΙΝ|ΣΥΝΔ|ΣΥΝ|ΣΥΝΘ|ΧΩΡ|ΠΟΝ|ΒΡ|ΚΑΘ|ΕΥΘ|ΕΚΘ|ΝΕΤ|ΡΟΝ|ΑΡΚ|ΒΑΡ|ΒΟΛ|ΩΦΕΛ|ΑΒΑΡ|ΒΕΝ|ΕΝΑΡ|ΑΒΡ|ΑΔ|ΑΘ|ΑΝ|ΑΠΛ|ΒΑΡΟΝ|ΝΤΡ|ΣΚ|ΚΟΠ|ΜΠΟΡ|ΝΙΦ|ΠΑΓ|ΠΑΡΑΚΑΛ|ΣΕΡΠ|ΣΚΕΛ|ΣΥΡΦ|ΤΟΚ|Υ|Δ|ΕΜ|ΘΑΡΡ|Θ|ΠΑΡΑΚΑΤΑΘ|ΠΡΟΣΘ|ΣΥΝΘ)$)$`).MatchString(st); matched || ends_on_vowel2(st) {
			stem += "ΕΤ"
		}
	}

	// Step 5d
	step5D_pattern := "^(.+?)(ΟΝΤΑΣ|ΩΝΤΑΣ)$"
	matches_5D := getMatchesFromInput(step5D_pattern, stem)
	if len(matches_5D) > 0 {
		st := matches_5D[1]
		stem = st
		if matched := regexp.MustCompile("(ΑΡΧ)$").MatchString(st); matched {
			stem += "ΟΝΤ"
		}
		if matched := regexp.MustCompile("(ΚΡΕ)$").MatchString(st); matched {
			stem += "ΩΝΤ"
		}
	}

	// Step 5e
	step5E_pattern := "^(.+?)(ΟΜΑΣΤΕ|ΙΟΜΑΣΤΕ)$"
	matches_5E := getMatchesFromInput(step5E_pattern, stem)
	if len(matches_5E) > 0 {
		st := matches_5E[1]
		stem = st
		if matched := regexp.MustCompile("(ΟΝ)$").MatchString(st); matched {
			stem += "ΟΜΑΣΤ"
		}
	}

	// Step 5f
	step5F_pattern := "^(.+?)(ΙΕΣΤΕ)$"
	matches_5F := getMatchesFromInput(step5F_pattern, stem)
	if len(matches_5F) > 0 {
		st := matches_5F[1]
		stem = st
		if matched := regexp.MustCompile(`(^(Π|ΑΠ|ΣΥΜΠ|ΑΣΥΜΠ|ΑΚΑΤΑΠ|ΑΜΕΤΑΜΦ)$)$`).MatchString(st); matched {
			stem += "IEΣΤ"
		}
	}

	step5F_2_pattern := "^(.+?)(ΕΣΤΕ)$"
	matches_5F_2 := getMatchesFromInput(step5F_2_pattern, stem)
	if len(matches_5F_2) > 0 {
		st := matches_5F_2[1]
		stem = st
		if matched := regexp.MustCompile(`(^(ΑΛ|ΑΡ|ΕΚΤΕΛ|Ζ|Μ|Ξ|ΠΑΡΑΚΑΛ|ΑΡ|ΠΡΟ|ΝΙΣ)$)$`).MatchString(st); matched {
			stem += "EΣΤ"
		}
	}

	// Step 5g
	step5G_pattern := "^(.+?)(ΗΘΗΚΑ|ΗΘΗΚΕΣ|ΗΘΗΚΕ)$"
	matches_5G := getMatchesFromInput(step5G_pattern, stem)
	if len(matches_5G) > 0 {
		st := matches_5G[1]
		stem = st
	}

	step5G_2_pattern := "^(.+?)(ΗΚΑ|ΗΚΕΣ|ΗΚΕ)$"
	matches_5G_2 := getMatchesFromInput(step5G_2_pattern, stem)
	if len(matches_5G_2) > 0 {
		st := matches_5G_2[1]
		stem = st
		if matched := regexp.MustCompile(`^(ΣΚΩΛ|ΣΚΟΥΛ|ΝΑΡΘ|ΣΦ|ΟΘ|ΠΙΘ)$`).MatchString(st); matched || regexp.MustCompile(`^(ΔΙΑΘ|Θ|ΠΑΡΑΚΑΤΑΘ|ΠΡΟΣΘ|ΣΥΝΘ|)$`).MatchString(st) {
			stem += "ΗΚ"
		}
	}

	// Step 5h
	step5H_pattern := "^(.+?)(ΟΥΣΑ|ΟΥΣΕΣ|ΟΥΣΕ)$"
	matches_5H := getMatchesFromInput(step5H_pattern, stem)
	if len(matches_5H) > 0 {
		st := matches_5H[1]
		stem = st
		if matched := regexp.MustCompile(`^(ΦΑΡΜΑΚ|ΧΑΔ|ΑΓΚ|ΑΝΑΡΡ|ΒΡΟΜ|ΕΚΛΙΠ|ΛΑΜΠΙΔ|ΛΕΧ|Μ|ΠΑΤ|Ρ|Λ|ΜΕΔ|ΜΕΣΑΖ|ΥΠΟΤΕΙΝ|
        ΑΜ|ΑΙΘ|ΑΝΗΚ|ΔΕΣΠΟΖ|ΕΝΔΙΑΦΕΡ|ΔΕ|ΔΕΥΤΕΡΕΥ|ΚΑΘΑΡΕΥ|ΠΛΕ|ΤΣΑ)$`).MatchString(st); matched || regexp.MustCompile(`^(ΠΟΔΑΡ|ΒΛΕΠ|
        ΠΑΝΤΑΧ|ΦΡΥΔ|ΜΑΝΤΙΛ|ΜΑΛΛ|ΚΥΜΑΤ|ΛΑΧ|ΛΗΓ|ΦΑΓ|ΟΜ|ΠΡΩΤ)$`).MatchString(st) || ends_on_vowel(st) {
			stem += "ΟΥΣ"
		}
	}

	// Step 5i
	step5I_pattern := "^(.+?)(ΑΓΑ|ΑΓΕΣ|ΑΓΕ)$"
	matches_5I := getMatchesFromInput(step5I_pattern, stem)
	if len(matches_5I) > 0 {
		st := matches_5I[1]
		stem = st
		if matched := regexp.MustCompile(`^(ΑΒΑΣΤ|ΠΟΛΥΦ|ΑΔΗΦ|ΠΑΜΦ|Ρ|ΑΣΠ|ΑΦ|ΑΜΑΛ|ΑΜΑΛΛΙ|
      ΑΝΥΣΤ|ΑΠΕΡ|ΑΣΠΑΡ|ΑΧΑΡ|ΔΕΡΒΕΝ|ΔΡΟΣΟΠ|ΞΕΦ|ΝΕΟΠ|ΝΟΜΟΤ|ΟΛΟΠ|ΟΜΟΤ|ΠΡΟΣΤ|
      ΠΡΟΣΩΠΟΠ|ΣΥΜΠ|ΣΥΝΤ|Τ|ΥΠΟΤ|ΧΑΡ|ΑΕΙΠ|ΑΙΜΟΣΤ|ΑΝΥΠ|ΑΠΟΤ|ΑΡΤΙΠ|ΔΙΑΤ|ΕΝ|ΕΠΙΤ|
      ΚΡΟΚΑΛΟΠ|ΣΙΔΗΡΟΠ|Λ|ΝΑΥ|ΟΥΛΑΜ|ΟΥΡ|Π|ΤΡ|Μ)$`).MatchString(st); matched || regexp.MustCompile(`^(ΟΦ|ΠΕΛ|ΧΟΡΤ|ΛΛ|ΣΦ|
      ΡΠ|ΦΡ|ΠΡ|ΛΟΧ|ΣΜΗΝ)$`).MatchString(st) && !(regexp.MustCompile(`^(ΨΟΦ|ΝΑΥΛΟΧ)$`).MatchString(st) || regexp.MustCompile("ΚΟΛΛ").MatchString(st)) {
			stem += "ΑΓ"
		}
	}

	// Step 5j
	step5J_pattern := "^(.+?)(ΗΣΕ|ΗΣΟΥ|ΗΣΑ)$"
	matches_5J := getMatchesFromInput(step5J_pattern, stem)
	if len(matches_5J) > 0 {
		st := matches_5J[1]
		stem = st
		if matched := regexp.MustCompile(`^(Ν|ΧΕΡΣΟΝ|ΔΩΔΕΚΑΝ|ΕΡΗΜΟΝ|ΜΕΓΑΛΟΝ|ΕΠΤΑΝ|Ι)$`).MatchString(st); matched {
			stem += "ΗΣ"
		}
	}

	// Step 5k
	step5K_pattern := "^(.+?)(ΗΣΤΕ)$"
	matches_5K := getMatchesFromInput(step5K_pattern, stem)
	if len(matches_5K) > 0 {
		st := matches_5K[1]
		stem = st
		if matched := regexp.MustCompile(`^(ΑΣΒ|ΣΒ|ΑΧΡ|ΧΡ|ΑΠΛ|ΑΕΙΜΝ|ΔΥΣΧΡ|ΕΥΧΡ|ΚΟΙΝΟΧΡ|
                             ΠΑΛΙΜΨ)$`).MatchString(st); matched {
			stem += "ΗΣΤ"
		}
	}

	// Step 5l
	step5L_pattern := "^(.+?)(ΟΥΝΕ|ΗΣΟΥΝΕ|ΗΘΟΥΝΕ)$"
	matches_5L := getMatchesFromInput(step5L_pattern, stem)
	if len(matches_5L) > 0 {
		st := matches_5L[1]
		stem = st
		if matched := regexp.MustCompile(`^(Ν|Ρ|ΣΠΙ|ΣΤΡΑΒΟΜΟΥΤΣ|ΚΑΚΟΜΟΥΤΣ|ΕΞΩΝ)$`).MatchString(st); matched {
			stem += "ΟΥΝ"
		}
	}

	// Step 5m
	step5M_pattern := "^(.+?)(ΟΥΜΕ|ΗΣΟΥΜΕ|ΗΘΟΥΜΕ)$"
	matches_5M := getMatchesFromInput(step5M_pattern, stem)
	if len(matches_5M) > 0 {
		st := matches_5M[1]
		stem = st
		if matched := regexp.MustCompile(`^(ΠΑΡΑΣΟΥΣ|Φ|Χ|ΩΡΙΟΠΛ|ΑΖ|ΑΛΛΟΣΟΥΣ|ΑΣΟΥΣ)$`).MatchString(st); matched {
			stem += "ΟΥΜ"
		}
	}

	// Step 6a
	step6A_pattern := "^(.+?)(ΜΑΤΟΙ|ΜΑΤΟΥΣ|ΜΑΤΟ|ΜΑΤΑ|ΜΑΤΩΣ|ΜΑΤΩΝ|ΜΑΤΟΣ|ΜΑΤΕΣ|ΜΑΤΗ|ΜΑΤΗΣ|ΜΑΤΟΥ)$"
	matches_6A := getMatchesFromInput(step6A_pattern, stem)
	if len(matches_6A) > 0 {
		st := matches_6A[1]
		stem = st + "Μ"
		if matched := regexp.MustCompile("ΓΡΑΜ").MatchString(st); matched {
			stem += "Α"
		} else if matched := regexp.MustCompile(`^(ΓΕ|ΣΤΑ)$`).MatchString(st); matched {
			stem += "ΑΤ"
		}
	}

	// Step 6b
	step6B_pattern := "^(.+?)(ΟΥΑ)$"
	matches_6B := getMatchesFromInput(step6B_pattern, stem)
	if len(matches_6B) > 0 {
		st := matches_6B[1]
		stem = st + "ΟΥ"
	}
    if len(stem) == len(word) {
        stem = longStemList(stem)
    }

	// Step 7
	step7_pattern := "^(.+?)(ΕΣΤΕΡ|ΕΣΤΑΤ|ΟΤΕΡ|ΟΤΑΤ|ΥΤΕΡ|ΥΤΑΤ|ΩΤΕΡ|ΩΤΑΤ)$"
	matches_7 := getMatchesFromInput(step7_pattern, stem)
	if len(matches_7) > 0 {
		st := matches_7[1]
		stem = st
		if !regexp.MustCompile(`^(ΕΞ|ΕΣ|ΑΝ|ΚΑΤ|Κ|ΠΡ)$`).MatchString(st) {
			// do nothing
		}
		if matched := regexp.MustCompile(`^(ΚΑ|Μ|ΕΛΕ|ΛΕ|ΔΕ)$`).MatchString(st); matched {
			stem += "ΥΤ"
		}
	}

    return stem
}

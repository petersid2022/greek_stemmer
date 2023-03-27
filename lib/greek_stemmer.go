package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

func loadsettings(x string) []string {
    config_path, err := os.ReadFile("../config/stemmer.yml")
    if err != nil {
        log.Fatal(err)
    }
    config := make(map[interface{}]interface{})
    err2 := yaml.Unmarshal(config_path, &config)
    if err2 != nil {
        log.Fatal(err2)
    }

    out := []string{}
    for _, item := range config[x].([]interface{}) {
        out = append(out, item.(string))
    }
    return out 
}

func isgreek(x string) bool {
    alphabet :=  regexp.MustCompile("^[ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ]+$")
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

func main(){

    // Getting user input
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("word: ")
    input, err := reader.ReadString('\n')
    if err != nil {
        os.Exit(33)
    }
    word := strings.ReplaceAll(input, "\n", "")

    // Setting environment variables
    protected_words := loadsettings("protected_words")
    //step_0_exceptions := loadsettings("step_0_exceptions")
    stem := strings.Clone(word)

    // Step 0
    if utf8.RuneCountInString(stem) < 3 || !isgreek(stem) {
        fmt.Println(stem)
    }     

    if contains(protected_words, stem){
        fmt.Println(stem)
    }

    // Step 1
    



}

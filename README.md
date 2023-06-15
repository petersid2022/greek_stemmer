# Greek stemmer written in Go
WORK IN PROGRESS!
This repository contains an implementation of a Greek stemmer in Go programming language, based on the [paper](https://people.dsv.su.se/~hercules/papers/Ntais_greek_stemmer_thesis_final.pdf) written by Georgios Ntais. The stemmer is designed to extract the stem or root form of Greek words, enabling various natural language processing (NLP) tasks such as text analysis, information retrieval, and machine translation.

## Background
The Greek stemmer is developed as a computational linguistics tool that applies a set of rules and algorithms to transform inflected Greek words into their base or root form. This process, known as stemming, allows for improved analysis and comparison of Greek texts by reducing words to their essential forms.
The algorithm implemented in this project is based on the research paper "Development of a Stemmer for the Greek Language" by Georgios Ntais. The paper serves as a reference guide, providing insights into the Greek language morphology and the rules employed in the stemming process.

## Installation 
Make sure you are running ```go version 1.20.x```

## Usage
To use the Greek stemmer in your Go projects, follow these steps:

1. Install Go: Ensure that Go programming language is installed on your system.
2. Clone the repository: ```bash git clone https://github.com/petersid2022/greek-stemmer-go```
3. Import the package: ```go import "github.com/petersid2022/greek-stemmer-go"```
4. Call the stemmer function: Use the provided function to stem Greek words in your code.

Here's a simple example of how to use the Greek stemmer:

```go
package main

import (
	"fmt"
	"github.com/petersid2022/greek-stemmer-go"
)

func main() {
	word := "καλημέρας"
	stemmedWord := greekstemmer.Stem(word)
	fmt.Printf("Stemmed word: %s\n", stemmedWord)
}
```

## Contributing
Contributions to this Greek stemmer implementation are welcome! If you find any issues, have suggestions for improvements, or want to extend the functionality, please feel free to open an issue or submit a pull request.

## License

(The MIT License)

Copyright (c) 2023 Peter Sideris petersid2022@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

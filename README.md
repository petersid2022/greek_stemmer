# Greek stemmer library written in Go
This repository contains an implementation of a Greek stemmer library written in Go, based on this [paper](https://people.dsv.su.se/~hercules/papers/Ntais_greek_stemmer_thesis_final.pdf) written by Georgios Ntais. The stemmer is designed to extract the stem or root form of Greek words, enabling various natural language processing (NLP) tasks such as text analysis, information retrieval, and machine translation.

## Background
The Greek stemmer is developed as a computational linguistics tool that applies a set of rules and algorithms to transform inflected Greek words into their base or root form. This process, known as stemming, allows for improved analysis and comparison of Greek texts by reducing words to their essential forms.
The algorithm implemented in this project is based on the research paper "Development of a Stemmer for the Greek Language" by Georgios Ntais. The paper serves as a reference guide, providing insights into the Greek language morphology and the rules employed in the stemming process.

## References
* [Georgios Ntais Development of a Stemmer for the Greek Language](https://people.dsv.su.se/~hercules/papers/Ntais_greek_stemmer_thesis_final.pdf)
* [GreekStemmer by skroutz](https://github.com/skroutz/greek_stemmer/)

## Usage
To use the Greek stemmer in your Go projects, follow these steps:

1. Make sure you are running ```go version 1.20.x```
2. Go get the libary ```go get github.com/petersid2022/greek_stemmer@latest```
3. Import it ```import "github.com/petersid2022/greek_stemmer"```

## Example
```go
package main 

import (
    "fmt"
    "github.com/petersid2022/greek_stemmer"
)

func main(){
    x := greek_stemmer.GreekStemmer("ΑΠΑΓΩΓΗ")
    fmt.Println(x)
}
```

## Contributing
Contributions to this Greek stemmer implementation are welcome! If you find any issues, have suggestions for improvements, or want to extend the functionality, please feel free to open an issue or submit a pull request.

## License
This project is licensed under the MIT License. Please see the [LICENSE](./LICENSE) file for more details.

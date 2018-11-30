package main

import (
	"fmt"

	"github.com/DarkAnHell/FastPhish/pkg/analyzer"
	lev "github.com/DarkAnHell/FastPhish/pkg/analyzer/levenshtein"
)

// TODO: This is just for testing, later on we will launch
// every analyzer from here

// TODO: fuzzy logic for deciding upon scores
func main() {
	var a lev.Levenshtein

	stop := make(chan bool)
	out := make(chan analyzer.DomainScore)
	err := make(chan analyzer.DomainError)

	a.Process(
		"twitter.com", []string{"twistter.com", "twitter.com", "google.com", "twi†ter.com", "facebook.es", "random.link.valid"},
		stop,
		out,
		err)

	for v := range out {
		fmt.Printf("For domain %s, levenshtein is %d%% sure it is phishing for domain twitter.com\n", v.Domain, v.Score)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DarkAnHell/FastPhish/pkg/analyzer"
	lev "github.com/DarkAnHell/FastPhish/pkg/analyzer/levenshtein"
)

// TODO: This is just for testing, later on we will launch
// every analyzer from here

// TODO: fuzzy logic for deciding upon scores
func main() {

	// Parse config
	if len(os.Args) < 2 {
		log.Fatalf("missing JSON config file path")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open config file %s: %v", os.Args[1], err)
	}

	var a lev.Levenshtein

	err = a.Load(f)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	stop := make(chan bool)
	out := make(chan analyzer.DomainScore)
	err_chan := make(chan analyzer.DomainError)

	a.Process(
		"twitter.com", []string{"twistter.com", "twitter.com", "google.com", "twi†ter.com", "facebook.es", "random.link.valid"},
		stop,
		out,
		err_chan)

	for v := range out {
		fmt.Printf("For domain %s, levenshtein is %d%% sure it is phishing for domain twitter.com\n", v.Domain, v.Score)
	}
}

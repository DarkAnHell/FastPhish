package main

import (
	"fmt"
	"log"
	"os"

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

	var anal lev.Levenshtein
	err = anal.Load(f)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}


	against := []string{
		"twistter.com",
		"twitter.com",
		"google.com",
		"twi†ter.com",
		"facebook.es",
		"random.link.valid",
	}
	out := anal.Process("twitter.com", against)
	for _, v := range out {
		fmt.Printf("For domain %s, levenshtein is %d%% sure it is phishing for domain twitter.com\n", v.GetName(), v.GetScore())
	}
}

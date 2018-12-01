package levenshtein

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"

	"github.com/DarkAnHell/FastPhish/api"
)

// Levenshtein is just a placeholder to create this "class"
type Levenshtein struct {
	cfg *config
}

// Translates the score given by the algorithm into a usable score.
func (l Levenshtein) dumpScore(score int) int {
	// The same domain
	if score == 0 {
		return 0
	}

	// Closer to the domain (less score in Levenshtein), more likely to be phishing
	return 100 - int(math.Min(float64((score*100)/l.cfg.Threshold), 100.0))
}

// Load reads the configuration and applies changes to the object
func (l *Levenshtein) Load(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not read configuration %v", err)
	}

	var cfg *config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return fmt.Errorf("could not parse configuration: %v", err)
	}
	l.cfg = cfg

	return nil
}

// Process is the implementation of analyzer's Process
func (l Levenshtein) Process(input string, against []string) []api.DomainScore {
	result := make([]api.DomainScore, len(against))
	for index, domain := range against {
		mat := make([][]int, len(input))
		for i := range mat {
			mat[i] = make([]int, len(domain))
		}

		// Set matrix
		for i := 0; i < len(input); i++ {
			mat[i][0] = i
		}
		for i := 0; i < len(domain); i++ {
			mat[0][i] = i
		}

		// Get score
		for i := 1; i < len(input); i++ {
			for j := 1; j < len(domain); j++ {
				var cost int
				if input[i-1] != domain[j-1] {
					cost = l.cfg.Cost
				}

				// Store the minimum between deleting, inserting or subsitute a character
				mat[i][j] = min(mat[i-1][j]+1, mat[i][j-1]+1, mat[i-1][j-1]+cost)
			}
		}

		result[index] = api.DomainScore{
			Name: domain,
			Score: uint32(l.dumpScore(mat[len(input)-1][len(domain)-1])),
		}
	}
	return result
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

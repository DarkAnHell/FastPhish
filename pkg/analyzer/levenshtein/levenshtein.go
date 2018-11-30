package levenshtein

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"

	"github.com/DarkAnHell/FastPhish/pkg/analyzer"
)

// Levenshtein is just a placeholder to create this "class"
type Levenshtein struct {
	// Threshold to use for the activation fuzzy logic
	threshold int

	// Cost to adjust the significance of a letter changing
	cost int
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

// Translates the score given by the algorithm into a usable score
// (see) Analyze's docs)
func (l Levenshtein) translateScore(score int) int {
	// The same domain
	if score == 0 {
		return 0
	}

	// Closer to the domain (less score in Levenshtein), more likely to be phishing
	return 100 - int(math.Min(
		float64((score*100)/l.threshold),
		100.0))
}

// Load reads the configuration and applies changes to the object
func (l *Levenshtein) Load(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return err
	}

	l.threshold = cfg.Threshold
	l.cost = cfg.Cost

	return nil
}

// Process is the implementation of analyzer's Process
func (l Levenshtein) Process(
	input string,
	xpected []string,
	stop chan bool,
	out chan analyzer.DomainScore,
	err chan analyzer.DomainError) {

	// Using a co-routine to return as quickly as possible while we calculate
	go l.internalProcess(input, xpected, stop, out, err)

}

// internalProcess actually calculates the score
func (l Levenshtein) internalProcess(
	input string,
	xpected []string,
	stop chan bool,
	out chan analyzer.DomainScore,
	err chan analyzer.DomainError) {

	// NOTE: Use defer to any cleanup code
	defer close(stop)
	defer close(out)
	defer close(err)

	// Create vars
	var i, j, cost int

	// Pre-calculate what we can
	lenInput := len(input)

	// Iterate all domains and analyze them
	for _, domain := range xpected {
		select {
		default:
			// Pre-calculate for this iteration
			lenDomain := len(domain)

			//  Reset mat
			mat := make([][]int, lenInput)
			for i := range mat {
				mat[i] = make([]int, lenDomain)
			}

			// Expect some outliers or errors
			// TODO: Custom errors

			// Set matrix
			for i = 0; i < lenInput; i++ {
				mat[i][0] = i
			}
			for j = 0; j < lenDomain; j++ {
				mat[0][j] = j
			}

			// Get score
			for i = 1; i < lenInput; i++ {
				for j = 1; j < lenDomain; j++ {
					if input[i-1] == domain[j-1] {
						cost = 0
					} else {
						cost = l.cost
					}

					// Store the minimum between deleting, inserting or subsitute a character
					mat[i][j] = min(
						mat[i-1][j]+1,      // deletion
						mat[i][j-1]+1,      // insertion
						mat[i-1][j-1]+cost) // substitution
				}
			}

			// Return score
			out <- analyzer.DomainScore{Domain: domain, Score: l.translateScore(mat[lenInput-1][lenDomain-1])}

		case <-stop:
			// stop
			return
		}
	}
}

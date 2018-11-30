package levenshtein

import (
	"testing"

	"github.com/DarkAnHell/FastPhish/pkg/analyzer"
)

func auxiliarTestScores(t *testing.T, scores []int, thresholds []int, expected []int) {
	for i := range scores {
		got := translateScore(scores[i], thresholds[i])
		if got != expected[i] {
			t.Errorf("Translations incorrect, got: %d, want: %d. (for %d - %d )", got, expected[i], scores[i], thresholds[i])
		}
	}
}

func auxiliarTestDomains(t *testing.T, domains []string, against []string, expected []int) {
	var l Levenshtein

	stop := make(chan bool, len(domains))
	out := make(chan analyzer.DomainScore, len(domains))
	err := make(chan analyzer.DomainError, len(domains))

	for i := range domains {
		l.internalProcess(domains[i], []string{against[i]}, stop, out, err)
	}

	for y := 0; y < len(domains); y++ {
		select {
		case res := <-out:
			for i, d := range against {

				if d == res.Domain {
					if expected[i] != res.Score {
						t.Errorf("Error in domain checking -> Got: %d for %s against %s. Expected %d", res.Score, d, against[i], expected[i])
					}
					break
				}
			}
		case e := <-err:
			t.Errorf("Unexpected error %s", e)
		}
	}
}

func TestTranslateScoreZero(t *testing.T) {

	scores := []int{0, 0, 0, 100, 5}
	thresholds := []int{5, 1, 100, 10, 5}
	expected := []int{0, 0, 0, 0, 0}

	auxiliarTestScores(t, scores, thresholds, expected)
}

func TestTranslateScoreRandom(t *testing.T) {

	scores := []int{2, 7, 1, 0, 9}
	thresholds := []int{10, 5, 10, 10, 10}
	expected := []int{80, 0, 90, 0, 10}

	auxiliarTestScores(t, scores, thresholds, expected)
}

func TestScoreDomainsExact(t *testing.T) {
	domains := []string{"twitter.com", "google.com", "facebook.es", "random.link.valid"}
	checkAgainst := []string{"twitter.com", "google.com", "facebook.es", "random.link.valid"}
	expected := []int{0, 0, 0, 0}

	auxiliarTestDomains(t, domains, checkAgainst, expected)

}

func TestScoreDomainsPhishing(t *testing.T) {
	domains := []string{"twitter.com", "twi†ter.com", "google.com", "facebook.es", "random.link.valid"}
	checkAgainst := []string{"twittter.com", "twitter.com", "joogle.com", "facebock.es", "rand.link.valid"}
	expected := []int{90, 70, 90, 90, 80}

	auxiliarTestDomains(t, domains, checkAgainst, expected)

}

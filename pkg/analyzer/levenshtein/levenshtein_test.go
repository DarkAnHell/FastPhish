package levenshtein

import "testing"

func TestTranslateScoreZero(t *testing.T) {

	scores := []int{0, 0, 0}
	thresholds := []int{5, 1, 100}
	expected := []int{0, 0, 0}

	for i := range scores {
		got := translateScore(scores[i], thresholds[i])
		if got != expected[i] {
			t.Errorf("Translations incorrect, got: %d, want: %d.", got, expected[i])
		}
	}

}

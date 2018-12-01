package levenshtein

import (
	"testing"
)


func TestTranslateScoreZero(t *testing.T) {
	tt := []struct{
		score 	  int
		threshold int
		expected  int
	}{
		{0, 5, 0},
		{0, 1, 0},
		{0, 100, 0},
		{100, 10, 0},
		{5, 5, 0},
	}

	for _, tc := range tt {
		t.Run("translate score zero", func(t *testing.T) {
			l := &Levenshtein{
				cfg: &config{
					Threshold: tc.threshold,
				},
			}
			got := l.dumpScore(tc.score)
			if got != tc.expected {
				t.Fatalf("failed to calculate score: got %d expected %d", got, tc.expected)
			}
		})
	}
}

func TestTranslateScoreRandom(t *testing.T) {
	tt := []struct{
		score 	  int
		threshold int
		expected  int
	}{
		{2, 10, 80},
		{7, 5, 0},
		{1, 10, 90},
		{0, 10, 0},
		{9, 10, 10},
	}

	for _, tc := range tt {
		t.Run("translate score zero", func(t *testing.T) {
			l := &Levenshtein{
				cfg: &config{
					Threshold: tc.threshold,
				},
			}
			got := l.dumpScore(tc.score)
			if got != tc.expected {
				t.Fatalf("failed to calculate score: got %d expected %d", got, tc.expected)
			}
		})
	}
}

func TestScoreDomainsExact(t *testing.T) {
	tt := []struct{
		domain string
		against []string
		exp []int
	}{
		{"twitter.com", []string{"twitter.com"}, []int{0}},
		{"google.com", []string{"google.com"}, []int{0}},
		{"facebook.es", []string{"facebook.es"}, []int{0}},
		{"random.link.valid", []string{"random.link.valid"}, []int{0}},
	}

	for _, tc := range tt {
		t.Run("score domains exact match", func(t *testing.T) {
			l := &Levenshtein{
				cfg: &config{
					Cost: 1,
					Threshold: 10,
				},
			}
			got := l.Process(tc.domain, tc.against)
			if uint32(tc.exp[0]) != got[0].GetScore() {
				t.Fatalf("expected %d got %d", tc.exp[0], got[0].GetScore())
			}
		})
	}
}

func TestScoreDomainsPhishing(t *testing.T) {
	tt := []struct{
		domain string
		against []string
		exp []int
	}{
		{"twitter.com", []string{"twittter.com"}, []int{90}},
		{"twi†ter.com", []string{"twitter.com"}, []int{70}},
		{"google.com", []string{"joogle.com"}, []int{90}},
		{"facebook.es", []string{"facebock.es"}, []int{90}},
		{"random.link.valid", []string{"rand.link.valid"}, []int{80}},
	}

	for _, tc := range tt {
		t.Run("score domains exact match", func(t *testing.T) {
			l := &Levenshtein{
				cfg: &config{
					Cost: 1,
					Threshold: 10,
				},
			}
			got := l.Process(tc.domain, tc.against)
			if uint32(tc.exp[0]) != got[0].GetScore() {
				t.Fatalf("expected %d got %d", tc.exp[0], got[0].GetScore())
			}
		})
	}
}

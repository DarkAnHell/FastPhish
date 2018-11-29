// Package levenshtein analyzes a given domain against a given list
// of valid domains, using the levenshtein distance between them as an inverted score
//
// For example:
// Input doimain: "twitters.com"
// Given list: "twitter.com","google.com"
// Output: "1, 8"
//
// This means the algorithm thinks that "twitters.com" is ALMOST CERTAIN (distance = 1)
// that this is phishing for "twitter.com", but is very unlikely (distance = 8) that
// this is phishing for "google.com"
package levenshtein

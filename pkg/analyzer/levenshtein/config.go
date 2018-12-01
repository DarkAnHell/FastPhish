package levenshtein

// Config for levenshtein.
type Config struct {
	// cost to adjust the significance of a letter changing
	Cost int `json:"cost"`

	// Threshold to use for the activation fuzzy logic
	Threshold int `json:"threshold"`
}

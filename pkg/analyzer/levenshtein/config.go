package levenshtein

// config for levenshtein.
type config struct {
	// Cost to adjust the significance of a letter changing.
	Cost int `json:"cost"`
	// Threshold to use for the activation fuzzy logic.
	Threshold int `json:"threshold"`
}

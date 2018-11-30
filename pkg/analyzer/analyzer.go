package analyzer

// Analyzer defines the interface which must be used by any analyzer.
type Analyzer interface {

	// TODO: Change string to Domain
	// TODO: Input config

	// Process should take an input (domain) and use it asynchronously
	// to derive a score for it.
	//
	// This call MUST be processed asynchronously to allow a quick call to it,
	// and then retrieve the result when it's ready, since different analyzers will
	// require different times.
	//
	// The score MUST be an integer from 0 to 100, where 100 is absolutely and positively sure
	// of a phishing, and 0 is the opposite.
	// Keep in mind an absolute 100 or 0 is very unlikely, and should only be used when absolutely sure
	//
	// NOTES:
	// - Any error should be returned on err channel
	// - The analyzing process MUST stop if a message comes through the "stop" channel
	// - The result MUST be an int, and must be given to the "out" channel
	// - Example of out channel: ["twitter.com"]:1 , ["google-com"]:8
	Process(
		input string,
		compareAgainst []string,
		stop chan bool,
		out chan DomainScore,
		err chan DomainError)
}

// DomainScore holds a fake "key->value" structure for a domain and it's score
type DomainScore struct {
	Domain string
	Score  int
}

// DomainError holds a fake "key->value" structure for a domain and it's error
type DomainError struct {
	Domain string
	Err    error
}

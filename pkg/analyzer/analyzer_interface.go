package analyzer

// Analyzer defines the interface which must be used by any analyzer.
type Analyzer interface {

	// TODO: Change string to Domain

	// Process should take an input (domain) and use it asynchronously
	// to derive a score for it.
	//
	// This call MUST be processed asynchronously to allow a quick call to it,
	// and then retrieve the result when it's ready, since different analyzers will
	// require different times.
	//
	// NOTES:
	// - The analyzing process MUST stop if a message comes through the "stop" channel
	//
	// - The result MUST be an int, and must be given to the "out" channel
	Process(input string, stop chan bool, out chan int)
}

package analyzer

import "io"

// Analyzer defines the interface which must be used by any analyzer.
type Analyzer interface {

	// Process takes an input (domain) and use it asynchronously
	// to derive a score for it.
	//
	// This call MUST be processed asynchronously to allow a quick call to it,
	// and then retrieve the result when it's ready, since different analyzers will
	// require different times.
	//
	// The score MUST be an integer from 0 to 100, where 100 is absolutely and positively sure
	// of a phishing, and 0 is the opposite.
	// Keep in mind an absolute 100 or 0 is very unlikely, and should only be used when absolutely sure.
	//
	// NOTES:
	// - Any error should be returned on err channel
	// - The analyzing process MUST stop if a message comes through the "stop" channel
	// - The result MUST be an int, and must be given to the "out" channel
	// - Example of out channel: ["twitter.com"]:1 , ["google-com"]:8
	Process(input string, against []string)
	// Load changes the properties for the object according to the information
	// passed.
	//
	// This is used to load configuration.
	Load(r io.Reader) error
}

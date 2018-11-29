// Package analyzer defines the interface that the analyzer must implement
// in order for it to be used.
//
// An analyzer is any tool that is able to take as an input a certain domain
// and output a confidence score of how malicious the domain is.
//
// Any analyzer MUST follow the configuration specifics' of operation,
// which will be passed as input
package analyzer

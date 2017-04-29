package util

// Signal is a sentinel value that indicates an explicit lack of information.
// Signal can be used on pure signal channels, where the presence of a message on the channel is sufficient
// (e.g. "quit channels").
// (inspiration from http://peter.bourgon.org/go-in-production/)
type Signal struct{}

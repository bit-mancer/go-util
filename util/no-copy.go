package util

// NoCopy prevents structs that have an instance of NoCopy from being copied (no-copy check via vet -copylocks checker)
type NoCopy struct{}

// Lock is a trigger for vet -copylocks checker
func (*NoCopy) Lock() {} // nolint

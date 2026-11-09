package common

// SetLowMemory turns the low-memory relay path on for a device short of memory, in a build that was
// not tagged with_low_memory: one Android core serves phones and low-RAM TVs. Buffer sizes stay as
// built. Call it before relaying starts.
func SetLowMemory(enabled bool) {
	LowMemory = enabled
}

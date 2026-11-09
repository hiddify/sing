//go:build !with_low_memory

package common

// LowMemory makes the relay path trade speed for memory. Builds tagged with_low_memory start with it
// on; others can turn it on at run time, before any connection is relayed (see SetLowMemory).
var LowMemory = false

package fmtx

import "bytes"

// BytesToString Take a []byte{} and return a string based on null termination.
// This is useful for situations where the OS has returned a null terminated
// string to use.
// If this function happens to receive a byteArray that contains no nulls, we
// simply convert the array to a string with no bounding.
func BytesToString(byteArray []byte) string {
	n := bytes.IndexByte(byteArray, 0)
	if n < 0 {
		return string(byteArray)
	}
	return string(byteArray[:n])
}

package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// ReadAll 一次读取文件中得所有内容
func ReadAll(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ReadUintFromFile read unit from file
func ReadUintFromFile(path string) (uint64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

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

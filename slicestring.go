package main

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//Compar - Comparison of two slices
func Compar(small []string, big []string) (result bool) {
	var a int

	for i := 0; i < len(small); i++ {
		for _, b := range big {
			if b == small[i] {
				a++
			}
		}
	}

	if a != len(small) {
		return false
	}

	return true

}

//ToLow - Makes all cut letters small
func ToLow(old []string) (new []string) {
	for _, a := range old {
		new = append(new, strings.ToLower(a))
	}
	return new
}

// RandString - There could be a regular hex :)
func RandString(text string) string {
	a := md5.New()
	a.Write([]byte(strings.ToLower(text)))
	return hex.EncodeToString(a.Sum(nil))
}

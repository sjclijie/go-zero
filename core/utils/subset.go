package utils

import "math/rand"

func Subset(set []string, sub int) []string {
	rand.Shuffle(len(set), func(i, j int) {
		set[i], set[j] = set[j], set[i]
	})
	if len(set) <= sub {
		return set
	} else {
		return set[:sub]
	}
}

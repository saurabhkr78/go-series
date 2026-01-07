package main

// TASK 1: Word Frequency Counter
import (
	"fmt"
	"strings"
)

func freqcount(sentence string) map[string]int {
	words := strings.Fields(sentence)
	freq := make(map[string]int)

	for _, word := range words {
		freq[word]++
	}
	return freq
}

func main() {
	input := "go is the fast and modern lang"
	result := freqcount(input)

	for word, count := range result {
		fmt.Printf("%s occurs %d\n", word, count)
	}
}

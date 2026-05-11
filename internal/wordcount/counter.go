package wordcount

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type WordCount struct {
	Word  string
	Count int
}

func CountWords(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	wordCount := make(map[string]int)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		if word != "" {
			wordCount[word]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return wordCount, nil
}

func PrintResults(wordCount map[string]int) {
	words := make([]string, 0, len(wordCount))
	for word := range wordCount {
		words = append(words, word)
	}
	sort.Strings(words)

	for _, word := range words {
		fmt.Printf("%s: %d\n", word, wordCount[word])
	}
}

func PrintResultsSortedByFrequency(wordCount map[string]int) {
	wordCounts := make([]WordCount, 0, len(wordCount))
	for word, count := range wordCount {
		wordCounts = append(wordCounts, WordCount{Word: word, Count: count})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count == wordCounts[j].Count {
			return wordCounts[i].Word < wordCounts[j].Word
		}
		return wordCounts[i].Count > wordCounts[j].Count
	})

	for _, wc := range wordCounts {
		fmt.Printf("%s: %d\n", wc.Word, wc.Count)
	}
}
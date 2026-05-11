package wordcount

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
)

type WordCount struct {
	Word  string
	Count int
}

func CountWordsParallel(filename string, numWorkers int) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	// Канал для строк
	lineChan := make(chan string, 10000)
	
	// Канал для результатов
	resultChan := make(chan map[string]int, numWorkers)
	
	var wg sync.WaitGroup
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(lineChan, resultChan, &wg)
	}
	
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				lineChan <- line
			}
		}
		close(lineChan)
		
		if err := scanner.Err(); err != nil {
		}
	}()
	
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	
	totalCount := make(map[string]int)
	for partialResult := range resultChan {
		for word, count := range partialResult {
			totalCount[word] += count
		}
	}
	
	return totalCount, nil
}

func worker(lineChan <-chan string, resultChan chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	
	localCount := make(map[string]int)
	
	for line := range lineChan {
		localCount[line]++
	}
	
	resultChan <- localCount
}

func CountWords(filename string) (map[string]int, error) {
	return CountWordsParallel(filename, 6)
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
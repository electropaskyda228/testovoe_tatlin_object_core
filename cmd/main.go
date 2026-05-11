package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/electropaskyda228/namecount/internal/wordcount"
)

func main() {
	sortByFreq := flag.Bool("sort-by-freq", false, "Вывести слова в порядке убывания частоты встречаемости")
	numWorkers := flag.Int("workers", 6, "Количество горутин для обработки (по умолчанию 6)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Использование: %s [-sort-by-freq] [-workers N] <filename>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Пример: %s file.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Пример с флагом: %s -sort-by-freq file.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Пример с количеством воркеров: %s -workers 8 file.txt\n", os.Args[0])
		os.Exit(1)
	}

	filename := flag.Arg(0)

	wordCount, err := wordcount.CountWordsParallel(filename, *numWorkers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	if *sortByFreq {
		wordcount.PrintResultsSortedByFrequency(wordCount)
	} else {
		wordcount.PrintResults(wordCount)
	}
}
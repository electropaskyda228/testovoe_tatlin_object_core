package wordcount

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCountWords(t *testing.T) {
	content := `hello
world
hello
golang
test
world
hello
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	result, err := CountWords(tmpFile)
	if err != nil {
		t.Fatalf("Ошибка при подсчёте слов: %v", err)
	}

	expected := map[string]int{
		"hello":  3,
		"world":  2,
		"golang": 1,
		"test":   1,
	}

	for word, expectedCount := range expected {
		if result[word] != expectedCount {
			t.Errorf("Слово '%s': ожидалось %d, получено %d", word, expectedCount, result[word])
		}
	}

	if len(result) != len(expected) {
		t.Errorf("Ожидалось %d уникальных слов, получено %d", len(expected), len(result))
	}
}

func TestCountWordsEmptyFile(t *testing.T) {
	tmpFile := createTempFile(t, "")
	defer os.Remove(tmpFile)

	result, err := CountWords(tmpFile)
	if err != nil {
		t.Fatalf("Ошибка при подсчёте: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Для пустого файла ожидалась пустая map, получено %d записей", len(result))
	}
}

func TestCountWordsWithEmptyLines(t *testing.T) {
	content := `hello

world

hello
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	result, err := CountWords(tmpFile)
	if err != nil {
		t.Fatalf("Ошибка при подсчёте: %v", err)
	}

	expected := map[string]int{
		"hello": 2,
		"world": 1,
	}

	for word, expectedCount := range expected {
		if result[word] != expectedCount {
			t.Errorf("Слово '%s': ожидалось %d, получено %d", word, expectedCount, result[word])
		}
	}
}

func TestCountWordsFileNotFound(t *testing.T) {
	_, err := CountWords("/nonexistent/file.txt")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего файла, но её не было")
	}
}

func TestPrintResults(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	wordCount := map[string]int{
		"banana": 5,
		"apple":  10,
		"cherry": 3,
	}

	PrintResults(wordCount)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "apple: 10\nbanana: 5\ncherry: 3\n"
	if output != expected {
		t.Errorf("Ожидался вывод:\n%s\nПолучено:\n%s", expected, output)
	}
}

func TestPrintResultsSortedByFrequency(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	wordCount := map[string]int{
		"banana": 5,
		"apple":  10,
		"cherry": 3,
		"date":   5,
	}

	PrintResultsSortedByFrequency(wordCount)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "apple: 10\nbanana: 5\ndate: 5\ncherry: 3\n"
	if output != expected {
		t.Errorf("Ожидался вывод:\n%s\nПолучено:\n%s", expected, output)
	}
}

func TestCountWordsLargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем тест с большим файлом в коротком режиме")
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "large.txt")

	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Не удалось создать файл: %v", err)
	}

	lines := 1000000
	for i := 0; i < lines; i++ {
		if i%1000 == 0 {
			_, err := file.WriteString(fmt.Sprintf("test_word_%d\n", i))
			if err != nil {
				file.Close()
				t.Fatalf("Ошибка записи в файл: %v", err)
			}
		} else {
			_, err := file.WriteString("common_word\n")
			if err != nil {
				file.Close()
				t.Fatalf("Ошибка записи в файл: %v", err)
			}
		}
		
		if i%100000 == 0 && i > 0 {
			t.Logf("Создано %d строк...", i)
		}
	}
	file.Close()

	t.Logf("Файл создан, начинаем подсчёт...")
	
	result, err := CountWords(tmpFile)
	if err != nil {
		t.Fatalf("Ошибка при обработке большого файла: %v", err)
	}

	expectedUniqueWords := lines/1000 + 1
	
	if len(result) != expectedUniqueWords {
		t.Errorf("Ожидалось примерно %d уникальных слов, получено %d", expectedUniqueWords, len(result))
	}
	
	if result["common_word"] != lines-lines/1000 {
		t.Errorf("Ожидалось %d вхождений 'common_word', получено %d", lines-lines/1000, result["common_word"])
	}
	
	t.Logf("Тест успешно завершён. Уникальных слов: %d", len(result))
}

func BenchmarkCountWords(b *testing.B) {
	tmpFile := createLargeFile(b, 100000)
	defer os.Remove(tmpFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountWords(tmpFile)
	}
}

func createTempFile(t testing.TB, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать временный файл: %v", err)
	}
	
	return tmpFile
}

func createLargeFile(t testing.TB, lines int) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "large_bench.txt")
	
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Не удалось создать файл: %v", err)
	}
	defer file.Close()
	
	for i := 0; i < lines; i++ {
		if i%1000 == 0 {
			file.WriteString(fmt.Sprintf("unique_word_%d\n", i))
		} else {
			file.WriteString("common\n")
		}
	}
	
	return tmpFile
}
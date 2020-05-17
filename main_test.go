package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_sortWord(t *testing.T) {
	testCases := map[string]string{
		"example": "aeelmpx",
		"macbook": "abckmoo",
		"a":       "a",
	}
	for name, result := range testCases {
		t.Run(name, func(t *testing.T) {
			if result != sortWord(name) {
				t.Fatalf("incorrect sorting for %s", name)
			}
		})
	}

}

func Test_LCDSort(t *testing.T) {
	words := []string{"race", "care", "acre"}
	sorted := []string{"acre", "care", "race"}
	for i, w := range LCDSort(words[:], 4) {
		if sorted[i] != w {
			t.Fatalf("should be %s, got %s", sorted[i], w)
		}
	}
}

func Test_anagrams(t *testing.T) {
	a := newAnagrams()
	words := []string{"bee", "asleep", "elapse", "please", "males", "meals"}
	angrms := [][]string{{"males", "meals"}, {"asleep", "elapse", "please"}}
	for _, w := range words {
		a.put(w)
	}
	for i, ams := range a.getListOf() {
		for j, am := range ams {
			if angrms[i][j] != am {
				t.Fatalf("should be %s, got %s", angrms[i][j], am)
			}
		}
	}
}

func BenchmarkSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, err := os.Open("words.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		a := newAnagrams()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			a.put(strings.ToLower(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		for _, words := range a.getListOf() {
			LCDSort(words, len(words[0]))
		}
	}
}

func BenchmarkConcurrent4Workers(b *testing.B) {
	numOfWorkers := 4
	for i := 0; i < b.N; i++ {
		concurrent(numOfWorkers)
	}
}

func BenchmarkConcurrent8Workers(b *testing.B) {
	numOfWorkers := 8
	for i := 0; i < b.N; i++ {
		concurrent(numOfWorkers)
	}
}

func BenchmarkConcurrent16Workers(b *testing.B) {
	numOfWorkers := 16
	for i := 0; i < b.N; i++ {
		concurrent(numOfWorkers)
	}
}

func concurrent(numOfWorkers int) {
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	a := mergeAnagrams(workerPool(file, numOfWorkers), numOfWorkers)

	for _, words := range a.getListOf() {
		LCDSort(words, len(words[0]))
	}

}

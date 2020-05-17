package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	alphabetLen   = 256
	buffer        = 100
	maxWordLength = 50
)

func main() {
	filename := flag.String("f", "sample.txt", "file name")
	flag.Parse()
	file, err := os.Open(*filename)
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
		fmt.Println(strings.Join(LCDSort(words, len(words[0])), " "))
	}
}

type anagrams struct {
	buckets []map[string][]string
}

func newAnagrams() *anagrams {
	return &anagrams{buckets: make([]map[string][]string, maxWordLength)}
}
func (a *anagrams) put(word string) {
	if len(word) == 0 {
		return
	}
	idx := len(word) - 1
	bucket := a.buckets[idx]
	if bucket == nil {
		bucket = make(map[string][]string)
		a.buckets[idx] = bucket
	}
	sorted := sortWord(word)
	bucket[sorted] = append(bucket[sorted], word)
}

func (a *anagrams) getListOf() [][]string {
	angrms := make([][]string, 0, len(a.buckets))
	for i := 0; i < len(a.buckets); i++ {
		for _, words := range a.buckets[i] {
			if len(words) < 2 {
				continue
			}
			angrms = append(angrms, words)
		}
	}
	return angrms
}

// sortWord sort letters in words for linear time
func sortWord(word string) string {

	count := make([]int, alphabetLen)
	for _, r := range word {
		count[r]++
	}
	sorted := make([]rune, 0, len(word))
	for i, e := range count {
		if e == 0 {
			continue
		}
		for j := 0; j < e; j++ {
			sorted = append(sorted, rune(i))
		}
	}
	return string(sorted)
}

// LCDSort sorts array of words the same length for linear time
// https://www.informit.com/articles/article.aspx?p=2180073&seqNum=2
func LCDSort(words []string, wLen int) []string {
	wsLen := len(words)
	aux := make([]string, wsLen)
	for d := wLen - 1; d >= 0; d-- {
		count := make([]int, alphabetLen+1)
		for i := 0; i < wsLen; i++ {
			count[words[i][d]-'a'+1] += 1
		}
		for i := 0; i < alphabetLen; i++ {
			count[i+1] += count[i]
		}
		for i := 0; i < wsLen; i++ {
			c := count[words[i][d]-'a']
			aux[c] = words[i]
			count[words[i][d]-'a'] += 1
		}
		for i := 0; i < wsLen; i++ {
			words[i] = aux[i]
		}
	}
	return words
}

// rest of file is code for concurrent program execution

func (a *anagrams) merge(other *anagrams) {
	for i := 0; i < maxWordLength; i++ {
		bucket := a.buckets[i]
		if bucket == nil {
			bucket = make(map[string][]string)
			a.buckets[i] = bucket
		}
		otherBucket := other.buckets[i]
		for key, words := range otherBucket {
			bucket[key] = append(bucket[key], words...)
		}
	}
}

func mergeAnagrams(in <-chan *anagrams, numOfWorkers int) *anagrams {
	a := newAnagrams()
	for i := 0; i < numOfWorkers; i++ {
		a.merge(<-in)
	}
	return a
}

func workerPool(file *os.File, numWorkers int) <-chan *anagrams {
	jobs := make(chan string, buffer)
	results := make(chan *anagrams, numWorkers)

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		go worker(&wg, jobs, results)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		jobs <- strings.ToLower(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	go func() {
		wg.Wait()
		close(jobs)
	}()
	return results
}

func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- *anagrams) {
	a := newAnagrams()
	for j := range jobs {
		a.put(j)
		wg.Done()
	}
	results <- a
}

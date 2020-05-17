# Desing

## The problem

Given an input file that contains one word per line, as an output construct a list of all anagrams from that input file. Print those words to the console, where all words that are an anagram should each other should be on the same line.

Consider aspects such as Maintainability, Scalability, Performance, etc

## The steps

### Read file

The only optimization of reading file is reading file line by line to not load the whole into memory.

### Process word 

During word processing we have to pay attention to two word characteristics: length and character set:

> Two words are defined as anagrams if they do share the same letters, but are in a different order

Firstly we group our words by _length_ we use array where array index corresponds word length:

> `car` to second element, `home` to third
	
Each arrays element represents bucket of words. How shall we organize it? Here we pay attention to the unique character set. 
If we sort letters inside each anagram word we get the same character set:
```
act, cat -> sort -> act, act
```
That means we can represent bucket as a map where _key_ is a word ordered char set, _value_ - an array of anagrams.

## Complexity

- Array access - constant time, O(1)
- Sort letters in the word: linear time complexity, O(n) where _n_ - length of the word
- Map key access - constant time, O(1) 
- Sort anagrams, [LSD string sort](https://www.informit.com/articles/article.aspx?p=2180073&seqNum=2): linear time complexity, O(n) where _n_ - length of anagrams array

## Optimization

I compared two version of program: _sequential_ and _concurrent_ on a big(466 551 words file)

My macbook has 4 physical and 8 logical cpu:
```
sctl hw.physicalcpu hw.logicalcpu
```
> hw.physicalcpu: 4
> hw.logicalcpu: 8

```
go test -bench .
goos: darwin
goarch: amd64
pkg: github.com/sergii4/anagrams
BenchmarkSequential-8            	       2	 609577077 ns/op
BenchmarkConcurrent4Workers-8    	       2	 545940748 ns/op
BenchmarkConcurrent8Workers-8    	       2	 570692494 ns/op
BenchmarkConcurrent16Workers-8   	       2	 600963959 ns/op
```

Benchmark show that 4 concurrent program with 4 worker in worker pool is slightly faster than sequential:
> 545940748 ns/op vs 609577077 ns/op
and more workers we have the slower program works

My conclusion that concurrent and parallel (4, 8, 16 workers) program execution doesn't really speed up program and can even slows it down cause of *fork/join* overhead. 

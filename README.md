# Run

Execute from the command line to run the program:
```
go run main.go -f sample.txt
```
where `-f` - the file name to check the anagrams(_sample.txt_ uses by default, so you can just run: `go run main.go`)

To run the program on a text file containing over 466k English words use `words.txt` as a file name argument: `-f words.txt`

If you need an executable, build it with:
```
go build
```
and then run with:
```
./anagrams -f words.txt
```

# Test

To test the program:
```
go test -v
```

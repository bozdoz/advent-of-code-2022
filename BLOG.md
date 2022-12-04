# What Am I Learning Each Day?

### Day 4

**Difficulty: 0/10**

**Time: ~10 min**

Seems like the easiest puzzle so far.  I was waiting for the catch in part 2; didn't happen.

Just swapped a `contains` check with an `overlaps` check.  Surprisingly simple.

I was thinking today about simplifying the day-*.go files to call a predictable function, like **"parseData"**, since that's a pretty common pattern.

Some duplicate code in the parsers, but I'm usually fine to repeat myself twice.

#### Day 4 Update

Abstracting some of the duplicate code to a utility function to `RunSolvers`.  I think this is a pretty good structure, helps for testing in the future, and helps reduce bloated boilerplate:

```go
// A day is just a file reader and the functions to call
// with the input content
type Day[T any] struct {
	FileReader func(string) T
	Fncs       []func(T) int
}
```

TIL about alias declaration vs type definition:

```go
type dataType []string // is a type definition
type dataType = []string // is an alias declaration
```

With the alias, there's no need for all the conversions I was doing.  So I can reduce duplication, and remove all the conversions!  Also it was necessary to work with the `RunSolvers` function.  This also made me have to alter the `getInputFile` function.  I may need to revisit that to adjust `depth` across the board.

### Day 3

**Difficulty: 1/10**

**Time: ~30 min**

I restructured a bit; really going harder on the notion that I should just call `panic` instead of handling errors. This is not a production app, and this could be cleaner and clearer.

Heavy use of runes to get letter scores:

```go
// go characters 'a' is 97, but should be 1, and 'A' is 65, and should be 27
func getLetterScore(l rune) int {
	out := int(l) - 96

	if out < 0 {
		out += 58
	}

	return out
}
```

Curious if there was a better way to do this, but there was no way I was going to enum or map that out.

First part I just split each block in halves and iterated one while checking the other:

```go
first, last := items[:half], items[half:]

for _, letter := range first {
	if strings.ContainsRune(last, letter) {
```

I think it reads well, and I'm not sure how else I could have done it.

For the second part I got to use my `types.Set` utility, and I believe this is the first time using a `map` with a `struct{}` value:

```go
type Set[T comparable] map[T]struct{}

func (set *Set[T]) Add(item T) {
	(*set)[item] = struct{}{}
}
```

I definitely find the syntax rough, so I'm glad to abstract this into its own type.  Note also that I used a generic, which **had** to be `comparable` in order to be used in `map`.

I used a for loop `label`, maybe for the first time, because I did the second part with 3 for loops and a switch statement.  The idea was, for the three strings in the group, add all of the first to a set, check the second against the first set and add those to a set, and check the third against the second set, and continue the outmost loop.  Felt pretty simple.

### Day 2

**Difficulty: 1/10**

**Time: ~20 min**

I think I spent too much time on structure today: creating constants and maps for input data parsing.  Probably I didn't need to create a `struct` for `tournament` either. 

Though I did enjoy writing the `iota`'s:

```go
const (
	ROCK = iota + 1
	PAPER
	SCISSORS
)

const (
	LOSS = iota * 3
	DRAW
	WIN
)
```

I've actually never done a rock, paper, scissors program, so that was mildly interesting.  I fought about making my constants 0-based instead of making `ROCK=1`, so that I could work with modulus to wrap around when `ROCK` beats `SCISSORS`.

Instead I just had to include an if statement:

```go
var yourScore, opponent int
// ...
case win:
  choice := opponent + 1
  if choice == 4 {
    choice = 1
  }
  yourScore += WIN + choice
```

Didn't spend any time on utilities, or the test, or in the `day-02.go` file; so I think that's a huge win.

### Day 1

**Difficulty: 1/10**

**Time: ~15 min**

Decided not to worry about panicking, going forward.

Did a lot of prep before Day 1 to make sure I had a structure that would require the least amount of extra work each day.

Here's how it's broken down:

```bash
./create-day.sh 01
```

Then, copy over the example input into `example.txt`, then copy over the answer into `day-01_test.go`:

```go
// fill in the answers for each part (as they come)
var answers = map[int]int{
	1: 24000,
	2: 45000,
}
```

Then, update the `dataType` and `fileReader` for each day:

```go
// today's input data type
type dataType []string

// how to read today's inputs
var fileReader = utils.ReadNewLineGroups
```

And everything else is ready to **go**.

Now I can run `go test ./01`, until eventually I run `go run ./01`, to get the answers and timings.

I've also decided I should keep track of my time and difficulty of each day.

Today, I was a bit frustrated just at splitting empty new lines, trimming the last line, then splitting new lines, converting to ints, and summing.  Just a bunch of work to parse the input.

#### Day 1 Update

Added a min heap to day 1 and some benchmark tests to see how bad my sorting implementation was:

```bash
> go test ./01 -bench=.
goos: linux
goarch: amd64
pkg: github.com/bozdoz/advent-of-code-2022/01
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
BenchmarkSort-2            13512             86409 ns/op
BenchmarkHeap-2            12837             94467 ns/op
PASS
```

Doesn't seem that bad: sorting takes 8,000 nanoseconds longer? that's 0.008 milliseconds!
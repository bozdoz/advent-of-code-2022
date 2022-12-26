# What Am I Learning Each Day?

### Day 17

**Difficulty: 8/10** ★★★★★★★★☆☆

**Time: ~7 hrs**

I like to plan out my strategy in comments before writing the code, then adding code one step at a time

```go
// shape begins falling 2 from left, and 3 from bottom
// 1. add shape to space
// 2. air pushes left or right
// 3. shape falls 1 unit
// 4. check for collision
// 5. repeat at step 2, or add next shape at step 1
```

Part 1 originally in 1.187ms.  After refactor, Part 1 took 680ms, and Part 2 took 3.3 seconds.

I really wanted to try a bitmask today, because I have never really used one. I think it was a fine idea, though hard to keep track of.  Added some printing functions to help, and a `binaryStringToInt` function.  I also had to push the shapes to the left-hand-side, according to the width of each, and the width of the space (7):

```go
shape{
	[]int{
		binaryStringToInt("1111") << 3,
	},
	4,
}
```

The collision function is where the bitmask helped:

```go
// test if space row & shape row == 0
if (*space)[j]&(shape.outline[i]>>x) != 0 {
	return true
}
```

And adding the shape to the space like so:

```go
// add shape to space
for i := 0; i < len(shape.outline); i++ {
	j := y + i
	// shifting by x
	alignedPart := shape.outline[i] >> x
	if j > len(*space)-1 {
		*space = append(*space, alignedPart)
	} else {
		// adding shape
		(*space)[j] ^= alignedPart
	}
}
```

I ran into a lot of issues with Part 2, related to how I was padding the space before adding the shape, and how I was checking for a drop pattern (checked the space originally, since I could visibly see it in printouts).  Checking the x sequence was easier, and creating generator functions for shape and air pushing made them more isolated from depending on `i` iterations:

```go
func generator[T any](arr []T) func() T {
	i := -1
	l := len(arr)
	return func() T {
		i++
		return arr[i%l]
	}
}
```

And I learned a bit about rolling hashes to detect patterns:

```go

// simple hash function from:
// https://golangprojectstructure.com/hash-functions-go-code/
func djb2(data []int) uint64 {
	hash := uint64(5381)

	for _, b := range data {
		hash += uint64(b) + hash + hash<<5
	}

	return hash
}

// try rolling hash?
// we believe `source` ENDS with a pattern
func findRepeatingPatternFromEnd(source []int) (length int) {
	end := len(source) - 1
	lowerLimit := 5
	upperLimit := len(source)/2 - 1

	for w := lowerLimit; w < upperLimit; w++ {
		a := source[end-w:]
		b := source[end-(w*2)-1 : end-w]

		if djb2(a) == djb2(b) {
			// hashes match; we have a pattern
			return w + 1
		}
	}

	return
}
```

Overall, I learned to hate today, even though I started really enjoying it.  I feel like I was so close yet so far.

### Day 16

**Difficulty: 7/10** ★★★★★★★☆☆☆

**Time: ~3 hrs**

So this reminded me of last year's [Day 23](https://github.com/bozdoz/advent-of-code-2021/blob/354349e4943eba626edd877507c85f5df25d235b/23/shuffle.go).  What I did was create a priority queue, with caching, and storing each state, while discovering each possible next state.

It didn't seem very performant, but I'm happy enough with the results:

```sh
1 | 1915 (38.370372572s)
2 | 2772 (10.330658624s)
```

I really tried to avoid using pointers, and I think it didn't matter much, except for updating pressure:

```go
func (cur *state) addPressure(valves map[string]valve) {
	for open := range cur.valvesOpen {
		cur.pressure += valves[open].flow
	}
}
```

For Part 2, I ran the same function once, returning unique paths according to which valves were open.  I cached the paths according to those valves, and repeatedly updated the best pressure value:

```go
// logic here is that if there are 2 of us, we should each open ~50%
percentOpen := float64(len(state.valvesOpen)) / viableValveCount

// test passes with 40% and 60%, but
// actual puzzle does not pass after 8 seconds:
// 		bestDuoPath Count: 146
// trying 30% and 70% with actual puzzle worked:
// 		bestDuoPath Count: 2422
if percentOpen > 0.3 && percentOpen < 0.7 {
	// we can cache here again on valvesOpen, and update max(pressure)
	key := state.hashValvesOpen()
	pathPressure, ok := bestDuoPaths[key]

	if !ok || state.pressure > pathPressure {
		bestDuoPaths[key] = state.pressure
	}
}
```

And I contemplated implementing a bitmask for `valvesOpen`, which is a unique Set:

```go
func (cur state) hashValvesOpen() string {
	open := make([]string, 0, len(cur.valvesOpen))

	// tempted to do a bitmask
	for key := range cur.valvesOpen {
		open = append(open, key)
	}

	sort.Slice(open, func(i, j int) bool {
		return open[i] < open[j]
	})

	asString := fmt.Sprint(open)

	// lazily omit the "[]" from the Sprint
	return asString[1 : len(asString)-1]
}
```

A bitmask would help so I could have a set, but not need to order it to create the hash (since maps wouldn't be ordered identically).

Also, today I added the `-part` flag to the `utils.runSolvers` function, so that I wouldn't have to wait for part 1 to finish before I could start part 2.

### Day 15

**Difficulty: 7/10** ★★★★★★★☆☆☆

**Time: ~3 hrs**

Finally successful at:

```sh
1 | 5809294 (392.930042ms)
2 | 10693731308112 (121.933359ms)
```

Tried removing the sensor sorting, and got *worse*—almost 2x worse—performance (388.931311ms).

```go
// start with most touched sensor
sort.Slice(sensors, func(i, j int) bool {
	return sensors[i].touches > sensors[j].touches
})
```

Obviously, this brute-force was going to take too long; but I did try to at least chunk the search into quadrants:

```go
scanArea := func(xmin, xmax, ymin, ymax int) {
	for x := xmin; x <= xmax; x++ {
		for y := ymin; y <= ymax; y++ {
			// check if x,y is definitely not a beacon
			// but also definitely not a scanned beacon
			if space.beacons.Has(image.Point{x, y}) {
				continue
			}
			if space.couldBeBeacon(x, y) {
				done <- image.Point{x, y}
			}
		}
	}
}

// run the scanner at four spots, in four quadrants
go scanArea(space.xmin, space.xmax/2, space.ymin, space.ymax/2)
go scanArea(space.xmin, space.xmax/2, space.ymax/2, space.ymax)
go scanArea(space.xmax/2, space.xmax, space.ymin, space.ymax/2)
go scanArea(space.xmax/2, space.xmax, space.ymax/2, space.ymax)

i := 0
for {
	select {
	// sanity check stdout
	case <-time.After(1 * time.Second):
		i++
		fmt.Printf("%d... ", i)
	case point := <-done:
		fmt.Println("found", point)
		return point.X*4000000 + point.Y
	}
}
```

This year, I felt like I finally *understood* Manhattan Distance; like I wouldn't need to look up the formula again.  Last year, I remember adding it to a solution, because I was frantically trying suggestions from reddit.

Today I felt quite confident adding channels and goroutines.  Added similar `for...select` loops as yesterday, to track infinite loops, but I also added a goroutine just to avoid dealing with return values; I like the concept of just calling a function, and saying: "just tell this channel when you're done"...:

```go
// maybe lazy: passing channel to avoid checking return values
found := make(chan image.Point)

// need a goroutine for channel to receive the value
go func() {
	for _, sensor := range sensors {
		space.scanClockwise(sensor, found)
	}
}()

return <-found
```

Otherwise, I was worried of writing `scanClockwise` like this:

```go
func (space *space) scanClockwise(sensor *sensor) (point image.Point, found bool) {
	// didn't find it, most of the time
	return image.Point{}, false
}
```

And then checking that `bool` value on each iteration.  Seemed cleaner (and maybe idiomatic) to just wait for a value to be sent to the channel.

### Day 14

**Difficulty: 2/10** ★★☆☆☆☆☆☆☆☆

**Time: ~60 min**

First run:

```sh
1 | 774 (5.421103ms)
2 | 22499 (267.526678ms)
```

First data structure:

```go
// TODO: maybe these can just be generic "obstacles"
type rocks = types.Set[image.Point]
type sand = types.Set[image.Point]

type area struct {
	rocks  rocks
	sand   sand
	bottom int
	isFull bool
}
```

After minor refactor (**~2x faster**): 

```sh
1 | 774 (3.31752ms)
2 | 22499 (134.481686ms)
```

With simplified data structure:

```go
type obstacle = types.Set[image.Point]

type area struct {
	obstacle obstacle
	bottom   int
}
```

This helps by preventing maintaining/checking of two separate Sets.  For the answer, instead of returning `len(area.sand)`, we can just count iterations of `area.dropSand()`:

```go
for area.dropSand() != nil {
	ans++
}
```

I also enjoyed using go routines and `time.After` to prevent infinite loops while developing/testing, doing a timeout function for the **first time**, maybe using `select` for the **first time**:

```go
// prevent infinite loops
done := make(chan bool)
go func() {
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("--- timeout! ---")
		os.Exit(1)
	case <-done:
	}
}()

for area.dropSand() != nil {
	ans++
}

// stop timeout
done <- true
```

Overall quite a fun puzzle I think.  Quite straightforward and understandable.  No tricks.  No surprises.  Not much of a challenge, but interesting to implement.  Love the idea of collision detection, due to interest in video games.  The End.

### Day 13

**Difficulty: 6/10** ★★★★★★☆☆☆☆

**Time: ~4 hrs**

Did a lot of `go test -run TestSorted -timeout 800us ./13`.  I was having a hard time parsing (partially because I borrowed code from last year's [Day 18](https://github.com/bozdoz/advent-of-code-2021/blob/7851692ea585dac1a2c139dc65755926d58d0bbb/18/pair.go)).  The parsing was surprisingly difficult, as I kept having issues with slices of `any`, pointers to slices of `any`, and copying of slices of `any` to pointers.

**First time** using `recover`, because I had to debug why my type assertions were failing:

```go
// catch panics just in case the type assertions are incorrect
defer func() {
	if rec := recover(); rec != nil {
		log.Printf("panic: (a) %[1]T %[1]v (b) %[2]T %[2]v\n", a, b)
	}
}()
```

**First time** I think using the type assertions (maybe called something else), instead of type switches (also maybe called something else):

```go
// figure out types
aAsInt, aok := a.(int)
bAsInt, bok := b.(int)

// quality variable names here
if aok && bok {
	// BOTH INTS
```

Copied over types.Stack from last year's solutons too.

### Day 12

**Difficulty: 2/10** ★★☆☆☆☆☆☆☆☆

**Time: ~45 min**

Today I copied from [2021 Day 15](https://github.com/bozdoz/advent-of-code-2021/blob/main/15/cave.go), where I first used **Dijkstra's algorithm**, with some improvements:

1. removed the `start` field in the grid struct (start is any cell with distance==0)
2. way better at iterating and creating a grid now
3. used maxint instead of `height * width * 10`, or `math.Inf()`; got maxint from stackoverflow: `const max = int(^uint(0) >> 1)`

One thing I kept identical was the `updateNeighbours` function.

The priority queue I kept mostly similar, except the logic of the puzzle has changed slightly (which is why I found today's puzzle vastly simple).

```go
if neighbour.visited || neighbour.height-square.height > 1 {
	// already visited, or
	// we can only walk up a square at most 1 higher
	continue
}

// update dijkstra's distance
neighbour.distance = utils.Min(
	// we've walked one extra step (+1)
	square.distance+1,
	neighbour.distance,
)
```

The first change is that we can only walk up a square that's at most 1 higher than our current; and the second change is that the distance just increases by 1 each square. Otherwise, it's about the same as 2021 Day 15.

For Part 2, I again threw a `part` flag on the parser function, and added a small if statement:

```go
if part == 2 && char == 'a' {
	// all a's are starting spots
	square.distance = 0
}
```

I was also giving up on prefixing types with "T":

```go
square := &square{
	distance: max,
}
```

Instead, I'm doing what I have feared, and naming variables identical to their types, because it's frustrating to come up with unique names for each.

### Day 11

**Difficulty: 7/10** ★★★★★★★☆☆☆

**Time: ~150 min**

Part one was easy, but Part Two I did not understand.  I roughly know about modulus, and I was pointed in the direction of **Chinese Remainder theorem**, but I did not understand the application of it.  I also have to look up **pairwise coprime**.

I found the input parsing very lazy today.

Lots of this: 

```go
op := line[len("  Operation: new = "):]
```

I also ran into an issue where `val^2` apparently doesn't equal `val*val`, so I had to use math.Pow, which is always super annoying:

```go
case EXPONENT:
	val = int(math.Pow(float64(val), float64(monkey.operationNum)))
```

I think this is the **first time** I used a `for range` loop:

```go
for range [10000]struct{}{} {
	monkeys.inspect(2)
}
```

Maybe dumb.

### Day 10

**Difficulty: 5/10** ★★★★★☆☆☆☆☆

**Time: ~60 min**

Some issues with understanding the puzzle, but seemed vastly easy compared to understanding the rope movement of Day 9.  

Somehow I had issues with tests.  I had a program test that succeeded individually, but not when I run the ExampleOne test. So I must have some issue with parallelism that I might want to sort out eventually.

### Day 9

**Difficulty: 7/10** ★★★★★★★☆☆☆

**Time: ~180 min**

I'm at a complete loss with part 2.  This is a day where the examples and explanation doesn't seem to cover it. In particular, I'm looking at the deliberately unexplained diagram for "Up 8":

```
.........H..
.........1..
.........2..
.........3..
........54..
.......6....
......7.....
.....8......
....9.......
```

And wondering how on earth "5" is adjacent to "4", since that would never happen on a two-node rope.  So the explanation just doesn't seem to suffice for part 2, and I will have to take time to think on this, or seek out an answer elsewhere.

The only warning given was "*be careful: more types of motion are possible than before*".

Blech. I played around with some snake games and videos posted to reddit, and found this reddit solution which gaave me my answer:

https://github.com/EbbDrop/AoCropeSnake/blob/main/src/main.rs#L19-L24

And this:

https://docs.rs/num/latest/num/fn.signum.html

Which I already had a function for.

I guess I can also say today is the **first time** I've used `image.Point`, as it is a 2d vector type with basic methods for adding and subtracting.  Easier than writing my own I think, barely.

Dealing with the linked list today was a little difficult, as I had values that repeatedly needed updating, so I had to cast back to its original value type:

```go
prevPoint := prev.Value.(image.Point)
curPoint := cur.Value.(image.Point)
diff := prevPoint.Sub(curPoint)
```

I also updated a string character by index with:

```go
out[y] = out[y][0:x] + label + out[y][x+1:]
```

### Day 8

**Difficulty: 4/10** ★★★★☆☆☆☆☆☆

**Time: ~60 min**

Today got me mostly by appearing challenging.  I thought it was going to do one of those "Now extrapolate the data 10 times in every direction!"  But today was also a lot of duplication.

Today was the first grid.  I maintained a row-column approach, and may later move this grid to a common type.

I also named custom types prefixed by "T", because I've been having an annoying time naming variables like:

```go
type grid [][]int

// bleeeeeecch
var grid grid
```

I used a `rune` to set the base of the stringified numbers in the input data:

```go
// used for converting rune to int
const zero rune = '0'

// ...
grid[r][c] = int(char - zero)
```

I assume this is more performant than `strconv.Atoi`, but what do I know?

I also updated `types.Set` to take an arbitrary number of items:

```go
func (set *Set[T]) Add(items ...T) {
	for _, item := range items {
		(*set)[item] = struct{}{}
	}
}
```

And this didn't affect the other implementations, which is great.

I also was highly tempted to run a function inside a for loop, because of all the duplication:

```go
for i := r + 1; ; i++ {
	if should_break(i, c, i == g.height-1) {
		break
	}
}
```

something like:

```go
for i := r + 1; !should_break(i, c, i == g.height-1); i++ {
}
```

but I worried that I'd get calls about this.

I mean, already it's a ridiculously named function, that's only meant to de-duplicate the tree height checking, score accumulator.  Anyway...


### Day 7

**Difficulty: 3/10** ★★★☆☆☆☆☆☆☆

**Time: ~60 min**

Today was mostly challenging to create a directory structure.  I kind of anticipated far too many issues that probbaly never came up (multple `ls` calls, multiple `cd` calls, inaccurate `cd` calls, etc).

I think I had a bit too much duplication, but not too bad.  The entire parser felt quite straight-forward, and I didn't really have any errors going through it.

I created a custom type, which I think helped quite a bit:

```go
type dir []*entry

func (*dir).current() *entry
func (*dir).move(dir string) *entry
```

As I mentioned above, I could have returned an `error` in `move` if the directory was invalid.

Again, I'm puzzled why the Part 2's haven't been more challenging than the Part 1's, but maybe it's because of how much I'm anticipating issues.

### Day 6

**Difficulty: 2/10** ★★☆☆☆☆☆☆☆☆

**Time: ~40 min**

For a brief moment, I thought about getting the first four letters and comparing them:

```go
a, b, c, d := data[i], data[i+1], data[i+2], data[i+3]
```

Until, obviously, I knew this would be a pain, even for 4 letters.

This is the **first time** I've used a doubly-linked list, with `container/list`, and also the **first time** declaring multiple variables in a for loop:

```go
var l *list.List

cur := data[i]

// check if current is in list already, moving backwards
for e, j := l.Back(), l.Len(); e != nil; e, j = e.Prev(), j-1 {
	if e.Value == cur {
		// match means we need to clean the list
		removeFront(l, j)

		// add element
		l.PushBack(cur)
		continue outer
	}
}
```

The syntax is a little awkward, but somewhat similar to other languages:

```go
e, j := l.Back(), l.Len()
```

I liked using the list package, since it helped me avoid writing `Push` and `Pop` methods, or dealing with slice/array memory allocation issues (even though I knew the capacity of the list could only be 4-1 or 14-1).

I did another for loop label, which makes me wonder if I'm getting lazy...

Apparently today is the first day I've used `t.Run` in tests:

```go
runs := map[string]int{
	"bvwbjplbgvbhsrlpgdmjqwftvncz":      5,
	"nppdvjthqldpwncqszvftbrmjlhg":      6,
	"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 10,
	"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  11,
}

for k, v := range runs {
	t.Run(fmt.Sprintf("%q should be %d", k, v), func(t *testing.T) {
```

So, doing this, we can have multiple test cases inside of a test suite (parent function).

The timing of today's puzzle seemed good too:

```bash
> go run ./06
1 | 1282 (99.665µs)
2 | 3513 (507.232µs)
```

From matching 4 sequential chars to 14 was only ~5 times slower.

### Day 5

**Difficulty: 4/10** ★★★★☆☆☆☆☆☆

**Time: ~60 min**

Parsing was annoying today.  I also went back and forth about how to add the data to the data structure, and how to process it.  Seemed easiest to parse the data in reverse order:

```go
lines := strings.Split(drawing, "\n")

// -2 because we omit the stack numbers
for i := len(lines) - 2; i >= 0; i-- {
	// we can ignore trailing space
	v := strings.TrimRightFunc(lines[i], unicode.IsSpace)

	for j := 0; j < len(v); j += 4 {
		// get letter inside of "[]"
		letter := v[j+1]

		if letter != space {
			stack := j / 4
			// map initialization continues to be annoying
			_, ok := (*stacks)[stack]

			if !ok {
				(*stacks)[stack] = &crates{}
			}

			(*stacks)[stack].push(letter)
		}
	}
}
```

The map initialization continues to bother me.  Everytime I have a nested data structure, with pointers, I have this pain where I have to check that maps or arrays are initialized, or `nil`.  It's likely that I could have gotten away with avoiding pointers.  Not really interested in going back to refactor to find out.

I implemented a remove and add function for the crates:

```go
type crates []byte

func (c *crates) pop() byte {
	old := *c
	n := len(old)
	out := old[n-1]
	*c = old[0 : n-1]

	return out
}

func (c *crates) push(b byte) {
	*c = append(*c, b)
}
```

And a function to move slices of crates:

```go
// move
l := len(*getFrom)
*addTo = append(*addTo, (*getFrom)[l-num:l]...)
// remove
*getFrom = (*getFrom)[:l-num]
```

Also, this is the first day that doesn't output `int` type; which forced me to redo my generic `runSolvers` already.

### Day 4

**Difficulty: 0/10** ☆☆☆☆☆☆☆☆☆☆

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

**Difficulty: 1/10** ★☆☆☆☆☆☆☆☆☆

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

For the second part I got to use my `types.Set` utility, and I believe this is the **first time** using a `map` with a `struct{}` value:

```go
type Set[T comparable] map[T]struct{}

func (set *Set[T]) Add(item T) {
	(*set)[item] = struct{}{}
}
```

I definitely find the syntax rough, so I'm glad to abstract this into its own type.  Note also that I used a generic, which **had** to be `comparable` in order to be used in `map`.

I used a for loop `label`, maybe for the **first time**, because I did the second part with 3 for loops and a switch statement.  The idea was, for the three strings in the group, add all of the first to a set, check the second against the first set and add those to a set, and check the third against the second set, and continue the outmost loop.  Felt pretty simple.

### Day 2

**Difficulty: 1/10** ★☆☆☆☆☆☆☆☆☆

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

**Difficulty: 1/10** ★☆☆☆☆☆☆☆☆☆

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
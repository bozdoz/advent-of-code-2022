# What Am I Learning Each Day?

- [Day 25](#day-25) ★★★
- [Day 24](#day-24) ★★★★★★★
- [Day 23](#day-23) ★★★★★★
- [Day 22](#day-22) ★★★★★★★★★★
- [Day 21](#day-21) ★★
- [Day 20](#day-20) ★★★★★★
- [Day 19](#day-19) ★★★★★★★★★★
- [Day 18](#day-18) ★★★★★
- [Day 17](#day-17) ★★★★★★★★
- [Day 16](#day-16) ★★★★★★★
- [Day 15](#day-15) ★★★★★★★
- [Day 14](#day-14) ★★
- [Day 13](#day-13) ★★★★★★
- [Day 12](#day-12) ★★
- [Day 11](#day-11) ★★★★★★★
- [Day 10](#day-10) ★★★★★
- [Day 9](#day-9)  ★★★★★★★
- [Day 8](#day-8)  ★★★★
- [Day 7](#day-7)  ★★★
- [Day 6](#day-6)  ★★
- [Day 5](#day-5)  ★★★★
- [Day 4](#day-4)  ☆
- [Day 3](#day-3)  ★
- [Day 2](#day-2)  ★
- [Day 1](#day-1)  ★

### Day 25

**Difficulty: 3/10** ★★★☆☆☆☆☆☆☆

**Time: ~1 hr*

For Part 1 anyway (I assume there is no Part 2), today was pretty easy.  Just had to look up how to convert base 10 to base 5, which predictably had to do with quotients and remainders, and the conversion to the snafu-string-system was fulfilled with carrying over a single `1` if the value was `>2`.  Sometimes this caused the value to be `5`, which also had to wrap around back to `0`:

```go
// got this from https://meteorconverter.com/conversions/number-bases/10-to-5?input=600
func decimalToSnafu(dec int) (snafu string) {
	remainders := []int{}

	for dec > 0 {
		quotient := dec / 5
		remaining := dec % 5
		remainders = append(remainders, remaining)
		dec = quotient
	}

	out := []string{}
	carried := 0

	for _, val := range remainders {
		// 3 == 5-2 = "1="
		// 4 == 5-1 = "1-"
		withCarried := val + carried

		if withCarried > 2 {
			carried = 1
		} else {
			carried = 0
		}

		switch withCarried {
		case 3:
			out = append(out, "=")
		case 4:
			out = append(out, "-")
		default:
			// sometimes this could be 5, and should be 0
			out = append(out, fmt.Sprint(withCarried%5))
		}
	}

	if carried != 0 {
		out = append(out, "1")
	}

	// reverse
	for i := len(out) - 1; i >= 0; i-- {
		snafu += out[i]
	}

	return
}
```

This was a little awkward, since I'm appending to a `remainders` slice, then an `out` slice, then reversing that `out` slice to append to the `snafu string`.  And that's it!

**Part 1 ran in 437.804µs**

### Day 24

**Difficulty: 7/10** ★★★★★★★☆☆☆

**Time: ~2 hrs*

**Part 1 runs in 7.5seconds**.  Not sure what to give for priority queue priority; currently trying: `distance+minutes`, where distance is manhattan distance

Part 2 took about 8 min to refactor.  Somehow I made Part 1 run much quicker (maybe by removing `state.steps` and a single if statement?). Anyway, new times:

```sh
1 | 251 (308.40433ms)
2 | 758 (5.262698437s)
```

This is totally agreeable.

I tried to simplify the `state` type even more, but it didn't seem to affect performance.  Even conditionally calculating the manhattan distance didn't help.  

Part 2 didn't require much at all in terms of changes to my blizzard logic: only needed to add an argument for `minutes int`; so, instead of starting the pathfinding at 0 minutes always, we can set a given time.  And then Part 2 is easy:

```go
vly := parseInput(data)

// just swap the start and end to go back and forth
swap := func() {
	start := vly.start
	end := vly.end
	vly.start = end
	vly.end = start
}

first := vly.pathFinder(0)

swap()

// start at the minute you finished the last trip
andBack := vly.pathFinder(first)

swap()

finally := vly.pathFinder(andBack)

return finally
```

My data structures were separated between a static, unchanging struct, and a dynamic one:

```go
type valley struct {
	height, width int
	start, end    image.Point
	walls         types.Set[image.Point] // walls don't move
	blizzards     map[rune]map[image.Point]struct{}
}
```

I could probably refactor `blizzards` above to be a `types.Set[image.Point]` also, since I decided it too should not change. The only struct that changes is `state`:

```go
type state struct {
	position image.Point
	minutes  int
}
```

And with that I can determine priority (via manhattan distance and time), and whether there's a blizzard overlap (by taking the position, wrapping around the area in each of the 4 directions, and determine if a blizzard is at a starting position `X` number of minutes away):

```go
width := vly.width - 2   // ignores walls
height := vly.height - 2 // ignores walls

// check if '<' is `minutes` to the right of `position`
check := pos
// Y is actually column; also +/- 1 to handle walls
check.Y = (check.Y+min-1)%width + 1
_, ok = vly.blizzards[left][check]

if ok {
	updateCache(hash, true)
	return true
}
```

I'm frustrated by how complex (though barely) the algorithm is: `check.Y = (check.Y+min-1)%width + 1`.  I'd love to clean that up with better variable names (`Y` in place of `Column` is quite upsetting).  The latter part would look cleaner as a `types.Set` by just doing: `if vly.blizzards[left].Has(check) {}`.

I cached this because I thought it would be quite expensive, though haven't tested without the cache.

A cheap way to avoid allowing the path finder to exit the map was to add a wall above the start:

```go
// add wall above start to avoid moving out
valley.walls.Add(image.Point{-1, 1})
```

This allowed me to not bother with introducing this as extra logic in the path finder.

**Actually**, as of writing this, I realized I should add the same kind of wall at the exit, since Part 2 starts at the exit:

```go
// add wall below end to avoid moving out
valley.walls.Add(image.Point{height, width - 2})
```

**And that saved me 4 seconds! 🤦‍♀️**

```sh
1 | 251 (319.13061ms)
2 | 758 (848.897889ms)
```

Finally, the algorithm I used is a priority queue, probably A* (even though I can't really distinguish between A* and Dijkstra, of BFS for that matter). I believe it's A* because it's "intelligent"™️:

```go
// prioritize states by distance and minutes?
// this PQ is in ASC order 🤷‍♀️
pq.PushValue(&next, distance+next.minutes)
```

But yes, this is an algorithm that seems to be used quite heavily in AOC: think I've used it 8 times over the last two years.

I was very happy with today's puzzle. It was straight-forward.  I'm glad I treated the blizzards as static, and avoided trying to update all of them every minute.  Laziness makes us better developers.

Apparently my laziness has gotten the better of me, as the test for Part 1 no longer passes.  I suspect that my measure of priority is not accurate.  I've just skipped this test and flagged it as TODO.

### Day 23

**Difficulty: 6/10** ★★★★★★☆☆☆☆

**Time: ~3 hrs**

I'm really unsure how long I took on this one.  I went back and forth with it for awhile.  Today seemed very easy to parse, and a rather simple data structure too:

```go
type grid map[int]map[int]struct{}

func parseInput(data inType) grid {
	grid := grid{}

	for r, line := range data {
		for c, char := range line {
			if char == '#' {
				grid.set(r, c)
			}
		}
	}

	return grid
}
```

I then added added some CRUD operations for the `grid` type:

```go
func (g *grid) set(r, c int) {
	if (*g)[r] == nil {
		(*g)[r] = map[int]struct{}{}
	}
	(*g)[r][c] = struct{}{}
}

func (g *grid) get(r, c int) bool {
	if (*g)[r] == nil {
		(*g)[r] = map[int]struct{}{}
	}
	_, ok := (*g)[r][c]

	return ok
}

func (g *grid) delete(r, c int) {
	delete((*g)[r], c)
}
```

I leaned on `image.Point` again for dealing with where the elves plan to move to:

```go
// map moveTo -> moveFrom
planned := map[image.Point]image.Point{}
contested := map[image.Point]struct{}{}
```

I thought iterating over these maps might take too much memory, but I was surprised by the results:

```sh
1 | 4181 (16.571719ms)
2 | 973 (1.493149216s)
```

Part 2 actually took me maybe less than a minute to complete; part 1 and part 2 together:

```go
// stop at 10 if part 1
if part == 1 && rounds == checkAt {
	// get square from bounds
	bounds := g.bounds()

	height := bounds[2] - bounds[0] + 1
	width := bounds[1] - bounds[3] + 1
	size := height * width
	empty := size - elves

	return empty
}

// all the elves stopped moving
if elves == still {
	return rounds
}
```

Super simple.  The only problem I seemed to have was dealing with the alternating order of directional checks around the elves:

```go
for r, row := range *g {
	for c := range row {
		elves++

		neighbours := make([]bool, 12)

		// top
		neighbours[0] = g.get(r-1, c-1)
		neighbours[1] = g.get(r-1, c)
		neighbours[2] = g.get(r-1, c+1)
		// bottom
		neighbours[3] = g.get(r+1, c-1)
		neighbours[4] = g.get(r+1, c)
		neighbours[5] = g.get(r+1, c+1)
		// left
		neighbours[6] = neighbours[3]
		neighbours[7] = g.get(r, c-1)
		neighbours[8] = neighbours[0]
		// right
		neighbours[9] = neighbours[5]
		neighbours[10] = g.get(r, c+1)
		neighbours[11] = neighbours[2]

		if utils.All(neighbours, func(x bool, _ int) bool {
			// check if all are empty
			return !x
		}) {
			// elf doesn't do anything
			still++
			continue
		}

		// alternate between which direction we check first
		j := neighbourIter % 4
		for i := j; i < j+4; i++ {
			k := i % 4
			arr := neighbours[k*3 : k*3+3]

			if utils.All(arr, func(x bool, i int) bool {
				// check if empty
				return !x
			}) {
				// elf can move here
			  // ...
```

This seemed super wasteful to run all of this for every elf! I honestly expected the program to fail, and that I would have to rethink the problem.

I repeatedly used the `All` function which I wrote for day 19 (still incomplete):

```go
// python's all() function, kinda
func All[T any](arr []T, cb func(x T, i int) bool) bool {
	for i, v := range arr {
		if !cb(v, i) {
			return false
		}
	}
	return true
}
```

I also, again, wrote a sanity check grid printer, and again used the same logic to update a string at an index:

```go
out[i] = out[i][0:j] + "#" + out[i][j+1:]
```

The grid was quite useful to view, though it was hard to distinguish which point was which, as they were all **"#"**, and the data structure was so minimal (being just `map[int]map[int]struct{}`)

### Day 22

**Difficulty: 10/10** ★★★★★★★★★★

**Time: ~24 hrs**

I found this one to be a hard one to wrap my head around.  Part 1 is 230 LOC. I did manage to get most of the logic inside one function:

```go
// rope around if there is empty space
func (b *board) getNext(r, c int, dir face) (next [2]int, hitWall bool) {
	cur := [2]int{r, c}

	for {
		next = cur
		moveOne(&next, dir)

		// negative numbers
		if dir > 1 {
			if next[0] < 0 {
				next[0] = b.height + next[0]
			} else if next[1] < 0 {
				next[1] = b.width + next[1]
			}
		}

		// wrap around
		next[0] %= b.height
		next[1] %= b.width

		cell := b.get(next[0], next[1])

		if cell == open {
			return
		}
		if cell == wall {
			// return original
			return [2]int{r, c}, true
		}

		cur = next
	}
}
```

Which is actually helped a lot by the `moveOne` function:

```go
func moveOne(pos *[2]int, dir face) {
	switch dir {
	case right:
		pos[1]++
	case left:
		pos[1]--
	case up:
		pos[0]--
	case down:
		pos[0]++
	}
}
```

I opted for a pointer here to avoid dealing with carrying over return values and re-assigning.  I'm unsure about the board.get method:

```go
func (b *board) get(r, c int) (val tile) {
	return b.grid[r][c]
}
```

But I thought it might prove useful later one (so far, no). 

I really liked my data structure for today:

```go

type tile rune

const (
	empty tile = ' '
	open  tile = '.'
	wall  tile = '#'
)

type face int

const (
	right face = iota
	down
	left
	up
)

type cmd int

const (
	move cmd = iota
	rotate
)

type instruction struct {
	command cmd
	value   int // if rotate, then -1 for L and +1 for R
}
```

I turned the instructions into -1 or +1 to rotate, and the `type face` was just iota's.

Anyway, I suspect Part 2 will be hell.

**Part 1 ran in 5ms**

Part 2 is folding the 2d map into a cube. Here's my idea: 

1. breadth-first search from each cube face edge
1. nearest edge of a cube face that isn't already a neighbour is the neighbouring edge
1. we can get the from-direction easily from the initial step, and to-direction is the last step's direction

![cube face neighbours and their edges](https://user-images.githubusercontent.com/1410985/213936714-c54a0f7c-af49-48c5-b2f6-55a2c8f873a3.png "cube face neighbours and their edges")

I accomplished this by using `image.Rectangle` to find non-empty tiles and get their associated squares (cube faces).  Then used a priority queue where each squares neighbouring square is searched with priority of `0`; on each steps search, the priority is increased by `1`, so that we can actually get the nearest cube face for **each** square in the same `loop`.

![cube folding](https://user-images.githubusercontent.com/1410985/212558442-ed9b1524-9cae-482e-888f-b80334084e99.png "cube folding")


```go
neighboursFound := 0
// 6 faces have 4 neighbours each
// we can break after 6*4 neighbours found
neighboursTotal := 6 * 4

for neighboursFound != neighboursTotal {
	cur := pq.Get()

	// stop checking for square neighbours if they have 4
	if cur.square.neighboursFound == 4 {
		continue
	}

	sq, isSquare := squares[cur.currentRect]

	if sq == cur.square {
		// thou shalt not be neighbours with thyself
		continue
	}

	if isSquare {
		// found something!
		added := cur.square.addNeighbour(sq, cur.fromDir, cur.toDir)

		// it's not added if the edge is already taken or the
		// square is already a neighbour
		if added {
			neighboursFound += 2
		}
		continue
	}

	// else, continue path-finding
	for i, vec := range neighbours {
		dir := direction(i)
		nextRect := cur.currentRect.Add(vec)

		pq.PushValue(&state{
			square:      cur.square,
			fromDir:     cur.fromDir,
			toDir:       dir,
			currentRect: nextRect,
			priority:    cur.priority + 1,
		}, cur.priority+1)
	}
}
```

I originally ran into issues with rotating; given the squares have a size of even numbers and I must have run into issues with transforming the correct answer into an `int`; then ran into the same issues with rounding after finding another algorithm for rotating.

**Issue with floats (off by one error):**

```go
// rotate a point around an origin algorithm
px, py := float64(curPoint.X), float64(curPoint.Y)
origin := rect.Min.Add(size.Div(2))
ox, oy := float64(origin.X)-0.5, float64(origin.Y)-0.5
// angle in radians 🤦‍♂️
// inverse to go clockwise
theta := -(math.Pi / 2) * float64(rotations)

// TODO: Need math.Round here!
rx := math.Cos(theta)*(px-ox) - math.Sin(theta)*(py-oy) + ox
ry := math.Sin(theta)*(px-ox) + math.Cos(theta)*(py-oy) + oy

next = [2]int{int(rx), int(ry)}
fmt.Println(rx, ry, next)
```

The `TODO` there says it all, I hope; here's what I was getting in the test:

```sh
=== RUN   TestCubeFoldFifty
9.000000000000005 99 [9 99]
199 8.999999999999998 [199 8]
    day-22_test.go:76: wanted [0 59], got: [0 58]
--- FAIL: TestCubeFoldFifty (2.26s)
```

I think my logic was straight-forward; however, it was complex.  Also my data structure probably should have changed.  In too many places, I was swapping between `r int, c int` and `[2]int` and `image.Point`.  Given the 3d-nature of the puzzle, I probably should have stuck with `image.Point`, as it became incredibly useful to work with it and `image.Rectangle`.

Here's part of the ridiculous swapping, and the use of `image.Rectangle.Mod`:

```go
next = [2]int{int(rx), int(ry)}

// move in new direction by one
moveOne(&next, newDir)

// annoying transitions here
curPoint = image.Point{next[0], next[1]}

// mod to next cube face (mod is virtually teleporting)
curPoint = curPoint.Mod(neighbour.sq.rect)

// set next for outer for loop to check for walls
next = [2]int{curPoint.X, curPoint.Y}
```

So, let's try to visualize this, suppose we're moving from A, moving up, which should line up with B, moving right:

```sh
  B.
  ..
A...
....
```

We take A's square, and rotate it (90deg clockwise in this case) so that the edge is lined up with B's square (the inverse direction really, so they'd touch if they were next to each other):

```sh
  B.
  ..
.A..
....
```

Then, the `moveOne` function moves in the new direction (**right**):

```sh
  B.
  ..
..A.
....
```

Then, we use `image.Rectangle.Mod` to run the modulus operator on both `X` and `Y` of a given `image.Point`, essentially moving the point `A` into `B`'s square, in the equivalent place:

```sh
  A.
  ..
....
....
```

This works in non-2x2-squares too, obviously.

This solution runs fast enough:

```sh
1 | 131052 (8.646415ms)
2 | 4578 (1.600803891s)
```

### Day 21

**Difficulty: 2/10** ★★☆☆☆☆☆☆☆☆

**Time: ~60 mins**

Part 1 was quite simple.  Pretty simple to parse by just using `strings.Fields(line)` and checking the length of the fields.  The structure was a bit awkward, because the monkeys either had a value or had to derive a value:

```go
type monkey struct {
	value              int
	hasValue           bool
	needs              [2]string
	operator, neededBy string
}

type monkeys map[string]monkey

// ...

fields := strings.Fields(line)
name := fields[0][:len(fields[0])-1]

monkey := monkey{}

if len(fields) == 2 {
	// monkey has number
	monkey.hasValue = true
	monkey.value = utils.ParseInt(fields[1])
} else {
	// monkey depends on other monkeys
	a, op, b := fields[1], fields[2], fields[3]
	monkey.needs = [2]string{a, b}
	monkey.operator = op
}

monkeys[name] = monkey
```

I thought the `hasValue` field was quite annoying; mostly because empty int's are `0`, which would also be valid.  I've stuck with the idea of identifying objects via their string key's instead of defining `needs` as `[2]*monkey`.  This made for some awkward code later, as I was constantly referencing `monkeys[name]`, but it made for much easier parsing.

Today was very lazy: it was kind of like I was talking before thinking, and trying to figure out what I said.  

Part 1's `getMonkey` is a recursive function, which I find very clean.  I used it to update the monkeys for part 2, so that I didn't have to derive values again.

The `whatToYell` function is very sloppy and filled with comments, mostly to try to figure out my flailing ideas.  I knew I could figure out the value by creating a stack from "humn", and then iterating back down the stack, from "root", to figure out what each value is supposed to be and then doing simple math on it.  There were a lot of variables which were hard to name, and hard to keep track of.  

```go
// number we need the next monkey to yell
var nextNeed int

switch cur.operator {
case "=":
	// we need to equal the other
	nextNeed = other
case "/":
	// mult or divide depending on whether we're looking for
	// numerator or denominator
	if index == 1 {
		// x / other = prevNeed
		nextNeed = other * prevNeed
	} else {
		// other / x = prevNeed
		nextNeed = other / prevNeed
	}
case "+":
	nextNeed = prevNeed - other
case "*":
	nextNeed = prevNeed / other
case "-":
	if index == 1 {
		// x - other = prevNeed
		nextNeed = prevNeed + other
	} else {
		// other - x = prevNeed
		nextNeed = other - prevNeed
	}
}
```

This seemed very easy.

Times: `~2ms` each part

### Day 20

**Difficulty: 6/10** ★★★★★★☆☆☆☆

**Time: ~90 mins**

Ordering the list for Part 1 was at first hard to grasp, but I seemed to luck out.  I iterated the list, kept track if it was `visited` or not, and then either incremented `i` by 1 or 0 (to check the likely next unvisited value).

```go
for i := 0; i < len(ordered); {
	s := ordered[i]

	if s.visited {
		i++
		continue
	}
	newI := (i + s.value) % (len(ordered) - 1)
	if newI < 0 {
		newI %= len(ordered)
		newI = len(ordered) + newI - 1
	}
	if newI == 0 {
		newI = len(ordered) - 1
	}

	// remove
	ordered = append(ordered[:i], ordered[i+1:]...)
	// insert
	ordered = append(ordered[:newI], append([]sorting{{
		value:   s.value,
		visited: true,
	}}, ordered[newI:]...)...)

	// don't adjust i; we revisit this index
}
```

Inserting array values, and removing array values is still a little annoying.  We'll see how painful Part 2 is later.

Part 2 was a little harder to grasp.  I found it very hard to debug too, since all the numbers are obscured by the `decryption key`.  Also, keeping track of a list, which is being shuffled 10 times is difficult.  I eventually settled on keeping two lists: a slice of original nodes which represented the numbers, and a `container/list` linked list which had pointers to the numbers, and represented the current order.  This worked!  I did however have to get the current index on each iteration, and cast(?) the `any` types back to get values.  This ran in just over a second for part 2:

```sh
1 | 27726 (82.975042ms)
2 | 4275451658004 (1.400359475s)
```

I realized this was over-engineered and could just have two slices instead, keeping a lot of the logic for new indices as I had in Part 1:

```go
func reOrder(orig []int, mixes int) []int {
	length := len(orig)

	// original order with pointers
	originalNodes := make([]*int, length)
	// current order with pointers to the original
	currentOrder := make([]*int, length)

	for i := range orig {
		originalNodes[i] = &orig[i]
		currentOrder[i] = originalNodes[i]
	}

	for mixes > 0 {
		mixes--
		for i := 0; i < length; i++ {
			node := originalNodes[i]

			// need to get current index; so have to loop over the slice, I believe
			oldI := indexOf(currentOrder, node)
			newI := (oldI + *node) % (length - 1)

			if newI < 0 {
				newI = length - 1 + newI
			}

			// remove
			currentOrder = append(currentOrder[:oldI], currentOrder[oldI+1:]...)
			// insert
			currentOrder = append(currentOrder[:newI], append([]*int{node}, currentOrder[newI:]...)...)
		}
	}

	out := make([]int, length)

	for i, v := range currentOrder {
		out[i] = *v
	}

	return out
}
```

After converting into just two slices, I got the times down a little bit (half the time for part 2): 

```sh
1 | 27726 (77.528845ms)
2 | 4275451658004 (722.292438ms)
```

This was challenging, because I thought I could keep track of the index internally somewhere; but what I didn't appreciate is that I'd have to alter the index of every node between `oldIndex` and `newIndex` when moving. Adding an `indexOf` function seemed like a fine cop-out.

### Day 19

**Difficulty: 10/10** ★★★★★★★★★★

**Time: >10 hrs**

**Part 1 finished at 4.68s**.  This was very difficult.  I couldn't get the tests running in a decent time.  I tried to prune invalid states and messed up a few times. The data structure I picked was rather difficult.

I started with `[4][4]int` as the blueprint for how to build a robot: 4 types of robot which could require any of 4 types of resource, and then finally a count of those resources needed.

I switched this with a map:

```go
type resource int

const (
	ore resource = 1 << iota
	clay
	obsidian
	geode
)

type blueprint struct {
	// robot_type->cost_type->cost: ore,clay,obsidian,geode
	robots map[resource]map[resource]int
	// bitmask saves time with comparing
	robotbitmask map[resource]int
	// figure out when we should stop buying robots
	maxRobots map[resource]int
}
```

This was still annoying.  I still had to iterate everything. In `getNextStates` I did use the bitmasks to save on some for loops:

```go
// BASICALLY: which robots could we buy next and when can we buy them?
// and can it happen before the `end`?
for robot, bitmask := range bp.robotbitmask {
	// if we have every robot that is needed
	if bitmask&cur.robotbits == bitmask {
```

Otherwise, I'd have to iterate the current state's robots to see if it were possible to buy any robot.

Next states was not minute-by-minute but robot-purchase-by-robot-purchase:

```go
maxTime := 0
for res, num := range bp.robots[robot] {
	// costs `num` of `res`
	time := int(math.Ceil(float64(num-cur.resources[res]) / float64(cur.robots[res])))

	if (cur.time + time) >= end-1 {
		// we can't actually buy this before the end of time
		// skip ahead to the end
		maxTime = end - clone.time - 1
		break
	}

	// we have to wait for the resources which take the longest
	if time > maxTime {
		maxTime = time
	}
	// reduce resources by cost
	// (could go negative, but we'll increase it soon)
	clone.resources[res] -= num
}

// increment time
totalTime := maxTime + 1
clone.time += totalTime
```

For pruning invalid states, I knew I could keep a cache of visited states and ignore those that I already knew had less geodes.  I also removed states that had no geodes at the end or couldn't buy a geode robot at the last minute.

I tried to keep track of earliest geode robot purchase (as I read on Reddit), and prune any state that didn't have a geode robot at that time, and this passed the tests, but it failed on my input data.  My answer with that pruning was 10 less than it should have been.

#### Day 19 Update (Part 2)

Wow.  I had to overhaul the code I wrote, and switch from breadth-first-search to depth-first-search. Finally decent results.

```sh
1 | 1144 (1.550023196s)
2 | 19980 (2.933193562s)
```

I simplified state quite a bit:

```go
type state struct {
	time              int
	resources, robots map[resource]int
}
```

This was my **first time** recursing inside a closure:

```go

// *first time* recurse inside a closure :D
var dfs func(st state)

func (bp blueprint) bestPath(timeLimit int) (best int) {
	cache := map[string]struct{}{}

	dfs = func(st state) {
		// ...
		dfs(next)
	}
	// ...
}
```

I iterated the resources in order this time, to prioritize states with geodes more:

```go
var resources = [4]resource{geode, obsidian, clay, ore}

// ... 

for _, robot := range resources {
```

The biggest time saver was checking if a state could *possibly* get better (i.e. could buy a geode robot and get max geodes for the rest of time):

```go
// pruning
// could it possibly be better?
if best > 0 {
	max := getMaxGeodesFromTimeLeft(next, timeLimit-next.time)
	if max <= best {
		// ! this saved perhaps the most time
		continue
	}
}

// ...

func getMaxGeodesFromTimeLeft(st state, time int) int {
	// current geodes
	cur := st.resources[geode]
	// previous robots accumulate geodes
	rate := st.robots[geode] * time
	// new robots accumulate geodes
	// example: given diff == 5 -> 4+3+2+1 == 5*4/2
	maxBots := (time * (time - 1)) / 2

	return cur + rate + maxBots
}
```

And(!) my hash function was inaccurate.  I was previously omitting `geode` robots and resources so that I could compare with previous values to keep the best ones.  That produced incorrect answers.

I'm glad it's over.

One thing that I was told on Reddit, paraphrased: BFS is good for specific targets, and DFS is good for some "best of" scenario, where there is no explicit end goal.

### Day 18

**Difficulty: 5/10** ★★★★★☆☆☆☆☆

**Time: ~2 hrs**

```sh
1 | 4536 (123.357198ms)
2 | 2606 (14.013846ms)
```

I originally had a data structure of `[]cube` but opted for `map[cube]struct{}` for ease of checking if a cube was "lava" or not.  Most of today was verbose, with lots of lines for min/max, in order to get the bounds of the cube I needed to "flood fill".

I did another recursive function within a closure.

Overall seemed quite easy, though I did read up on other solutions before I started this one, so I knew what I was about to get into.

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

I should say I initially tried to put walls and the floor into the space variable, but second guessed it during refactoring.  Seeing as Part 1 is now 600x slower, I might revisit this if and when I care.

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

#### Update Day 16

So, I apparently did not set up the priority queue correctly:

```diff
// sets index automatically
func (pq *PriorityQueue[T]) PushValue(value *T, priority int) {
	newItem := &Item[T]{
		value,
		priority,
		0,
	}

-	pq.Push(pq, newItem)
+	heap.Push(pq, newItem)
}
```

After fixing this, the speed improved by 10s for Part 1:

```sh
1 | 1915 (29.5305s)
2 | 2772 (8.6158s)
```

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
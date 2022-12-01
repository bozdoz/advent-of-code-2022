# What Am I Learning Each Day?

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

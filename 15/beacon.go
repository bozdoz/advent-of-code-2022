package main

import (
	"fmt"
	"image"
	"sort"

	"github.com/bozdoz/advent-of-code-2022/types"
	"github.com/bozdoz/advent-of-code-2022/utils"
)

type rel int

const (
	DISJOINT rel = iota - 1 // sensors aren't close at all
	TOUCH                   // sensors touch (specifically they are 1-pixel apart)
	OVERLAP                 // sensors overlap (either actually touch, or overlap)
)

type sensor struct {
	coords    image.Point
	manhattan int                // distance from closest beacon
	touches   int                // count of sensors this touches (used for sorting)
	overlaps  types.Set[*sensor] // lazy tracking to ensure de-duping
}

type space struct {
	sensors                []*sensor
	beacons                types.Set[image.Point] // part 1 requires tracking beacons
	xmin, xmax, ymin, ymax int
}

func parseInput(data inType) space {
	sensors := []*sensor{}
	beacons := types.Set[image.Point]{}

	var xmin, xmax, ymin, ymax int

	for i, line := range data {
		var sx, sy, bx, by int

		num, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sx, &sy, &bx, &by)

		if num != 4 || err != nil {
			panic(fmt.Sprintf("parsed %d vars from: %q, with error: %v", num, line, err))
		}

		coords := image.Point{sx, sy}
		beacon := image.Point{bx, by}

		// determine area where beacons cannot be
		// AND determine min/max of space
		d := utils.ManhattanDistance(coords, beacon)

		sensors = append(sensors, &sensor{
			coords:    coords,
			manhattan: d,
			overlaps:  types.Set[*sensor]{},
		})

		beacons.Add(beacon)

		min := coords.Sub(image.Point{d, d})
		max := coords.Add(image.Point{d, d})

		if i == 0 {
			xmin, ymin = min.X, min.Y
			xmax, ymax = max.X, max.Y
		} else {
			xmin = utils.Min(xmin, min.X)
			ymin = utils.Min(ymin, min.Y)
			xmax = utils.Max(xmax, max.X)
			ymax = utils.Max(ymax, max.Y)
		}
	}

	// keep track of sensor touches (edges are along a common path)
	// keep track of overlaps
	// scan sensors sorted by num of touches, against overlapping sensors
	for i := 0; i < len(sensors); i++ {
		a := sensors[i]

		for j := 0; j < len(sensors); j++ {
			if i == j {
				continue
			}
			b := sensors[j]

			switch sensorsTouchOrOverlap(*a, *b) {
			case TOUCH:
				a.touches++
				b.touches++
			case OVERLAP:
				a.overlaps.Add(b)
				b.overlaps.Add(a)
			}
		}
	}

	space := space{
		sensors: sensors,
		beacons: beacons,
		xmin:    xmin,
		xmax:    xmax,
		ymin:    ymin,
		ymax:    ymax,
	}

	return space
}

// we want to keep track of sensors that don't *actually* touch, but
// share the same 1-pixel outer edge; but also,
// track which sensors *actually* touch or overlap areas
func sensorsTouchOrOverlap(a, b sensor) rel {
	between := utils.ManhattanDistance(a.coords, b.coords)

	// -1 to omit the sensor space itself
	dist := between - a.manhattan - b.manhattan - 1

	switch {
	// check that exactly 1 pixel is between the sensor areas
	case dist == 1:
		return TOUCH
	case dist < 1:
		return OVERLAP
	default:
		return DISJOINT
	}
}

// part 1 check if pixels could be a beacon or not
// note: beacons return true because beacons could also be beacons
func (space *space) couldBeBeacon(x, y int) bool {
	point := image.Point{x, y}

	// check if definitely is beacon
	if space.beacons.Has(point) {
		// is, in fact, beacon, so definitely could be
		return true
	}

	// run manhattan all over the sensors
	for _, sensor := range space.sensors {
		// sensor cannot be beacon
		if point.Eq(sensor.coords) {
			return false
		}

		d := utils.ManhattanDistance(point, sensor.coords)

		if d <= sensor.manhattan {
			return false
		}
	}

	return true
}

// there's only one pixel that isn't covered by the sensors
// that is the missing beacon
func (space *space) findMissingBeacon() image.Point {
	sensors := space.sensors

	// start with most touched sensor
	sort.Slice(sensors, func(i, j int) bool {
		return sensors[i].touches > sensors[j].touches
	})

	// maybe lazy: passing channel to avoid checking return values
	found := make(chan image.Point)

	// need a goroutine for channel to receive the value
	go func() {
		for _, sensor := range sensors {
			space.scanClockwise(sensor, found)
		}
	}()

	return <-found
}

// diagonal directions, clock-wise
var (
	tr = image.Point{1, 1}
	rb = image.Point{-1, 1}
	bl = image.Point{-1, -1}
	lt = image.Point{1, -1}
)

var directions = [...]image.Point{tr, rb, bl, lt}

func (space *space) scanClockwise(sensor *sensor, found chan image.Point) {
	// add 1 to manhattan, because the missing beacon MUST be
	// exactly 1 pixel beyond the sensor scan
	d := sensor.manhattan + 1
	top := sensor.coords.Add(image.Point{0, -d})
	right := sensor.coords.Add(image.Point{d, 0})
	bottom := sensor.coords.Add(image.Point{0, d})
	left := sensor.coords.Add(image.Point{-d, 0})

	vertices := [...]image.Point{top, right, bottom, left}

	for i := range vertices {
		var start, finish image.Point

		start = vertices[i]

		if i == len(vertices)-1 {
			// left -> top
			finish = vertices[0]
		} else {
			finish = vertices[i+1]
		}

		dir := directions[i]

		// iterate diagonal pixels around the sensor
		for !start.Eq(finish) {
			// check overlapping sensors within space
			if space.contains(start) && !withinSensorRange(start, sensor.overlaps) {
				// found missing beacon!
				found <- start
			}

			start = start.Add(dir)
		}
	}
}

// check if point is inside of min/max space coordinates
func (space *space) contains(p image.Point) bool {
	return p.X >= space.xmin && p.X <= space.xmax && p.Y >= space.ymin && p.Y <= space.ymax
}

// check point against overlapping (not touching and not disjoint) sensors
func withinSensorRange(point image.Point, sensors types.Set[*sensor]) bool {
	// run manhattan all over the sensors
	for sensor := range sensors {
		d := utils.ManhattanDistance(point, sensor.coords)

		if d <= sensor.manhattan {
			return true
		}
	}

	return false
}

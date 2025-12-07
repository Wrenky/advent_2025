package main

import (
	"advent/helpers"
	"advent/helpers/grid"
	log "log/slog"
	"strings"
	"time"

	"github.com/samber/lo"
)

// This needs to change to match the input
func parseInput(input string) grid.Grid[string] {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		return strings.Split(line, "")
	})
}

func FindStart(g grid.Grid[string]) grid.Coord {
	_, y, _ := lo.FindIndexOf(g[0], func(i string) bool {
		return i == "S"
	})
	return grid.Coord{X: 0, Y: y}
}

func Beam(g grid.Grid[string], beams []grid.Coord, splitCount int) int {
	if len(beams) == 0 {
		return splitCount
	}
	nextBeams := lo.Flatten(lo.FilterMap(beams, func(b grid.Coord, _ int) ([]grid.Coord, bool) {
		next := grid.Add(b, grid.DOWN)
		if !grid.InBounds(g, next) {
			// Done!
			return []grid.Coord{}, false
		}
		if grid.Get(g, next) == "^" {
			// split!
			return []grid.Coord{
				grid.Add(next, grid.LEFT),
				grid.Add(next, grid.RIGHT),
			}, true
		}
		return []grid.Coord{next}, true
	}))
	if len(nextBeams) == 0 {
		return splitCount
	}
	return Beam(g, lo.Uniq(nextBeams), splitCount+len(nextBeams)-len(beams))
}

// Memoization cache
var cache = make(map[grid.Coord]int)

func Timelines(g grid.Grid[string], beam grid.Coord) int {
	if result, ok := cache[beam]; ok {
		return result
	}
	next := grid.Add(beam, grid.DOWN)
	var result int
	switch {
	case !grid.InBounds(g, next):
		// Ending case!
		result = 1
	case grid.Get(g, next) == "^":
		// Splitter case- Do both sides
		l := grid.Add(next, grid.LEFT)
		r := grid.Add(next, grid.RIGHT)
		result = Timelines(g, l) + Timelines(g, r)
	default:
		// Straight case
		result = Timelines(g, next)
	}
	cache[beam] = result
	return result
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	g := parseInput(data)
	for _, v := range g {
		log.Debug("", "line", v)
	}
	start := FindStart(g)
	log.Debug("Starting coord", "start", start, "startVal", grid.Get(g, start))

	// Part 1
	pre1 := time.Now()
	splits := Beam(g, []grid.Coord{grid.Add(start, grid.DOWN)}, 0)
	post1 := time.Now()
	log.Info("Part1", "answer", splits, "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	timelines := Timelines(g, grid.Add(start, grid.DOWN))
	post2 := time.Now()
	log.Info("Part2", "answer", timelines, "time", post2.Sub(pre2))
}

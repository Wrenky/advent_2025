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

func AllDirections(p grid.Coord) []grid.Coord {
	return lo.Map(grid.ALLDIRECTIONS, func(c grid.Coord, _ int) grid.Coord {
		return grid.Add(p, c)
	})
}

// Can we forklift this point?
func Forkliftable(g grid.Grid[string], p grid.Coord) bool {
	if !grid.InBounds(g, p) {
		log.Warn("Invlalid point passed in")
		return false
	}
	if grid.Get(g, p) != "@" {
		return false
	}
	// Check that all the surrounding points are @s
	return len(lo.Filter(AllDirections(p), func(cand grid.Coord, _ int) bool {
		return grid.InBounds(g, cand) && (grid.Get(g, cand) == "@")
	})) < 4
}

// ================================================================
// Part 2 only, remove all forkliftable points
// Forklift the points and return an updated grid + removed count
func ForkRemove(g grid.Grid[string]) (grid.Grid[string], int) {

	removable := lo.Filter(grid.AllPoints(g), func(p grid.Coord, _ int) bool {
		return Forkliftable(g, p)
	})

	// New grid!
	gnext := grid.Clone(g)
	for _, c := range removable {
		gnext[c.X][c.Y] = "."
	}
	return gnext, len(removable)
}

// Part2: Forkremove until we cant!
func Removal(g grid.Grid[string], removed int) int {
	gnext, r := ForkRemove(g)
	if r == 0 {
		return removed
	}
	return Removal(gnext, removed+r)
}

//================================================================

// ================================================================
// Part 1 code, count forkliftable on a single grid
func CountForkliftable(g grid.Grid[string]) int {
	return len(lo.Filter(grid.AllPoints(g), func(p grid.Coord, _ int) bool {
		return Forkliftable(g, p)
	}))
}

//================================================================

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	g := parseInput(data)
	for _, v := range g {
		log.Debug("", "row", v)
	}

	pre1 := time.Now()
	count := CountForkliftable(g)
	post1 := time.Now()
	log.Info("Part1", "answer", count, "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	removed := Removal(g, 0)
	post2 := time.Now()
	log.Info("Part2", "answer", removed, "time", post2.Sub(pre2))
}

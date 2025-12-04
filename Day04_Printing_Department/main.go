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

var (
	UP    = grid.Coord{X: -1, Y: 0}
	DOWN  = grid.Coord{X: 1, Y: 0}
	LEFT  = grid.Coord{X: 0, Y: -1}
	RIGHT = grid.Coord{X: 0, Y: 1}
)

func AllDirections(p grid.Coord) []grid.Coord {
	up := grid.Add(p, UP)
	down := grid.Add(p, DOWN)
	return []grid.Coord{
		up,
		grid.Add(up, LEFT),
		grid.Add(up, RIGHT),
		grid.Add(p, LEFT),
		grid.Add(p, RIGHT),
		down,
		grid.Add(down, LEFT),
		grid.Add(down, RIGHT),
	}
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

// Forklift the points and return an updated grid + removed count
func ForkRemove(g grid.Grid[string]) (grid.Grid[string], int) {

	removable := lo.Flatten(lo.Map(g, func(row []string, i int) []grid.Coord {
		return lo.FilterMap(row, func(x string, j int) (grid.Coord, bool) {
			p := grid.Coord{
				X: i,
				Y: j,
			}
			return p, Forkliftable(g, p)
		})
	}))

	// New grid!
	gnext := grid.Clone(g)
	for _, c := range removable {
		gnext[c.X][c.Y] = "."
	}
	return gnext, len(removable)
}

// Forkremove until we cant!
func Removal(g grid.Grid[string], removed int) int {
	gnext, r := ForkRemove(g)
	if r == 0 {
		return removed
	}
	return Removal(gnext, removed+r)
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	g := parseInput(data)
	for _, v := range g {
		log.Debug("", "row", v)
	}

	pre1 := time.Now()
	count := lo.Sum(lo.Map(g, func(row []string, i int) int {
		return len(lo.Filter(row, func(x string, j int) bool {
			return Forkliftable(g, grid.Coord{X: i, Y: j})
		}))
	}))
	post1 := time.Now()
	log.Info("Part1", "answer", count, "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	removed := Removal(g, 0)
	post2 := time.Now()
	log.Info("Part2", "answer", removed, "time", post2.Sub(pre2))
}
